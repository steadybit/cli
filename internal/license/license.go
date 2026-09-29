// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package license implements the `license` commands.
package license

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/resource"
	"github.com/steadybit/cli/v6/internal/table"
)

var errNotAdmin = errors.New("The license needs an admin access token.")

type ShowOptions struct {
	Type string
}

type summary struct {
	License *struct {
		LicenseType string `json:"licenseType"`
		OrderNumber string `json:"orderNumber"`
		ValidFrom   string `json:"validFrom"`
		ValidTo     string `json:"validTo"`
	} `json:"license"`
	Expires   *time.Time `json:"expires"`
	TenantKey string     `json:"tenantKey"`
	Features  []feature  `json:"features"`
}

type feature struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Usage     *int64 `json:"usage"`
	SoftLimit *int64 `json:"softLimit"`
	HardLimit *int64 `json:"hardLimit"`
}

// Show prints the license of the tenant and how much of each limit is used.
func Show(ctx context.Context, c *platform.Client, o ShowOptions) error {
	if _, err := output.ResolveDatatype(o.Type, ""); err != nil {
		return err
	}
	body, _, err := platform.Read(c.GetLicenseSummary(ctx))
	if platform.IsStatus(err, http.StatusForbidden) {
		return errNotAdmin
	}
	if err != nil {
		return platform.Failed(err, "Failed to get the license")
	}
	if resource.Machine(o.Type) {
		return resource.PrintJSONValue(body, o.Type)
	}
	var s summary
	if err := json.Unmarshal(body, &s); err != nil {
		return fmt.Errorf("Failed to read the license: %w", err)
	}
	if s.License == nil || s.License.LicenseType == "" || s.License.LicenseType == "NONE" {
		fmt.Println("The tenant has no license.")
		return nil
	}
	fmt.Println(sentence(s, time.Now()))

	// The platform sends the features in no particular order.
	sort.Slice(s.Features, func(i, j int) bool { return s.Features[i].Name < s.Features[j].Name })
	// Features without a limit are only on or off; they make a list, not table rows.
	var included []string
	t := table.New(
		table.Column{Name: "feature", Title: "Limit", Alignment: table.Left},
		table.Column{Name: "used", Title: "Used", Alignment: table.Right},
		table.Column{Name: "limit", Title: "Licensed", Alignment: table.Right},
	)
	limited := false
	for _, f := range s.Features {
		if f.Type == "SIMPLE" {
			included = append(included, f.Name)
			continue
		}
		limited = true
		limit, color := limitOf(f)
		used := ""
		if f.Usage != nil {
			used = fmt.Sprint(*f.Usage)
		}
		t.AddRow(color, table.Cell("feature", f.Name), table.Cell("used", used), table.Cell("limit", limit))
	}
	if limited {
		t.Print()
	}
	if len(included) > 0 {
		fmt.Printf("Included: %s\n", strings.Join(included, ", "))
	}
	return nil
}

// sentence describes the license from the parts the platform sent: an order number or
// a tenant key can be missing, and must not leave a gap or a dangling "of tenant".
func sentence(s summary, now time.Time) string {
	text := title(s.License.LicenseType) + " license"
	if s.License.OrderNumber != "" {
		text += " " + s.License.OrderNumber
	}
	if s.TenantKey != "" {
		text += " of tenant " + s.TenantKey
	}
	switch from, to := s.License.ValidFrom, s.License.ValidTo; {
	case from != "" && to != "":
		text += ", valid from " + from + " to " + to
	case from != "":
		text += ", valid from " + from
	case to != "":
		text += ", valid to " + to
	}
	return text + expiry(s.Expires, now) + "."
}

// limitOf is the licensed amount of a feature and whether its usage is highlighted.
// The feature's type says which limit applies; the platform can send the other field
// too. A hard limit is highlighted once reached, since nothing more can be added; a
// soft one only once exceeded.
func limitOf(f feature) (string, table.Color) {
	highlight := func(over bool) table.Color {
		if over {
			return table.Red
		}
		return table.Default
	}
	switch f.Type {
	case "SOFT_LIMIT":
		if f.SoftLimit != nil {
			return fmt.Sprintf("%d (soft)", *f.SoftLimit), highlight(f.Usage != nil && *f.Usage > *f.SoftLimit)
		}
	case "HARD_LIMIT":
		if f.HardLimit != nil {
			return fmt.Sprint(*f.HardLimit), highlight(f.Usage != nil && *f.Usage >= *f.HardLimit)
		}
	}
	return "unlimited", table.Default
}

func title(licenseType string) string {
	return strings.ToUpper(licenseType[:1]) + strings.ToLower(licenseType[1:])
}

// expiry warns of a license that has run out or is about to: a pipeline's runs stop with it.
func expiry(expires *time.Time, now time.Time) string {
	if expires == nil {
		return ""
	}
	left := expires.Sub(now)
	switch {
	case left <= 0:
		return ", expired"
	case left < 30*24*time.Hour:
		days := int(left.Hours()/24) + 1
		if days == 1 {
			return ", expires in 1 day"
		}
		return fmt.Sprintf(", expires in %d days", days)
	}
	return ""
}

type ReportOptions struct {
	Output string
}

// Report downloads the license usage report, a zip archive of the tenant's usage over
// each license period, as the platform names it unless an output file is given. Only a
// file given with -o is overwritten: the platform's name is not the user's choice, and
// could be that of any file in the current directory.
func Report(ctx context.Context, c *platform.Client, o ReportOptions) error {
	// The report covers every license period; building it takes longer than an API response.
	ctx = platform.WithTimeout(ctx, 5*time.Minute)
	content, resp, err := platform.Read(c.GetReport(ctx))
	if platform.IsStatus(err, http.StatusForbidden) {
		return errNotAdmin
	}
	if err != nil {
		return platform.Failed(err, "Failed to download the license report")
	}
	file, flags := o.Output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC
	if file == "" {
		file, flags = fileName(resp.Header.Get("Content-Disposition")), os.O_WRONLY|os.O_CREATE|os.O_EXCL
	}
	if dir := filepath.Dir(file); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(file, flags, 0o644)
	if errors.Is(err, os.ErrExist) {
		return fmt.Errorf("%s already exists. Pass -o %s to overwrite it, or -o another file.", file, file)
	}
	if err != nil {
		return err
	}
	_, err = f.Write(content)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	fmt.Printf("License report written to %s.\n", file)
	return nil
}

// fileName takes the name the platform gives the report, but only as one path segment
// and not as a hidden file: the header must not decide where on disk the file goes.
func fileName(disposition string) string {
	_, params, err := mime.ParseMediaType(disposition)
	if err != nil {
		return "license-report.zip"
	}
	name := strings.TrimLeft(output.PathSegment(strings.ReplaceAll(params["filename"], `\`, "/")), ".")
	if name == "" || name == "_" {
		return "license-report.zip"
	}
	return name
}

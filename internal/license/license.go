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
	fmt.Printf("%s license %s of tenant %s, valid from %s to %s%s.\n", title(s.License.LicenseType), s.License.OrderNumber, s.TenantKey,
		s.License.ValidFrom, s.License.ValidTo, expiry(s.Expires, time.Now()))

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
		limit, color := "unlimited", table.Default
		switch {
		case f.HardLimit != nil:
			limit = fmt.Sprint(*f.HardLimit)
			if f.Usage != nil && *f.Usage > *f.HardLimit {
				color = table.Red
			}
		case f.SoftLimit != nil:
			limit = fmt.Sprintf("%d (soft)", *f.SoftLimit)
			if f.Usage != nil && *f.Usage > *f.SoftLimit {
				color = table.Red
			}
		}
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
		return fmt.Sprintf(", expires in %d days", int(left.Hours()/24)+1)
	}
	return ""
}

type ReportOptions struct {
	Output string
}

// Report downloads the license usage report, a zip archive of the tenant's usage over
// each license period, as the platform names it unless an output file is given.
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
	file := o.Output
	if file == "" {
		file = fileName(resp.Header.Get("Content-Disposition"))
	}
	if dir := filepath.Dir(file); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	if err := os.WriteFile(file, content, 0o644); err != nil {
		return err
	}
	fmt.Printf("License report written to %s.\n", file)
	return nil
}

// fileName takes the name the platform gives the report, but only as one path segment:
// the header must not decide where on disk the file goes.
func fileName(disposition string) string {
	_, params, err := mime.ParseMediaType(disposition)
	if err != nil || params["filename"] == "" {
		return "license-report.zip"
	}
	return output.PathSegment(strings.ReplaceAll(params["filename"], `\`, "/"))
}

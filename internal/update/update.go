// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package update tells someone at a terminal, at most once a day, that a newer release
// of the CLI exists. With npm gone, nothing else would: a downloaded binary stays as it is.
package update

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

// LatestURL redirects to the tag of the latest release. Asking where it redirects costs
// none of the GitHub API's rate limit. Tests point it elsewhere.
var LatestURL = "https://github.com/steadybit/cli/releases/latest"

// Now is the clock the daily check is measured against. Tests replace it.
var Now = time.Now

// Interactive reports whether a notice would be seen, and not end up in a CI log or in
// output a script reads. Tests replace it.
var Interactive = func() bool {
	return !turnedOff(os.Getenv) && term.IsTerminal(int(os.Stderr.Fd()))
}

// turnedOff is whether the environment rules the check out: asked to, or in CI.
func turnedOff(getenv func(string) string) bool {
	switch strings.ToLower(getenv("STEADYBIT_NO_UPDATE_CHECK")) {
	case "", "0", "false":
	default:
		return true
	}
	for _, ci := range []string{"CI", "JENKINS_URL", "TF_BUILD", "BUILDKITE"} {
		if getenv(ci) != "" {
			return true
		}
	}
	return false
}

const every = 24 * time.Hour

// How long a command waits, once a day, for the check before it ends without it.
const patience = time.Second

// The request gives up earlier than the command waits, so that even a request that runs
// out of time leaves room to record the check before the command ends.
const fetchTimeout = 800 * time.Millisecond

// fetch asks where the latest release is. Tests replace it, so that none depends on how
// fast a refused connection fails, which differs between systems.
var fetch = fetchLatest

var release = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)

type state struct {
	CheckedAt time.Time `json:"checkedAt"`
	Latest    string    `json:"latest"`
}

// Start checks for a newer release when a check is due, while the command runs. The
// function it returns prints the notice, after the command, from what is known by then.
// A build that is not a release, such as one from `go install ...@main`, is never told.
func Start(current, cacheFile string) (notice func(io.Writer)) {
	if !release.MatchString(current) || !Interactive() {
		return func(io.Writer) {}
	}
	known := read(cacheFile)
	done := make(chan state, 1)
	// A check dated in the future, from a wrong clock or a copied cache, is due: without
	// that it would not be repeated until the clock caught up.
	if since := Now().Sub(known.CheckedAt); since >= 0 && since < every {
		done <- known
	} else {
		go func() {
			checked := state{CheckedAt: Now(), Latest: known.Latest}
			ctx, cancel := context.WithTimeout(context.Background(), fetchTimeout)
			defer cancel()
			if latest, err := fetch(ctx); err == nil {
				checked.Latest = latest
			}
			// Written even when offline, so an unreachable GitHub is not asked again
			// before tomorrow.
			write(cacheFile, checked)
			done <- checked
		}()
	}
	return func(w io.Writer) {
		select {
		case known = <-done:
		case <-time.After(patience):
		}
		if newer(known.Latest, current) {
			fmt.Fprintf(w, "\nA new release of the Steadybit CLI is available: %s → %s\n%s\n", current, known.Latest, howToUpdate())
		}
	}
}

func fetchLatest(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, LatestURL, nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	_ = resp.Body.Close()
	tag := strings.TrimPrefix(path.Base(resp.Header.Get("Location")), "v")
	if !release.MatchString(tag) {
		return "", fmt.Errorf("no release tag in %q", resp.Header.Get("Location"))
	}
	return tag, nil
}

// newer compares release versions by number, so that 6.10.0 is after 6.9.0.
func newer(latest, current string) bool {
	l, c := release.FindStringSubmatch(latest), release.FindStringSubmatch(current)
	if l == nil || c == nil {
		return false
	}
	for i := 1; i <= 3; i++ {
		ln, _ := strconv.Atoi(l[i])
		cn, _ := strconv.Atoi(c[i])
		if ln != cn {
			return ln > cn
		}
	}
	return false
}

// howToUpdate names the way this binary was installed, when that can be told from where
// it lives.
func howToUpdate() string {
	if executable, err := os.Executable(); err == nil {
		if resolved, err := filepath.EvalSymlinks(executable); err == nil {
			executable = resolved
		}
		if strings.Contains(executable, string(filepath.Separator)+"Caskroom"+string(filepath.Separator)) {
			return "Update with: brew upgrade steadybit"
		}
	}
	return "Get it from " + LatestURL
}

func read(file string) state {
	var s state
	if content, err := os.ReadFile(file); err == nil {
		_ = json.Unmarshal(content, &s)
	}
	return s
}

// write replaces the cache in one step, through a file renamed over it, so that a command
// ending mid-write, or two running at once, never leave half a file behind.
func write(file string, s state) {
	content, err := json.Marshal(s)
	if err != nil {
		return
	}
	dir := filepath.Dir(file)
	if os.MkdirAll(dir, 0o755) != nil {
		return
	}
	tmp, err := os.CreateTemp(dir, filepath.Base(file)+".*")
	if err != nil {
		return
	}
	_, err = tmp.Write(content)
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil || os.Rename(tmp.Name(), file) != nil {
		_ = os.Remove(tmp.Name())
	}
}

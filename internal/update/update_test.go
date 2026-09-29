// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package update

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// github answers like the releases page: a redirect to the latest tag.
func github(t *testing.T, tag string) *atomic.Int32 {
	var asked atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		asked.Add(1)
		w.Header().Set("Location", "https://github.com/steadybit/cli/releases/tag/"+tag)
		w.WriteHeader(http.StatusFound)
	}))
	t.Cleanup(server.Close)
	original, originalInteractive, originalNow := LatestURL, Interactive, Now
	LatestURL, Interactive = server.URL, func() bool { return true }
	t.Cleanup(func() { LatestURL, Interactive, Now = original, originalInteractive, originalNow })
	return &asked
}

func notice(current, cache string) string {
	var out bytes.Buffer
	Start(current, cache)(&out)
	return out.String()
}

func TestTellsOfANewerRelease(t *testing.T) {
	github(t, "v6.10.0")
	cache := filepath.Join(t.TempDir(), "update-check.json")

	out := notice("6.9.1", cache)

	assert.Contains(t, out, "A new release of the Steadybit CLI is available: 6.9.1 → 6.10.0\n")
	assert.Contains(t, out, "Get it from ")
}

func TestAsksAtMostOnceADay(t *testing.T) {
	asked := github(t, "v6.1.0")
	cache := filepath.Join(t.TempDir(), "update-check.json")
	now := time.Date(2026, 9, 29, 9, 0, 0, 0, time.UTC)
	Now = func() time.Time { return now }

	first := notice("6.0.1", cache)
	now = now.Add(23 * time.Hour)
	second := notice("6.0.1", cache)
	now = now.Add(2 * time.Hour)
	notice("6.0.1", cache)

	// The second knew of 6.1.0 from the first, without asking.
	assert.Contains(t, first, "6.0.1 → 6.1.0")
	assert.Contains(t, second, "6.0.1 → 6.1.0")
	assert.EqualValues(t, 2, asked.Load())
}

func TestSaysNothingWhenUpToDateOrNotARelease(t *testing.T) {
	asked := github(t, "v6.0.1")
	dir := t.TempDir()

	assert.Empty(t, notice("6.0.1", filepath.Join(dir, "a.json")))
	assert.Empty(t, notice("6.1.0", filepath.Join(dir, "b.json")))
	// go install ...@main and local builds are not releases, and are left alone.
	assert.Empty(t, notice("6.0.0-20260928130623-cbee9e0fadeb", filepath.Join(dir, "c.json")))
	assert.Empty(t, notice("dev", filepath.Join(dir, "d.json")))
	assert.EqualValues(t, 2, asked.Load())
}

func TestStaysQuietWhereNobodyWouldSeeIt(t *testing.T) {
	asked := github(t, "v7.0.0")
	Interactive = func() bool { return false }

	assert.Empty(t, notice("6.0.1", filepath.Join(t.TempDir(), "update-check.json")))
	assert.Zero(t, asked.Load())
}

// Offline, the check is not retried on every command until the next day.
func TestAnUnreachableGitHubIsNotAskedAgainThatDay(t *testing.T) {
	github(t, "v6.1.0")
	LatestURL = "http://127.0.0.1:1/releases/latest"
	cache := filepath.Join(t.TempDir(), "update-check.json")

	assert.Empty(t, notice("6.0.1", cache))
	checked := read(cache)
	assert.False(t, checked.CheckedAt.IsZero())
	assert.Empty(t, checked.Latest)
}

func TestComparesVersionsByNumber(t *testing.T) {
	assert.True(t, newer("6.10.0", "6.9.9"))
	assert.True(t, newer("7.0.0", "6.99.99"))
	assert.False(t, newer("6.0.1", "6.0.1"))
	assert.False(t, newer("6.0.0", "6.0.1"))
	assert.False(t, newer("", "6.0.1"))
}

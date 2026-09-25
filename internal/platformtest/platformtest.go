// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package platformtest runs commands against a fake platform: an httptest server whose
// endpoints a test declares, which records every request it receives.
package platformtest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/steadybit/cli/v6/internal/platform"
)

type Request struct {
	Method string
	Path   string
	Query  map[string][]string
	Header http.Header
	Body   []byte
}

// JSON decodes the request body.
func (r Request) JSON(t *testing.T) any {
	t.Helper()
	var v any
	if err := json.Unmarshal(r.Body, &v); err != nil {
		t.Fatalf("request body is not JSON: %s", r.Body)
	}
	return v
}

type Reply struct {
	Status  int
	JSON    any
	Body    string
	Headers map[string]string
}

type Platform struct {
	t        *testing.T
	server   *httptest.Server
	mu       sync.Mutex
	routes   map[string]func(Request) Reply
	requests []Request
	Client   *platform.Client
	URL      string
}

// New starts a fake platform and points the CLI configuration at it.
func New(t *testing.T) *Platform {
	t.Helper()
	// The fake platform meters nothing, and a suite would soon exhaust the real allowance.
	platform.SetLimiter(platform.NewRateLimiter(platform.Bucket{Burst: 1 << 30, RefillTokens: 1 << 30, RefillInterval: time.Second}, nil))
	p := &Platform{t: t, routes: map[string]func(Request) Reply{}}
	p.server = httptest.NewServer(http.HandlerFunc(p.serve))
	t.Cleanup(p.server.Close)
	p.URL = p.server.URL
	Home(t)
	t.Setenv("STEADYBIT_URL", p.server.URL)
	t.Setenv("STEADYBIT_TOKEN", "test-token")
	client, err := platform.New()
	if err != nil {
		t.Fatal(err)
	}
	p.Client = client
	return p
}

// Handle answers "METHOD /path"; a path segment written as * matches any value.
func (p *Platform) Handle(route string, reply func(Request) Reply) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.routes[route] = reply
}

// Reply answers a route with a fixed reply.
func (p *Platform) Reply(route string, reply Reply) {
	p.Handle(route, func(Request) Reply { return reply })
}

func (p *Platform) Requests(route string) []Request {
	p.mu.Lock()
	defer p.mu.Unlock()
	method, path, _ := strings.Cut(route, " ")
	var matching []Request
	for _, r := range p.requests {
		if r.Method == method && matches(path, r.Path) {
			matching = append(matching, r)
		}
	}
	return matching
}

func matches(pattern, path string) bool {
	a, b := strings.Split(pattern, "/"), strings.Split(path, "/")
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != "*" && a[i] != b[i] {
			return false
		}
	}
	return true
}

func (p *Platform) serve(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	request := Request{Method: r.Method, Path: r.URL.EscapedPath(), Query: r.URL.Query(), Header: r.Header, Body: body}
	p.mu.Lock()
	p.requests = append(p.requests, request)
	// The most specific route answers: `/schedules/v2` over `/schedules/*`. Map order is
	// random, so taking the first match made tests pass or fail by chance.
	var handler func(Request) Reply
	best := -1
	for route, h := range p.routes {
		method, path, _ := strings.Cut(route, " ")
		if method == r.Method && matches(path, request.Path) {
			if exact := strings.Count(path, "/") - strings.Count(path, "*"); exact > best {
				handler, best = h, exact
			}
		}
	}
	p.mu.Unlock()
	if handler == nil {
		p.t.Errorf("unexpected request %s %s", r.Method, request.Path)
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	reply := handler(request)
	for k, v := range reply.Headers {
		w.Header().Set(k, v)
	}
	status := reply.Status
	if status == 0 {
		status = http.StatusOK
	}
	switch {
	case reply.JSON != nil:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(reply.JSON)
	default:
		w.WriteHeader(status)
		_, _ = w.Write([]byte(reply.Body))
	}
}

// Home gives the test an empty home directory. Go reads USERPROFILE for it on Windows
// and HOME elsewhere, so both are set.
func Home(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	return home
}

// Stdout captures what fn prints.
func Stdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()
	original := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()
	runErr := fn()
	_ = w.Close()
	os.Stdout = original
	return <-done, runErr
}

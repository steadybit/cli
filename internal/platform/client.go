// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package platform builds the generated API client with what every request needs:
// authentication, a User-Agent, request logging, and retries that are safe to make.
package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/steadybit/cli/api"
	"github.com/steadybit/cli/internal/config"
	"github.com/steadybit/cli/internal/output"
)

var Version = "dev"

// ErrNoAccessToken is reported with the setup help, before any request is made.
var ErrNoAccessToken = errors.New("no API access token")

func MissingTokenHelp() string {
	return strings.TrimSpace(fmt.Sprintf(`
No API access token configuration was found for Steadybit platform access.
You can configure API access tokens through configuration profiles or
environment variables (%s). We recommend configuration profiles
for local CLI usage. You can add a configuration profile via

                  %s
`, output.Bold("STEADYBIT_TOKEN"), output.Bold("steadybit config profile add")))
}

// APIError carries the status and the problem body of a failed request.
type APIError struct {
	Method, URL string
	Status      int
	Body        []byte
}

func (e *APIError) Error() string {
	body := string(e.Body)
	if body == "" {
		body = "<no body>"
	}
	return fmt.Sprintf("Steadybit API at %s %s responded with unexpected status code: %d - %s", e.Method, e.URL, e.Status, body)
}

// ProblemType is the `type` of an RFC 7807 problem body, if there is one.
func (e *APIError) ProblemType() string {
	var problem struct {
		Type string `json:"type"`
	}
	_ = json.Unmarshal(e.Body, &problem)
	return problem.Type
}

type Client struct {
	*api.ClientWithResponses
	BaseURL   string
	http      *http.Client
	authorize api.RequestEditorFn
}

// Get fetches a path the spec has no operation for, such as the Location of a run.
func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}
	if err := c.authorize(ctx, req); err != nil {
		return nil, err
	}
	return c.http.Do(req)
}

var Verbose bool

func New() (*Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	if cfg.APIAccessToken == "" {
		return nil, ErrNoAccessToken
	}
	base, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL '%s': %w", cfg.BaseURL, err)
	}
	httpClient := &http.Client{
		// No client-wide timeout: it would also count the time a request waits for the
		// rate limiter, which under a dump is far longer than any request takes. The
		// transport bounds each attempt instead.
		Transport: &transport{next: http.DefaultTransport, base: base},
		// Requests carry the access token; following a redirect could hand it to another host.
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	authorize := func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "accessToken "+cfg.APIAccessToken)
		req.Header.Set("Accept", "application/json, */*")
		req.Header.Set("User-Agent", "steadybit@"+Version)
		return nil
	}
	client, err := api.NewClientWithResponses(cfg.BaseURL, api.WithHTTPClient(httpClient), api.WithRequestEditorFn(authorize))
	if err != nil {
		return nil, err
	}
	return &Client{ClientWithResponses: client, BaseURL: cfg.BaseURL, http: httpClient, authorize: authorize}, nil
}

// Check turns any non-2xx response into an APIError.
func Check(resp *http.Response, body []byte) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return &APIError{Method: resp.Request.Method, URL: resp.Request.URL.String(), Status: resp.StatusCode, Body: body}
}

const maxRateLimitWait = 2 * time.Minute

const defaultTimeout = 30 * time.Second

type timeoutKey struct{}

// WithTimeout gives the requests made with ctx a longer deadline than the default 30
// seconds, as artifact downloads need.
func WithTimeout(ctx context.Context, d time.Duration) context.Context {
	return context.WithValue(ctx, timeoutKey{}, d)
}

// cancelOnClose ends an attempt's deadline once its body has been read, so the deadline
// bounds the download of the body and not just the wait for the status line.
type cancelOnClose struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (c cancelOnClose) Close() error {
	defer c.cancel()
	return c.ReadCloser.Close()
}

var idempotent = map[string]bool{"GET": true, "HEAD": true, "OPTIONS": true, "PUT": true, "DELETE": true}

// transport retries what is safe to retry: a 429 for any method, since the request was
// rejected rather than applied, and a transport failure only for idempotent methods, as
// a POST that failed in transit may still have started an experiment run.
type transport struct {
	next http.RoundTripper
	base *url.URL
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Absolute URLs from the platform, such as the Location of a run, keep their path
	// but are sent to the configured origin, so the token never leaves that host.
	req.URL.Scheme, req.URL.Host = t.base.Scheme, t.base.Host

	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
		_ = req.Body.Close()
	}
	var waited time.Duration
	for attempt := 1; ; attempt++ {
		if body != nil {
			req.Body = io.NopCloser(bytes.NewReader(body))
		}
		Limiter().Acquire()
		timeout := defaultTimeout
		if d, ok := req.Context().Value(timeoutKey{}).(time.Duration); ok {
			timeout = d
		}
		ctx, cancel := context.WithTimeout(req.Context(), timeout)
		sent := req.WithContext(ctx)
		logRequest(sent, body)
		resp, err := t.next.RoundTrip(sent)
		if err != nil {
			cancel()
			if !idempotent[req.Method] || attempt >= 4 {
				return nil, fmt.Errorf("Failed to call Steadybit API at %s %s: %w", req.Method, req.URL, err)
			}
			time.Sleep(jitter(time.Duration(attempt) * time.Second))
			continue
		}
		logResponse(resp)
		if resp.StatusCode != http.StatusTooManyRequests {
			resp.Body = cancelOnClose{ReadCloser: resp.Body, cancel: cancel}
			return resp, nil
		}
		wait := time.Second
		for _, h := range []string{"RateLimit-Reset", "Retry-After"} {
			if seconds, err := strconv.Atoi(resp.Header.Get(h)); err == nil && seconds > 0 {
				wait = time.Duration(seconds) * time.Second
				break
			}
		}
		if waited+wait > maxRateLimitWait {
			resp.Body = cancelOnClose{ReadCloser: resp.Body, cancel: cancel}
			return resp, nil
		}
		_ = resp.Body.Close()
		cancel()
		time.Sleep(wait)
		waited += wait
	}
}

func jitter(d time.Duration) time.Duration {
	return d/2 + time.Duration(rand.Int64N(int64(d/2)+1))
}

func logRequest(req *http.Request, body []byte) {
	if !Verbose {
		return
	}
	fmt.Printf("> HTTP %s %s\n", req.Method, req.URL)
	for name, values := range req.Header {
		value := strings.Join(values, ", ")
		if strings.EqualFold(name, "Authorization") || strings.EqualFold(name, "Cookie") {
			value = "<redacted>"
		}
		fmt.Printf("> %s: %s\n", name, value)
	}
	fmt.Println(">")
	if len(body) > 0 {
		fmt.Println(string(body))
	}
	fmt.Println()
}

func logResponse(resp *http.Response) {
	if !Verbose {
		return
	}
	fmt.Fprintf(os.Stdout, "< HTTP %s\n\n", resp.Status)
}

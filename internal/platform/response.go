// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package platform

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/steadybit/cli/internal/output"
)

// Read takes a generated client call's result and returns the body, failing on any
// status outside 2xx.
func Read(resp *http.Response, err error) ([]byte, *http.Response, error) {
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}
	return body, resp, Check(resp, body)
}

// Decode reads a response into target.
func Decode(resp *http.Response, err error, target any) (*http.Response, error) {
	body, resp, err := Read(resp, err)
	if err != nil {
		return resp, err
	}
	return resp, json.Unmarshal(body, target)
}

// ReadDocument reads a response as an order-preserving document.
func ReadDocument(resp *http.Response, err error) (*output.Document, *http.Response, error) {
	body, resp, err := Read(resp, err)
	if err != nil {
		return nil, resp, err
	}
	doc, err := output.ParseDocument(body)
	return doc, resp, err
}

func IsStatus(err error, status int) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.Status == status
}

// Failed reports a failed request as the TypeScript CLI did: the message, the request
// error, and the platform's problem body pretty-printed, which is what names the
// violated constraint.
func Failed(err error, format string, args ...any) error {
	message := err.Error()
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		if problem := apiErr.problemJSON(); problem != "" {
			message += ": " + problem
		}
	}
	return fmt.Errorf("%s: %s", fmt.Sprintf(format, args...), message)
}

func (e *APIError) problemJSON() string {
	body := bytes.TrimSpace(e.Body)
	if len(body) == 0 || !json.Valid(body) {
		return ""
	}
	switch body[0] {
	case '{':
		doc, err := output.ParseDocument(body)
		if err != nil {
			return ""
		}
		rendered, _ := doc.Render(output.JSON)
		return string(bytes.TrimSuffix(rendered, []byte("\n")))
	case '[':
		if string(body) == "[]" {
			return "[]"
		}
		// Arrays are rare in problem bodies; reindenting them as JavaScript would is
		// left to json.Indent, which agrees for plain values.
		var buf bytes.Buffer
		if json.Indent(&buf, body, "", "  ") == nil {
			return buf.String()
		}
	case '"':
		var s string
		if json.Unmarshal(body, &s) == nil && s != "" {
			return string(body)
		}
	}
	return ""
}

// PageSize is the most items a paged endpoint returns at once.
const PageSize int32 = 100

// AllPages follows nextPage until the last page: a listing cut at the first response
// would silently leave the rest out.
func AllPages[T any](fetch func(page, size int32) (*http.Response, error)) ([]T, error) {
	var items []T
	page := int32(0)
	for {
		var body struct {
			Items    []T    `json:"items"`
			NextPage *int32 `json:"nextPage"`
		}
		resp, err := fetch(page, PageSize)
		if _, err := Decode(resp, err, &body); err != nil {
			return nil, err
		}
		items = append(items, body.Items...)
		if body.NextPage == nil {
			return items, nil
		}
		page = *body.NextPage
	}
}

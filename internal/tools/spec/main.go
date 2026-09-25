// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Command spec keeps the committed platform spec and the client generated from it in
// step with the platform. The CLI's requests and responses are compiled against that
// client, so a breaking change in the platform fails the build instead of a customer's
// pipeline.
//
//	go run ./internal/tools/spec fetch   download the live spec into openapi/platform-api.json
//	go generate ./api                    regenerate api/platform.gen.go from the committed spec
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

const specFile = "openapi/platform-api.json"

// Examples in the spec end up in the generated doc comments, and one of them is a Slack
// incoming webhook URL, which GitHub push protection rightly refuses. Examples carry no
// type information, so they are masked before anything is written.
var slackWebhook = regexp.MustCompile(`https://hooks\.slack\.com/services/[^"\s\\]+`)

func main() {
	if len(os.Args) != 2 || os.Args[1] != "fetch" {
		fmt.Fprintln(os.Stderr, "usage: go run ./internal/tools/spec fetch")
		os.Exit(2)
	}
	if err := fetch(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func fetch() error {
	url := os.Getenv("STEADYBIT_SPEC_URL")
	if url == "" {
		url = "https://platform.steadybit.com/api/spec"
	}
	client := &http.Client{Timeout: time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("fetching the platform spec from %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetching the platform spec from %s failed with status %d", url, resp.StatusCode)
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var spec struct {
		OpenAPI string                     `json:"openapi"`
		Paths   map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(raw, &spec); err != nil {
		return fmt.Errorf("%s did not return JSON: %w", url, err)
	}
	if spec.OpenAPI == "" || spec.Paths == nil {
		return fmt.Errorf("%s did not return an OpenAPI document", url)
	}
	// The platform's own bytes, indented: re-encoding would reorder sections and escape
	// every < and > in the descriptions, and a diff of the file should read as a diff of
	// the API.
	var out bytes.Buffer
	if err := json.Indent(&out, raw, "", "  "); err != nil {
		return err
	}
	out.WriteByte('\n')
	if err := os.WriteFile(specFile, slackWebhook.ReplaceAll(out.Bytes(), []byte("https://hooks.slack.com/services/<redacted>")), 0o644); err != nil {
		return err
	}
	paths := spec.Paths
	fmt.Printf("Wrote %d paths from %s to %s\n", len(paths), url, specFile)
	return nil
}

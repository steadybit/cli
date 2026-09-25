// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Command npmpkg turns goreleaser's binaries into the npm packages: one per platform
// holding the binary, and `steadybit`, which depends on all of them optionally and starts
// the one npm installed. npm 12 no longer runs install scripts by default, so downloading
// the binary on install is not an option.
//
//	go run ./internal/tools/npmpkg -version 6.0.0 -dist dist -out npm/dist
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type platform struct {
	goos, goarch, nodeOS, nodeCPU string
}

var platforms = []platform{
	{"darwin", "arm64", "darwin", "arm64"},
	{"darwin", "amd64", "darwin", "x64"},
	{"linux", "arm64", "linux", "arm64"},
	{"linux", "amd64", "linux", "x64"},
	{"windows", "arm64", "win32", "arm64"},
	{"windows", "amd64", "win32", "x64"},
}

type artifact struct {
	Type   string `json:"type"`
	Path   string `json:"path"`
	Goos   string `json:"goos"`
	Goarch string `json:"goarch"`
}

func main() {
	version := flag.String("version", "", "the version to publish, without a leading v")
	dist := flag.String("dist", "dist", "goreleaser's output directory")
	out := flag.String("out", "npm/dist", "where to write the packages")
	flag.Parse()
	if err := run(strings.TrimPrefix(*version, "v"), *dist, *out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(version, dist, out string) error {
	if version == "" {
		return fmt.Errorf("-version is required")
	}
	content, err := os.ReadFile(filepath.Join(dist, "artifacts.json"))
	if err != nil {
		return err
	}
	var artifacts []artifact
	if err := json.Unmarshal(content, &artifacts); err != nil {
		return err
	}
	binaries := map[string]string{}
	for _, a := range artifacts {
		if a.Type == "Binary" {
			binaries[a.Goos+"/"+a.Goarch] = a.Path
		}
	}

	if err := os.RemoveAll(out); err != nil {
		return err
	}
	optional := map[string]string{}
	for _, p := range platforms {
		binary, ok := binaries[p.goos+"/"+p.goarch]
		if !ok {
			return fmt.Errorf("no %s/%s binary in %s", p.goos, p.goarch, dist)
		}
		name := fmt.Sprintf("@steadybit/cli-%s-%s", p.nodeOS, p.nodeCPU)
		dir := filepath.Join(out, fmt.Sprintf("cli-%s-%s", p.nodeOS, p.nodeCPU))
		executable := "steadybit"
		if p.goos == "windows" {
			executable += ".exe"
		}
		if err := copyFile(binary, filepath.Join(dir, "bin", executable), 0o755); err != nil {
			return err
		}
		if err := writeJSON(filepath.Join(dir, "package.json"), map[string]any{
			"name":            name,
			"version":         version,
			"description":     fmt.Sprintf("The Steadybit CLI binary for %s %s. Install `steadybit` instead.", p.nodeOS, p.nodeCPU),
			"license":         "MIT",
			"repository":      map[string]string{"type": "git", "url": "https://github.com/steadybit/cli.git"},
			"os":              []string{p.nodeOS},
			"cpu":             []string{p.nodeCPU},
			"files":           []string{"bin"},
			"preferUnplugged": true,
		}); err != nil {
			return err
		}
		optional[name] = version
	}

	main := filepath.Join(out, "steadybit")
	if err := copyFile("npm/steadybit/bin/steadybit.js", filepath.Join(main, "bin", "steadybit.js"), 0o755); err != nil {
		return err
	}
	for _, f := range []string{"README.md", "LICENSE", "CHANGELOG.md"} {
		if err := copyFile(f, filepath.Join(main, f), 0o644); err != nil {
			return err
		}
	}
	return writeJSON(filepath.Join(main, "package.json"), mainPackage{
		Name:                 "steadybit",
		Version:              version,
		Description:          "Command-line interface to interact with the Steadybit API",
		Keywords:             []string{"steadybit", "cli", "chaos engineering", "resilience engineering", "api", "gitops"},
		License:              "MIT",
		Author:               "Steadybit GmbH",
		Repository:           repo,
		Bin:                  map[string]string{"steadybit": "bin/steadybit.js"},
		Files:                []string{"bin", "README.md", "CHANGELOG.md", "LICENSE"},
		Engines:              map[string]string{"node": ">=18"},
		OptionalDependencies: optional,
	})
}

func copyFile(from, to string, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
		return err
	}
	src, err := os.Open(from)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}

// writeJSON writes a package.json without Go's HTML escaping, which would turn ">=" in
// engines into "\u003e=".
func writeJSON(file string, value any) error {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		return err
	}
	return os.WriteFile(file, buf.Bytes(), 0o644)
}

type repository struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

var repo = repository{Type: "git", URL: "https://github.com/steadybit/cli.git"}

type platformPackage struct {
	Name            string     `json:"name"`
	Version         string     `json:"version"`
	Description     string     `json:"description"`
	License         string     `json:"license"`
	Repository      repository `json:"repository"`
	OS              []string   `json:"os"`
	CPU             []string   `json:"cpu"`
	Files           []string   `json:"files"`
	PreferUnplugged bool       `json:"preferUnplugged"`
}

type mainPackage struct {
	Name                 string            `json:"name"`
	Version              string            `json:"version"`
	Description          string            `json:"description"`
	Keywords             []string          `json:"keywords"`
	License              string            `json:"license"`
	Author               string            `json:"author"`
	Repository           repository        `json:"repository"`
	Bin                  map[string]string `json:"bin"`
	Files                []string          `json:"files"`
	Engines              map[string]string `json:"engines"`
	OptionalDependencies map[string]string `json:"optionalDependencies"`
}

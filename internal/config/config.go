// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package config reads and writes the CLI configuration exactly as the TypeScript CLI
// did, so that an upgrade keeps every profile a user has already set up.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const DefaultBaseURL = "https://platform.steadybit.com"

type Profile struct {
	Name           string `json:"name"`
	APIAccessToken string `json:"apiAccessToken"`
	BaseURL        string `json:"baseUrl,omitempty"`
}

type Configuration struct {
	APIAccessToken string
	BaseURL        string
}

func dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".steadybit"), nil
}

func file(name string) (string, error) {
	d, err := dir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(d, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(d, name), nil
}

func Profiles() ([]Profile, error) {
	path, err := file("profiles.json")
	if err != nil {
		return nil, err
	}
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []Profile{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%s': %w", path, err)
	}
	var profiles []Profile
	if err := json.Unmarshal(content, &profiles); err != nil {
		return nil, fmt.Errorf("failed to parse file '%s' as JSON: %w", path, err)
	}
	return profiles, nil
}

func writeProfiles(profiles []Profile) error {
	path, err := file("profiles.json")
	if err != nil {
		return err
	}
	content, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		return err
	}
	// The file holds access tokens.
	return os.WriteFile(path, content, 0o600)
}

func AddProfile(profile Profile) error {
	profiles, err := Profiles()
	if err != nil {
		return err
	}
	kept := profiles[:0]
	for _, p := range profiles {
		if p.Name != profile.Name {
			kept = append(kept, p)
		}
	}
	return writeProfiles(append(kept, profile))
}

func RemoveProfile(name string) error {
	profiles, err := Profiles()
	if err != nil {
		return err
	}
	kept := profiles[:0]
	for _, p := range profiles {
		if p.Name != name {
			kept = append(kept, p)
		}
	}
	return writeProfiles(kept)
}

func activeProfileName() (string, error) {
	path, err := file("activeProfile")
	if err != nil {
		return "", err
	}
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to read file '%s': %w", path, err)
	}
	return strings.TrimSpace(string(content)), nil
}

// ActiveProfile is the selected profile, or the first one when none is selected.
func ActiveProfile() (*Profile, error) {
	profiles, err := Profiles()
	if err != nil || len(profiles) == 0 {
		return nil, err
	}
	name, err := activeProfileName()
	if err != nil {
		return nil, err
	}
	for i := range profiles {
		if profiles[i].Name == name {
			return &profiles[i], nil
		}
	}
	return &profiles[0], nil
}

func SetActiveProfile(name string) error {
	path, err := file("activeProfile")
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(name), 0o644)
}

// Load resolves the configuration: environment variables win over the active profile.
// An empty STEADYBIT_TOKEN counts as set, as it did in the TypeScript CLI, so that a
// pipeline can deliberately blank it out.
func Load() (Configuration, error) {
	cfg := Configuration{BaseURL: DefaultBaseURL}
	profile, err := ActiveProfile()
	if err != nil {
		return cfg, err
	}
	if profile != nil {
		cfg.APIAccessToken = profile.APIAccessToken
		if profile.BaseURL != "" {
			cfg.BaseURL = profile.BaseURL
		}
	}
	if token, ok := os.LookupEnv("STEADYBIT_TOKEN"); ok {
		cfg.APIAccessToken = token
	}
	if url, ok := os.LookupEnv("STEADYBIT_URL"); ok {
		cfg.BaseURL = url
	}
	cfg.BaseURL = strings.TrimSuffix(cfg.BaseURL, "/")
	return cfg, nil
}

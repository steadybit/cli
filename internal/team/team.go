// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package team implements the `team` commands.
package team

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/steadybit/cli/v6/api"
	"github.com/steadybit/cli/v6/internal/jsyaml"
	"github.com/steadybit/cli/v6/internal/output"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/resource"
	"github.com/steadybit/cli/v6/internal/table"
)

// The version is dropped as `service get` drops it. Members are sent back as the platform
// takes them, by username, email and role; the rest describes the user.
var (
	readOnly       = []string{"version"}
	memberReadOnly = []string{"name", "pictureUrl", "managedBy"}
)

func notFoundOr(err error, key, format string) error {
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("Team %s not found.", key)
	}
	return platform.Failed(err, format, key)
}

type ListOptions struct {
	Search string
	Type   string
}

func List(ctx context.Context, c *platform.Client, o ListOptions) error {
	var summaries struct {
		Teams []struct {
			Key, Name           string
			AllowedEnvironments []string `json:"allowedEnvironments"`
			Members             []any    `json:"members"`
		} `json:"teams"`
	}
	params := &api.GetTeamsParams{}
	if o.Search != "" {
		params.Search = &o.Search
	}
	resp, err := c.GetTeams(ctx, params)
	raw, err := resource.DecodeListed(resp, err, "teams", &summaries)
	if err != nil {
		return platform.Failed(err, "Failed to get the teams")
	}
	if resource.Machine(o.Type) {
		return resource.List(raw, o.Type, nil)
	}
	if len(summaries.Teams) == 0 {
		fmt.Println("No teams found.")
		return nil
	}
	t := table.New(
		table.Column{Name: "key", Title: "Key", Alignment: table.Left},
		table.Column{Name: "name", Title: "Name", Alignment: table.Left},
		table.Column{Name: "members", Title: "Members"},
		table.Column{Name: "environments", Title: "Environments"},
	)
	for _, team := range summaries.Teams {
		t.AddRow(table.Default, table.Cell("key", team.Key), table.Cell("name", team.Name), table.Cell("members", len(team.Members)),
			table.Cell("environments", len(team.AllowedEnvironments)))
	}
	t.Print()
	return nil
}

type GetOptions struct {
	Key, File, Type string
}

func Get(ctx context.Context, c *platform.Client, o GetOptions) error {
	doc, _, err := platform.ReadDocument(c.GetTeam(ctx, o.Key))
	if err != nil {
		return notFoundOr(err, o.Key, "Failed to get team %s")
	}
	resource.Strip(doc, readOnly...)
	if members, ok := doc.Value().Get("members"); ok {
		list, _ := members.([]any)
		for _, m := range list {
			if member, ok := m.(*jsyaml.Map); ok {
				for _, field := range memberReadOnly {
					member.Delete(field)
				}
			}
		}
	}
	if err := resource.Output(doc, o.File, o.Type); err != nil {
		return err
	}
	if o.File != "" {
		fmt.Printf("Team %s written to %s.\n", o.Key, o.File)
	}
	return nil
}

type ApplyOptions struct {
	Files     []string
	Recursive bool
}

// Apply upserts teams by their key, which is what names a team; there is no id to write back.
func Apply(ctx context.Context, c *platform.Client, o ApplyOptions) error {
	return resource.ApplyFiles(o.Files, o.Recursive, "team", func(file string, doc *output.Document) (resource.Applied, error) {
		key, _ := doc.Get("key")
		if key == "" {
			return resource.Applied{}, fmt.Errorf("Team file '%s' does not name the team key.", file)
		}
		resp, err := c.UpsertTeamWithBody(ctx, &api.UpsertTeamParams{}, "application/json", resource.Body(resource.Strip(doc, readOnly...).Value()))
		_, resp, err = platform.Read(resp, err)
		if err != nil {
			return resource.Applied{}, platform.Failed(err, "Failed to save team %s", key)
		}
		created := resp.StatusCode == http.StatusCreated
		fmt.Printf("Team %s %s.\n", key, resource.CreatedOrUpdated(created))
		return resource.Applied{Created: created}, nil
	})
}

type DeleteOptions struct {
	Key         string
	Experiments bool
	Yes         bool
}

func Delete(ctx context.Context, c *platform.Client, o DeleteOptions) error {
	question := fmt.Sprintf("Delete team %s?", o.Key)
	if o.Experiments {
		question = fmt.Sprintf("Delete team %s with all its experiments and their runs?", o.Key)
	}
	if ok, err := resource.Confirmed(o.Yes, question); !ok || err != nil {
		return err
	}
	_, _, err := platform.Read(c.DeleteTeam(ctx, o.Key, &api.DeleteTeamParams{PurgeIncludingExperiments: o.Experiments}))
	// The platform refuses without purging, even a team without experiments, and says
	// nothing about why.
	if !o.Experiments && platform.IsStatus(err, http.StatusBadRequest) {
		return fmt.Errorf("Team %s was not deleted. The platform deletes a team only with --purge-experiments, which deletes its experiments and their runs too.", o.Key)
	}
	if err != nil {
		return notFoundOr(err, o.Key, "Failed to delete team %s")
	}
	fmt.Printf("Team %s deleted.\n", o.Key)
	return nil
}

type member struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type members struct {
	Members []member `json:"members"`
}

func ListMembers(ctx context.Context, c *platform.Client, key, explicitType string) error {
	var result members
	resp, err := c.GetTeamMembers(ctx, key)
	raw, err := resource.DecodeListed(resp, err, "members", &result)
	if err != nil {
		return notFoundOr(err, key, "Failed to get the members of team %s")
	}
	if resource.Machine(explicitType) {
		return resource.List(raw, explicitType, nil)
	}
	printMembers(key, result.Members)
	return nil
}

func printMembers(key string, list []member) {
	if len(list) == 0 {
		fmt.Printf("Team %s has no members.\n", key)
		return
	}
	t := table.New(
		table.Column{Name: "username", Title: "Username", Alignment: table.Left},
		table.Column{Name: "name", Title: "Name", Alignment: table.Left},
		table.Column{Name: "email", Title: "Email", Alignment: table.Left},
		table.Column{Name: "role", Title: "Role", Alignment: table.Left},
	)
	for _, m := range list {
		t.AddRow(table.Default, table.Cell("username", m.Username), table.Cell("name", m.Name), table.Cell("email", m.Email), table.Cell("role", m.Role))
	}
	t.Print()
}

// MemberOptions names users by username or email; the platform takes either.
type MemberOptions struct {
	Key       string
	Usernames []string
	Emails    []string
	Role      string
	Validate  bool
	Yes       bool
}

func (o MemberOptions) update() (api.TeamMembersUpdateAO, error) {
	role := api.MemberUpdateAORole(strings.ToUpper(o.Role))
	if role != api.MemberUpdateAORoleMEMBER && role != api.MemberUpdateAORoleOWNER {
		return api.TeamMembersUpdateAO{}, fmt.Errorf("--role must be MEMBER or OWNER, not '%s'.", o.Role)
	}
	update := api.TeamMembersUpdateAO{Members: []api.MemberUpdateAO{}}
	for _, u := range o.Usernames {
		update.Members = append(update.Members, api.MemberUpdateAO{Username: &u, Role: role})
	}
	for _, e := range o.Emails {
		update.Members = append(update.Members, api.MemberUpdateAO{Email: &e, Role: role})
	}
	return update, nil
}

func (o MemberOptions) count() int { return len(o.Usernames) + len(o.Emails) }

var errNoMembers = errors.New("No members given. Pass --username or --email.")

func AddMembers(ctx context.Context, c *platform.Client, o MemberOptions) error {
	if o.count() == 0 {
		return errNoMembers
	}
	update, err := o.update()
	if err != nil {
		return err
	}
	var result members
	resp, err := c.AddTeamMembers(ctx, o.Key, &api.AddTeamMembersParams{ValidateMembers: &o.Validate}, update)
	if _, err := platform.Decode(resp, err, &result); err != nil {
		return notFoundOr(err, o.Key, "Failed to add members to team %s")
	}
	fmt.Printf("Team %s now has %d member(s).\n", o.Key, len(result.Members))
	return nil
}

func RemoveMembers(ctx context.Context, c *platform.Client, o MemberOptions) error {
	if o.count() == 0 {
		return errNoMembers
	}
	if ok, err := resource.Confirmed(o.Yes, fmt.Sprintf("Remove %d member(s) from team %s?", o.count(), o.Key)); !ok || err != nil {
		return err
	}
	request := api.TeamMembersRemoveAO{Usernames: resource.Optional(o.Usernames), Emails: resource.Optional(o.Emails)}
	var result members
	resp, err := c.RemoveTeamMembers(ctx, o.Key, request)
	if _, err := platform.Decode(resp, err, &result); err != nil {
		return notFoundOr(err, o.Key, "Failed to remove members from team %s")
	}
	fmt.Printf("Team %s now has %d member(s).\n", o.Key, len(result.Members))
	return nil
}

// SetMembers replaces the members: everyone not given is removed from the team.
func SetMembers(ctx context.Context, c *platform.Client, o MemberOptions) error {
	if o.count() == 0 {
		return errNoMembers
	}
	update, err := o.update()
	if err != nil {
		return err
	}
	if ok, err := resource.Confirmed(o.Yes, fmt.Sprintf("Make these %d the only members of team %s, removing everyone else?", o.count(), o.Key)); !ok || err != nil {
		return err
	}
	var result members
	resp, err := c.SetTeamMembers(ctx, o.Key, &api.SetTeamMembersParams{ValidateMembers: &o.Validate}, update)
	if _, err := platform.Decode(resp, err, &result); err != nil {
		return notFoundOr(err, o.Key, "Failed to set the members of team %s")
	}
	fmt.Printf("Team %s now has %d member(s).\n", o.Key, len(result.Members))
	return nil
}

type environments struct {
	Environments []struct {
		Name string `json:"name"`
	} `json:"environments"`
}

func (e environments) names() []string {
	names := make([]string, len(e.Environments))
	for i, env := range e.Environments {
		names[i] = env.Name
	}
	return names
}

func ListEnvironments(ctx context.Context, c *platform.Client, key, explicitType string) error {
	var result environments
	resp, err := c.GetTeamEnvironments(ctx, key)
	raw, err := resource.DecodeListed(resp, err, "environments", &result)
	if err != nil {
		return notFoundOr(err, key, "Failed to get the environments of team %s")
	}
	if resource.Machine(explicitType) {
		return resource.List(raw, explicitType, nil)
	}
	if len(result.Environments) == 0 {
		fmt.Printf("Team %s has no environments.\n", key)
		return nil
	}
	t := table.New(table.Column{Name: "name", Title: "Environment", Alignment: table.Left})
	for _, name := range result.names() {
		t.AddRow(table.Default, table.Cell("name", name))
	}
	t.Print()
	return nil
}

type EnvironmentOptions struct {
	Key          string
	Environments []string
	Validate     bool
	Yes          bool
}

func (o EnvironmentOptions) list() []api.TeamEnvironmentAO {
	list := make([]api.TeamEnvironmentAO, len(o.Environments))
	for i, name := range o.Environments {
		list[i] = api.TeamEnvironmentAO{Name: name}
	}
	return list
}

var errNoEnvironments = errors.New("No environments given. Pass --environment.")

func reportEnvironments(key string, result environments) {
	if len(result.Environments) == 0 {
		fmt.Printf("Team %s now has no environments.\n", key)
		return
	}
	fmt.Printf("Team %s now has %d environment(s): %s.\n", key, len(result.Environments), strings.Join(result.names(), ", "))
}

func AddEnvironments(ctx context.Context, c *platform.Client, o EnvironmentOptions) error {
	if len(o.Environments) == 0 {
		return errNoEnvironments
	}
	var result environments
	resp, err := c.AddTeamEnvironments(ctx, o.Key, &api.AddTeamEnvironmentsParams{ValidateEnvironments: &o.Validate}, api.TeamEnvironmentsUpdateAO{Environments: o.list()})
	if _, err := platform.Decode(resp, err, &result); err != nil {
		return notFoundOr(err, o.Key, "Failed to add environments to team %s")
	}
	reportEnvironments(o.Key, result)
	return nil
}

func RemoveEnvironments(ctx context.Context, c *platform.Client, o EnvironmentOptions) error {
	if len(o.Environments) == 0 {
		return errNoEnvironments
	}
	var result environments
	resp, err := c.RemoveTeamEnvironments(ctx, o.Key, api.TeamEnvironmentsUpdateAO{Environments: o.list()})
	if _, err := platform.Decode(resp, err, &result); err != nil {
		return notFoundOr(err, o.Key, "Failed to remove environments from team %s")
	}
	reportEnvironments(o.Key, result)
	return nil
}

// SetEnvironments replaces the environments a team may use.
func SetEnvironments(ctx context.Context, c *platform.Client, o EnvironmentOptions) error {
	if len(o.Environments) == 0 {
		return errNoEnvironments
	}
	if ok, err := resource.Confirmed(o.Yes, fmt.Sprintf("Make these %d the only environments of team %s?", len(o.Environments), o.Key)); !ok || err != nil {
		return err
	}
	var result environments
	resp, err := c.SetTeamEnvironments(ctx, o.Key, &api.SetTeamEnvironmentsParams{ValidateEnvironments: &o.Validate}, api.TeamEnvironmentsAO{Environments: o.list()})
	if _, err := platform.Decode(resp, err, &result); err != nil {
		return notFoundOr(err, o.Key, "Failed to set the environments of team %s")
	}
	reportEnvironments(o.Key, result)
	return nil
}

// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package user implements the `user` commands.
package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/steadybit/cli/v6/api"
	"github.com/steadybit/cli/v6/internal/platform"
	"github.com/steadybit/cli/v6/internal/resource"
	"github.com/steadybit/cli/v6/internal/table"
)

type InviteOptions struct {
	Emails []string
	Role   string
	Team   string
}

func Invite(ctx context.Context, c *platform.Client, o InviteOptions) error {
	if len(o.Emails) == 0 {
		return errors.New("No one to invite. Pass --email.")
	}
	var role *api.InvitationAORole
	if o.Role != "" {
		r := api.InvitationAORole(strings.ToUpper(o.Role))
		if r != api.InvitationAORoleADMIN && r != api.InvitationAORoleUSER {
			return fmt.Errorf("--role must be USER or ADMIN, not '%s'.", o.Role)
		}
		role = &r
	}
	request := api.InviteUsersRequestAO{Invitations: make([]api.InvitationAO, len(o.Emails))}
	for i, email := range o.Emails {
		if _, err := openapi_types.Email(email).MarshalJSON(); err != nil {
			return fmt.Errorf("'%s' is not an email address.", email)
		}
		request.Invitations[i] = api.InvitationAO{Email: openapi_types.Email(email), Role: role}
		if o.Team != "" {
			request.Invitations[i].TeamKey = &o.Team
		}
	}
	_, _, err := platform.Read(c.InviteUser(ctx, request))
	if platform.IsStatus(err, http.StatusForbidden) {
		return errors.New("Inviting users needs an admin access token.")
	}
	if err != nil {
		return platform.Failed(err, "Failed to invite %s", strings.Join(o.Emails, ", "))
	}
	fmt.Printf("%d user(s) invited. They receive an email to join.\n", len(o.Emails))
	return nil
}

type ListOptions struct {
	Query string
	Type  string
}

func List(ctx context.Context, c *platform.Client, o ListOptions) error {
	params := api.GetUsersParams{}
	if o.Query != "" {
		params.Query = &o.Query
	}
	raw, err := platform.AllPagesRaw(func(page, size int32) (*http.Response, error) {
		params.Page = api.PageRequestAO{Page: &page, Size: &size}
		return c.GetUsers(ctx, &params)
	})
	if platform.IsStatus(err, http.StatusForbidden) {
		return errors.New("Listing users needs an admin access token.")
	}
	if err != nil {
		return platform.Failed(err, "Failed to get the users")
	}
	return resource.List(raw, o.Type, func() error {
		var users []struct {
			Username  string   `json:"username"`
			Name      string   `json:"name"`
			Email     string   `json:"email"`
			Role      string   `json:"role"`
			Teams     []string `json:"teams"`
			ManagedBy string   `json:"managedBy"`
			LastLogin *string  `json:"lastLogin"`
		}
		if err := resource.DecodeEach(raw, &users); err != nil {
			return err
		}
		if len(users) == 0 {
			fmt.Println("No users found.")
			return nil
		}
		t := table.New(
			table.Column{Name: "username", Title: "Username", Alignment: table.Left},
			table.Column{Name: "name", Title: "Name", Alignment: table.Left},
			table.Column{Name: "email", Title: "Email", Alignment: table.Left},
			table.Column{Name: "role", Title: "Role", Alignment: table.Left},
			table.Column{Name: "teams", Title: "Teams", Alignment: table.Left},
			table.Column{Name: "managedBy", Title: "Managed by", Alignment: table.Left},
			table.Column{Name: "lastLogin", Title: "Last login", Alignment: table.Left},
		)
		for _, u := range users {
			managedBy, lastLogin := u.ManagedBy, "never"
			if managedBy == "" {
				managedBy = "Steadybit"
			}
			if u.LastLogin != nil && *u.LastLogin != "" {
				lastLogin = *u.LastLogin
			}
			t.AddRow(table.Default, table.Cell("username", u.Username), table.Cell("name", u.Name), table.Cell("email", u.Email),
				table.Cell("role", u.Role), table.Cell("teams", strings.Join(u.Teams, ", ")), table.Cell("managedBy", managedBy), table.Cell("lastLogin", lastLogin))
		}
		t.Print()
		return nil
	})
}

type DeleteOptions struct {
	Username string
	Yes      bool
}

// Delete removes a user from the tenant. Users managed by SCIM or an LDAP sync are
// removed in the identity provider instead; the platform refuses them.
func Delete(ctx context.Context, c *platform.Client, o DeleteOptions) error {
	if ok, err := resource.Confirmed(o.Yes, fmt.Sprintf("Remove user %s from the tenant? They lose access to it.", o.Username)); !ok || err != nil {
		return err
	}
	// as a query parameter, not a path segment: usernames may contain '/' (ADFS subjects are base64)
	_, _, err := platform.Read(c.RemoveUserFromTenant(ctx, &api.RemoveUserFromTenantParams{Username: o.Username}))
	if platform.IsStatus(err, http.StatusNotFound) {
		return fmt.Errorf("User %s is not a member of the tenant.", o.Username)
	}
	if err != nil {
		return platform.Failed(err, "Failed to remove user %s", o.Username)
	}
	fmt.Printf("User %s removed from the tenant.\n", o.Username)
	return nil
}

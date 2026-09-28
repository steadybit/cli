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

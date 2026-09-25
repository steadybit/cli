// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package killswitch implements the `killswitch` commands.
package killswitch

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/steadybit/cli/internal/output"
	"github.com/steadybit/cli/internal/platform"
	"github.com/steadybit/cli/internal/resource"
)

type StatusOptions struct {
	Type string
}

func Status(ctx context.Context, c *platform.Client, o StatusOptions) error {
	body, _, err := platform.Read(c.GetKillswitch(ctx))
	if err != nil {
		return platform.Failed(err, "Failed to get the kill switch status")
	}
	if o.Type != "" {
		doc, err := output.ParseDocument(body)
		if err != nil {
			return err
		}
		return resource.Output(doc, "", o.Type)
	}
	var status struct {
		Active           bool   `json:"active"`
		EngagedBy        string `json:"engagedBy"`
		Engaged          string `json:"engaged"`
		EngagedByDetails *struct {
			Name string `json:"name"`
		} `json:"engagedByDetails"`
	}
	if err := json.Unmarshal(body, &status); err != nil {
		return err
	}
	if !status.Active {
		fmt.Println("The kill switch is inactive: experiments can run.")
		return nil
	}
	by := status.EngagedBy
	if status.EngagedByDetails != nil && status.EngagedByDetails.Name != "" {
		by = status.EngagedByDetails.Name
	}
	fmt.Printf("The kill switch is active since %s, activated by %s: no experiment can run.\n", status.Engaged, by)
	return nil
}

// Activate stops every experiment running in the tenant, of every team, and keeps new
// ones from starting until the kill switch is deactivated.
func Activate(ctx context.Context, c *platform.Client, yes bool) error {
	if ok, err := resource.Confirmed(yes, "Activate the kill switch? It stops every running experiment of every team, and no experiment can run until it is deactivated."); !ok || err != nil {
		return err
	}
	if _, _, err := platform.Read(c.EngageKillswitch(ctx)); err != nil {
		return platform.Failed(err, "Failed to activate the kill switch")
	}
	fmt.Println("Kill switch activated. Running experiments are stopped, and none can run until `steadybit killswitch deactivate`.")
	return nil
}

func Deactivate(ctx context.Context, c *platform.Client) error {
	if _, _, err := platform.Read(c.DisengageKillswitch(ctx)); err != nil {
		return platform.Failed(err, "Failed to deactivate the kill switch")
	}
	fmt.Println("Kill switch deactivated. Experiments can run again; those it stopped are not restarted.")
	return nil
}

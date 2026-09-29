// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package gitops_test

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/steadybit/cli/v6/internal/gitops"
	"github.com/steadybit/cli/v6/internal/platformtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	environmentID = "019eacd7-fb2c-733a-bed5-99a935323d01"
	teamID        = "019eacd7-fb2c-733a-bed5-99a935323d02"
	hubID         = "019eacd7-fb2c-733a-bed5-99a935323d03"
	templateID    = "019eacd7-fb2c-733a-bed5-99a935323d04"
	importedID    = "019eacd7-fb2c-733a-bed5-99a935323d05"
	webhookID     = "019eacd7-fb2c-733a-bed5-99a935323d06"
	slackID       = "019eacd7-fb2c-733a-bed5-99a935323d07"
	providedHubID = "6a55640d-72d4-4194-a058-0afcf731dfac"
)

const (
	storedWebhook = `{"id":"` + webhookID + `","version":0,"scope":"GLOBAL","name":"Notify","url":"https://example.com/hook","secret":"****************","events":["*"],"targetAttributeIncludes":["*"],"headers":{"Authorization":"Bearer abc"}}`
	storedSlack   = `{"id":"` + slackID + `","version":1,"scope":"GLOBAL","name":"Chat","url":"https://hooks.example.com/T1/B2/xyz","events":["*"],"channel":"#chaos"}`
)

func fakeTenant(t *testing.T) *platformtest.Platform {
	p := platformtest.New(t)
	p.Reply("GET /api/properties/definitions", platformtest.Reply{JSON: map[string]any{"items": []any{map[string]any{"key": "tribe"}}}})
	p.Reply("GET /api/properties/definitions/tribe", platformtest.Reply{Body: `{"key":"tribe","label":"Tribe","dataType":"STRING","version":3}`})
	p.Reply("GET /api/environments", platformtest.Reply{JSON: map[string]any{"environments": []any{map[string]any{"id": environmentID}}}})
	p.Reply("GET /api/environments/"+environmentID, platformtest.Reply{Body: `{"id":"` + environmentID + `","name":"Prod","version":2,"predicate":{"operator":"AND","predicates":[]},"state":"READY"}`})
	p.Reply("GET /api/teams", platformtest.Reply{JSON: map[string]any{"teams": []any{map[string]any{"key": "ADM"}}}})
	p.Reply("GET /api/teams/ADM", platformtest.Reply{Body: `{"id":"` + teamID + `","key":"ADM","name":"Admins","allowedActions":["wait","service-validation"],"allowedEnvironments":["Prod"],"managedBy":"MANUAL","version":0,` +
		`"members":[{"username":"u1","name":"Jane","email":"jane@example.com","role":"OWNER","pictureUrl":"https://example.com/jane.png","managedBy":"MANUAL"}]}`})
	p.Reply("GET /api/hubs", platformtest.Reply{JSON: map[string]any{"hubs": []any{map[string]any{"id": providedHubID}, map[string]any{"id": hubID}}}})
	p.Reply("GET /api/hubs/"+providedHubID, platformtest.Reply{Body: `{"hubName":"Steadybit Reliability Hub","id":"` + providedHubID + `","templates":[{"id":"` + importedID + `","templateTitle":"Imported"}]}`})
	p.Reply("GET /api/hubs/"+hubID, platformtest.Reply{Body: `{"hubName":"Own Hub","repositoryUrl":"https://example.com/index.json","id":"` + hubID + `","version":4,"templates":[],"lastSync":"x","created":"c"}`})
	p.Reply("GET /api/experiments/templates", platformtest.Reply{JSON: map[string]any{"templates": []any{map[string]any{"id": templateID}, map[string]any{"id": importedID}}}})
	p.Reply("GET /api/experiments/templates/"+templateID, platformtest.Reply{Body: `{"id":"` + templateID + `","templateTitle":"Pod Crash","templateDescription":"d","placeholders":[],"tags":[],"lanes":[{"steps":[{"type":"wait","ignoreFailure":false,"parameters":{"duration":"5s"}}]}],"properties":{},"propertiesMetadata":[],"version":0,"created":"c","createdBy":{}}`})
	p.Reply("GET /api/integrations/webhook", platformtest.Reply{JSON: map[string]any{"content": []any{map[string]any{"id": webhookID}}}})
	p.Reply("GET /api/integrations/webhook/"+webhookID, platformtest.Reply{Body: storedWebhook})
	p.Reply("GET /api/integrations/slack", platformtest.Reply{JSON: map[string]any{"content": []any{map[string]any{"id": slackID}}}})
	p.Reply("GET /api/integrations/slack/"+slackID, platformtest.Reply{Body: storedSlack})
	p.Reply("GET /api/integrations/preflight", platformtest.Reply{JSON: map[string]any{"content": []any{}}})
	p.Reply("GET /api/integrations/preflight-action", platformtest.Reply{JSON: map[string]any{"content": []any{}}})
	p.Reply("GET /api/services/profiles", platformtest.Reply{JSON: map[string]any{"items": []any{
		map[string]any{"id": profileID, "name": "Shop", "origin": "CUSTOM"},
		map[string]any{"id": serviceID, "name": "Kubernetes Deployment", "origin": "STEADYBIT"},
	}}})
	p.Reply("GET /api/services/profiles/"+profileID, platformtest.Reply{Body: `{"id":"` + profileID + `","name":"Shop","origin":"CUSTOM","templates":[],"defaultProfile":false,"version":1}`})
	return p
}

func TestExportTenantLeavesOutWhatThePlatformProvides(t *testing.T) {
	p := fakeTenant(t)
	dir := t.TempDir()

	out, err := platformtest.Stdout(t, func() error { return gitops.Export(ctx, p.Client, gitops.ExportOptions{Directory: dir, Tenant: true}) })

	require.NoError(t, err)
	assert.Equal(t, "Exported the tenant to "+dir+": 1 experiment templates, 1 environments, 1 teams, 1 property definitions, 1 hubs, 2 integrations, 1 service profiles.\n", out)
	for file, content := range map[string]string{
		"property-definitions/tribe.yaml": "key: tribe\nlabel: Tribe\ndataType: STRING\n",
		"environments/prod.yaml":          "id: " + environmentID + "\nname: Prod\npredicate:\n  operator: AND\n  predicates: []\n",
		"teams/adm.yaml": "id: " + teamID + "\nkey: ADM\nname: Admins\nallowedActions:\n  - wait\n  - service-validation\nallowedEnvironments:\n  - Prod\nmanagedBy: MANUAL\n" +
			"members:\n  - username: u1\n    email: jane@example.com\n    role: OWNER\n",
		"hubs/own-hub.yaml": "hubName: Own Hub\nrepositoryUrl: https://example.com/index.json\nid: " + hubID + "\n",
		"templates/pod-crash.yaml": "id: " + templateID + "\ntemplateTitle: Pod Crash\ntemplateDescription: d\nplaceholders: []\ntags: []\nlanes:\n  - steps:\n      - type: wait\n        ignoreFailure: false\n" +
			"        parameters:\n          duration: 5s\nproperties: {}\npropertiesMetadata: []\n",
		// Credentials are masked: the platform's mask of the secret gives away its length.
		"integrations/webhook/notify.yaml": "id: " + webhookID + "\nscope: GLOBAL\nname: Notify\nurl: https://example.com/hook\nsecret: '********'\nevents:\n  - '*'\n" +
			"targetAttributeIncludes:\n  - '*'\nheaders:\n  Authorization: '********'\n",
		"integrations/slack/chat.yaml": "id: " + slackID + "\nscope: GLOBAL\nname: Chat\nurl: '********'\nevents:\n  - '*'\nchannel: '#chaos'\n",
		"service-profiles/shop.yaml":   "id: " + profileID + "\nname: Shop\norigin: CUSTOM\ntemplates: []\n",
	} {
		written, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(file)))
		require.NoError(t, err, file)
		assert.Equal(t, content, string(written), file)
	}
	for _, gone := range []string{"hubs/steadybit-reliability-hub.yaml", "templates/imported.yaml", "service-profiles/kubernetes-deployment.yaml", "integrations/preflight"} {
		assert.NoFileExists(t, filepath.Join(dir, filepath.FromSlash(gone)))
	}
	assert.Empty(t, p.Requests("GET /api/experiments/templates/"+importedID), "a template imported from a hub is not even fetched")

	out, err = platformtest.Stdout(t, func() error { return gitops.DiffProject(ctx, p.Client, dir) })
	require.NoError(t, err)
	assert.Equal(t, 8, strings.Count(out, "match the platform"), out)
}

// Applying an exported tenant leaves out the integrations that match: the masked secret
// in them could not be sent back, and leaving it out would remove it.
func TestApplyingAnExportedTenantSkipsMatchingIntegrations(t *testing.T) {
	p := fakeTenant(t)
	dir := t.TempDir()
	_, err := platformtest.Stdout(t, func() error { return gitops.Export(ctx, p.Client, gitops.ExportOptions{Directory: dir, Tenant: true}) })
	require.NoError(t, err)
	for _, route := range []string{"POST /api/properties/definitions", "POST /api/environments", "POST /api/teams", "POST /api/hubs", "POST /api/experiments/templates", "POST /api/services/profiles"} {
		p.Reply(route, platformtest.Reply{JSON: map[string]any{}})
	}

	out, err := platformtest.Stdout(t, func() error { return gitops.ApplyProject(ctx, p.Client, gitops.ApplyOptions{Directory: dir}) })

	require.NoError(t, err)
	assert.Contains(t, out, "notify.yaml matches webhook integration "+webhookID+".\n")
	assert.Contains(t, out, "chat.yaml matches Slack integration "+slackID+".\n")
	assert.Empty(t, p.Requests("POST /api/integrations/webhook"))
	assert.Empty(t, p.Requests("POST /api/integrations/slack"))
	order := []string{"Property definition tribe", "Environment", "Team ADM", "Hub", "Experiment template", "Service profile"}
	last := -1
	for _, line := range order {
		i := strings.Index(out, line)
		assert.Greater(t, i, last, "%s is applied after what it depends on:\n%s", line, out)
		last = i
	}
}

// A changed integration is applied with the credentials the platform reads back in the
// clear; its secret it never does, so that has to be put in.
func TestAChangedExportedIntegrationKeepsItsCredentials(t *testing.T) {
	p := fakeTenant(t)
	dir := t.TempDir()
	_, err := platformtest.Stdout(t, func() error { return gitops.Export(ctx, p.Client, gitops.ExportOptions{Directory: dir, Tenant: true}) })
	require.NoError(t, err)
	for _, other := range []string{"property-definitions", "environments", "teams", "hubs", "templates", "service-profiles", "integrations/webhook"} {
		require.NoError(t, os.RemoveAll(filepath.Join(dir, filepath.FromSlash(other))))
	}
	slack := filepath.Join(dir, "integrations", "slack", "chat.yaml")
	require.NoError(t, os.WriteFile(slack, []byte("id: "+slackID+"\nscope: GLOBAL\nname: Chat\nurl: '********'\nevents:\n  - '*'\nchannel: '#incidents'\n"), 0o644))
	p.Reply("POST /api/integrations/slack", platformtest.Reply{JSON: map[string]any{"id": slackID, "name": "Chat"}})

	out, err := platformtest.Stdout(t, func() error {
		return gitops.ApplyProject(ctx, p.Client, gitops.ApplyOptions{Directory: dir})
	})

	require.NoError(t, err)
	assert.Equal(t, "Slack integration Chat ("+slackID+") updated.\n", out)
	sent := p.Requests("POST /api/integrations/slack")[0].JSON(t).(map[string]any)
	assert.Equal(t, "https://hooks.example.com/T1/B2/xyz", sent["url"])
	assert.Equal(t, "#incidents", sent["channel"])
	written, _ := os.ReadFile(slack)
	assert.Contains(t, string(written), "url: '********'\n", "the credential is not written to the file")
}

func TestCredentialsAreNeverPrinted(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/integrations/webhook/"+webhookID, platformtest.Reply{Body: storedWebhook})
	file := write(t, "w.yml", "id: "+webhookID+"\nscope: GLOBAL\nname: Notify\nurl: https://example.com/hook\nsecret: n3w-s3cret\nevents:\n  - '*'\n"+
		"targetAttributeIncludes:\n  - '*'\nheaders:\n  Authorization: Bearer other\n")

	out, err := platformtest.Stdout(t, func() error {
		return gitops.DiffFiles(ctx, p.Client, gitops.Integrations["webhook"], []string{file}, false)
	})

	assert.ErrorIs(t, err, gitops.ErrDifferent)
	assert.Contains(t, out, "-secret: '********'\n+secret: '******** (from the file)'\n")
	assert.Contains(t, out, "-  Authorization: '********'\n+  Authorization: '******** (from the file)'\n")
	assert.NotContains(t, out, "s3cret")
	assert.NotContains(t, out, "Bearer")
}

// A team sent no actions may still wait and validate services, and a webhook sent no
// target attributes reports all of them. Neither is a difference.
func TestAHandWrittenTeamAndWebhookMatchRightAfterApply(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/teams/CLIX", platformtest.Reply{Body: `{"id":"` + teamID + `","key":"CLIX","name":"x","allowedActions":["wait","service-validation"],"allowedEnvironments":["Prod"],"managedBy":"MANUAL","version":0,"members":[]}`})
	p.Reply("GET /api/integrations/webhook/"+webhookID, platformtest.Reply{Body: `{"id":"` + webhookID + `","version":0,"scope":"GLOBAL","name":"n","url":"https://example.com","events":["*"],"targetAttributeIncludes":["*"],"headers":{}}`})
	team := write(t, "team.yml", "key: CLIX\nname: x\nallowedActions: []\nallowedEnvironments:\n  - Prod\n")
	webhook := write(t, "webhook.yml", "id: "+webhookID+"\nscope: GLOBAL\nname: n\nurl: https://example.com\nevents:\n  - '*'\ntargetAttributeIncludes: []\n")

	_, err := platformtest.Stdout(t, func() error { return gitops.DiffFiles(ctx, p.Client, gitops.Team, []string{team}, false) })
	require.NoError(t, err)
	_, err = platformtest.Stdout(t, func() error {
		return gitops.DiffFiles(ctx, p.Client, gitops.Integrations["webhook"], []string{webhook}, false)
	})
	require.NoError(t, err)

	p.Reply("GET /api/teams/CLIX", platformtest.Reply{Body: `{"key":"CLIX","name":"x","allowedActions":["wait"],"allowedEnvironments":["Prod"]}`})
	out, err := platformtest.Stdout(t, func() error { return gitops.DiffFiles(ctx, p.Client, gitops.Team, []string{team}, false) })
	assert.ErrorIs(t, err, gitops.ErrDifferent)
	assert.Contains(t, out, "-allowedActions:\n-  - wait\n+allowedActions: []\n")
}

func TestEnvironmentsAreFoundByName(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/environments", platformtest.Reply{JSON: map[string]any{"environments": []any{
		map[string]any{"id": teamID, "name": "Prod EU"},
		map[string]any{"id": environmentID, "name": "Prod"},
	}}})
	p.Reply("GET /api/environments/"+environmentID, platformtest.Reply{Body: `{"id":"` + environmentID + `","name":"Prod","version":2,"predicate":{"operator":"AND","predicates":[]},"state":"READY"}`})
	file := write(t, "e.yml", "name: Prod\npredicate:\n  operator: AND\n  predicates: []\n")

	out, err := platformtest.Stdout(t, func() error { return gitops.DryRun(ctx, p.Client, gitops.Environment, []string{file}, false) })

	require.NoError(t, err)
	assert.Equal(t, file+" matches environment "+environmentID+".\n", out)
	assert.Equal(t, []string{"Prod"}, p.Requests("GET /api/environments")[0].Query["search"])
}

func TestAnUnknownTeamWouldBeCreated(t *testing.T) {
	p := platformtest.New(t)
	p.Reply("GET /api/teams/NEW", platformtest.Reply{Status: http.StatusNotFound})
	file := write(t, "t.yml", "key: NEW\nname: New\nallowedActions: []\nallowedEnvironments: []\n")

	out, err := platformtest.Stdout(t, func() error { return gitops.DryRun(ctx, p.Client, gitops.Team, []string{file}, false) })

	require.NoError(t, err)
	assert.Equal(t, file+" would create a new team.\n", out)
}

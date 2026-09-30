#!/usr/bin/env bash
# SPDX-License-Identifier: MIT
# SPDX-FileCopyrightText: 2026 Steadybit GmbH

# Uses the CLI against a real platform, the way a pipeline does: experiments applied from
# files, run in parallel with a report, checked for drift, canceled by the SIGTERM a
# canceled job sends, refused while another runs, and found again by a gate. Everything
# only waits, belongs to one test team, and is deleted afterwards, also when a check fails.
#
# Then the rest of the platform the CLI covers: templates, run properties, schedules,
# services and their profiles with the team token, and environments, teams, access
# tokens, hubs, integrations and the audit log with an admin token. What an admin creates
# is named cli-e2e-ci-*, belongs to a team of its own when it has a scope, and targets
# nothing. The kill switch is only read, and no invitation is sent: both would reach
# beyond the test.
#
# Needs STEADYBIT_TOKEN (a team token) and STEADYBIT_URL, and `steadybit` on the PATH.
# STEADYBIT_E2E_TEAM and STEADYBIT_E2E_ENVIRONMENT name the team and its environment.
# STEADYBIT_E2E_ADMIN_TOKEN, an admin token, enables the checks that need one.
#
# Other test suites run experiments on the same platform at the same time, so every run
# here allows running in parallel, except the one whose refusal is the point of its check.

set -uo pipefail

: "${STEADYBIT_TOKEN:?a team access token is needed}"
TEAM=${STEADYBIT_E2E_TEAM:-CLI}
ENVIRONMENT=${STEADYBIT_E2E_ENVIRONMENT:-Global}
# Every experiment of this suite has a name starting with this, which is how a later run
# finds and deletes what a crashed one left behind.
MARK=cli-e2e-ci
RUN=${GITHUB_RUN_ID:-local-$$}
ADMIN_TOKEN=${STEADYBIT_E2E_ADMIN_TOKEN:-}
# The team the admin checks create; one run at a time, so its key can stay the same.
SCOPE_TEAM=CLIX
# A property of runs, made by the admin checks for the experiment they create.
PROPERTY=cliE2eCi

work=$(mktemp -d)
cd "$work" || exit 1
failures=0

check() {
  description=$1
  shift
  if "$@"; then
    echo "  ok    $description"
  else
    echo "  FAIL  $description"
    failures=$((failures + 1))
  fi
}

# Runs a command and checks its exit status; on a mismatch, shows what it printed.
exits_with() {
  expected=$1
  shift
  "$@" >out.log 2>&1
  actual=$?
  [ "$actual" -eq "$expected" ] || {
    echo "        expected exit $expected, got $actual from: $*"
    tail -n 20 out.log | sed 's/^/        /'
    return 1
  }
}

experiment() { # file name duration
  mkdir -p "$(dirname "$1")"
  cat >"$1" <<EOF
# Written by e2e/platform.sh; only waits.
name: $MARK-$RUN-$2
externalId: $MARK-$RUN-$2
team: $TEAM
environment: $ENVIRONMENT
lanes:
  - steps:
      - type: wait
        ignoreFailure: false
        parameters:
          duration: $3
EOF
}

key_of() { sed -n 's/^key: //p' "$1"; }
id_of() { sed -n 's/^id: //p' "$1"; }

# The CLI with the admin token.
admin() { STEADYBIT_TOKEN=$ADMIN_TOKEN steadybit "$@"; }

# Runs a command and checks that its output, or the file it wrote, has a line.
prints() { # expected command...
  expected=$1
  shift
  "$@" >out.log 2>&1 || { echo "        exit $? from: $*"; tail -n 20 out.log | sed 's/^/        /'; return 1; }
  grep -qF -- "$expected" out.log || { echo "        no '$expected' from: $*"; tail -n 20 out.log | sed 's/^/        /'; return 1; }
}

# Waits until the experiment has a run in one of the states, for up to a minute.
until_run_is() { # key state...
  key=$1
  shift
  for _ in $(seq 30); do
    [ "$(steadybit execution list --key "$key" --state "$@" --limit 1 --jq length 2>/dev/null)" = "1" ] && return 0
    sleep 2
  done
  return 1
}

# Deletes every experiment of this suite in the team, this run's and any a crashed run
# left, after canceling what still runs.
cleanup() {
  echo "cleanup"
  rm -rf sweep
  steadybit export --team "$TEAM" -d sweep >/dev/null 2>&1 || true
  for file in $(grep -l "^name: $MARK-" sweep/experiments/*.yaml 2>/dev/null); do
    key=$(key_of "$file")
    [ -n "$key" ] || continue
    for run in $(steadybit execution list --key "$key" --state CREATED PREPARED RUNNING --limit 0 --jq '.[].id' 2>/dev/null); do
      steadybit execution cancel -i "$run" >/dev/null 2>&1
    done
    for _ in $(seq 10); do
      steadybit experiment delete -k "$key" >/dev/null 2>&1 && break
      sleep 3
    done
    steadybit experiment get -k "$key" >/dev/null 2>&1 && echo "  could not delete $key" || echo "  deleted $key"
  done
  [ -n "$ADMIN_TOKEN" ] && sweep_admin
  rm -rf "$work"
}

# Deletes what the admin checks create, this run's and any a crashed run left: each kind
# is listed and only what is named cli-e2e-ci-* is deleted.
sweep_admin() {
  named() { admin "$@" --jq ".[] | select((.name // .hubName // .templateTitle // \"\") | startswith(\"$MARK-\")) | .id" 2>/dev/null; }
  for kind in webhook preflight preflight-action; do
    for id in $(named integration $kind list -t json); do admin integration $kind delete -i "$id" >/dev/null 2>&1 && echo "  deleted $kind $id"; done
  done
  for id in $(named hub list -t json); do admin hub delete -i "$id" --yes >/dev/null 2>&1 && echo "  deleted hub $id"; done
  for id in $(named access-token list --output json); do admin access-token delete -i "$id" --yes >/dev/null 2>&1 && echo "  deleted access token $id"; done
  for id in $(named service list -t json); do admin service delete -i "$id" >/dev/null 2>&1 && echo "  deleted service $id"; done
  for id in $(named service-profile list -t json); do admin service-profile delete -i "$id" >/dev/null 2>&1 && echo "  deleted service profile $id"; done
  for id in $(admin property association list --key "$PROPERTY" --output json --jq '.[].id' 2>/dev/null); do
    admin property association delete -i "$id" --delete-values --yes >/dev/null 2>&1 && echo "  deleted property association $id"
  done
  admin property definition delete -k "$PROPERTY" --yes >/dev/null 2>&1 && echo "  deleted property definition $PROPERTY"
  for id in $(named template list -t json); do admin template delete -i "$id" --yes >/dev/null 2>&1 && echo "  deleted template $id"; done
  if admin team get -k "$SCOPE_TEAM" >/dev/null 2>&1; then
    admin team delete -k "$SCOPE_TEAM" --purge-experiments --yes >/dev/null 2>&1 && echo "  deleted team $SCOPE_TEAM"
  fi
  for id in $(named environment list -t json); do admin environment delete -i "$id" --yes >/dev/null 2>&1 && echo "  deleted environment $id"; done
}
trap cleanup EXIT

echo "steadybit $(steadybit --version) against ${STEADYBIT_URL:-https://platform.steadybit.com}, team $TEAM"

# a runs long enough for a poll to see it running, which one check expects.
experiment experiments/a.yml a 15s
experiment experiments/b.yml b 8s
experiment long.yml long 90s

check "apply creates the experiments and writes their keys into the files" sh -c '
  steadybit experiment apply -f experiments -f long.yml >/dev/null &&
  grep -q "^key: " experiments/a.yml && grep -q "^key: " experiments/b.yml && grep -q "^key: " long.yml &&
  grep -q "^# Written by" experiments/a.yml
'
A=$(key_of experiments/a.yml)
LONG=$(key_of long.yml)

check "diff finds no drift right after apply" exits_with 0 steadybit experiment diff -f experiments
sed -i.bak 's/duration: 8s/duration: 9s/' experiments/b.yml && rm -f experiments/b.yml.bak
check "diff exits with 2 once a file changed" exits_with 2 steadybit experiment diff -f experiments
check "apply updates the experiment" exits_with 0 steadybit experiment apply -f experiments

check "a parallel run completes both experiments and reports them" sh -c '
  steadybit experiment run -f experiments --yes --parallel 2 --allowParallel --report report.xml >run.log 2>&1 || { tail -n 20 run.log; exit 1; }
  [ "$(grep -c "<testsuite " report.xml)" -eq 2 ] && grep -q "\[" run.log
'

# A canceled job sends SIGTERM: the run it started is canceled, and the CLI exits with 143.
steadybit experiment run -k "$LONG" --yes --allowParallel >term.log 2>&1 &
pid=$!
if until_run_is "$LONG" RUNNING; then
  kill -TERM "$pid"
  wait "$pid"
  status=$?
  check "SIGTERM exits with 143" test "$status" -eq 143
  check "SIGTERM cancels the run" until_run_is "$LONG" CANCELED
else
  kill -TERM "$pid" 2>/dev/null
  check "the long run started" false
fi

# While one experiment runs, the platform accepts another and cancels it right away;
# --no-wait has to notice instead of passing. This run alone does not allow running in
# parallel: its refusal is what is checked.
check "a run started with --no-wait begins" exits_with 0 steadybit experiment run -k "$LONG" --yes --no-wait --allowParallel
until_run_is "$LONG" RUNNING
check "--no-wait fails on a run the platform refused" exits_with 1 steadybit experiment run -k "$A" --yes --no-wait
check "execution list --fail-on-match gates on that canceled run" exits_with 1 \
  steadybit execution list --key "$A" --state CANCELED --ended-from "$(date -u +%F)" --limit 1 --fail-on-match
# What the run-experiment action relies on: the experiment found by its external id, an
# expected state reached before the end, and waiting while another experiment runs.
check "the experiment is found by its external id and passes at the expected state" exits_with 0 \
  steadybit experiment run --external-id "$MARK-$RUN-a" --yes --allowParallel --expect-state RUNNING --report expect.json
check "the report has the state reached and the run's API location" sh -c '
  grep -Eq "\"state\": *\"RUNNING\"" expect.json && grep -q "\"apiLocation\"" expect.json || { cat expect.json; exit 1; }
'
until_run_is "$A" COMPLETED CANCELED
check "a run that ends otherwise than expected fails" sh -c "
  steadybit experiment run -k $A --yes --allowParallel --expect-state FAILED >otherwise.log 2>&1
  status=\$?
  grep -q 'but failed was expected' otherwise.log && [ \$status -eq 1 ] || { tail -n 5 otherwise.log; exit 1; }
"
# The long run still goes, so both tries are refused: what is checked is that the CLI tries
# again instead of failing at once or, as --yes would otherwise do, running in parallel.
# Waiting for the platform to be free would depend on what other suites run at the time.
check "--busy-retries tries again while another experiment runs" sh -c "
  steadybit experiment run -k $A --yes --busy-retries 1 --busy-retry-interval 5s >busy.log 2>&1
  status=\$?
  grep -q 'trying again in 5s (1/1)' busy.log && [ \$status -eq 1 ]
"
check "execution list prints the platform's runs as JSON" sh -c "
  [ \"\$(steadybit execution list --team $TEAM --limit 2 --jq length 2>/dev/null)\" -ge 1 ]
"

# --- Templates, run properties, schedules, services -------------------------------------
#
# Tokens only ever go through the environment, never into a command line, so that a
# failing check, which shows its command, cannot show one.

check "target stats and actions read the platform" sh -c '
  steadybit target stats -t json >/dev/null && steadybit action list --kind BASIC -t json >/dev/null
'

# The CLI with a token read from a file an access-token command wrote.
with_token() { # file command...
  file=$1
  shift
  STEADYBIT_TOKEN=$(token_of "$file") "$@"
}
token_of() { sed -n 's/^ *"token": "\(.*\)",*$/\1/p' "$1"; }
json_id() { sed -n 's/^ *"id": "*\([^",]*\)"*,*$/\1/p' "$1" | head -n 1; }
json_key() { sed -n 's/^ *"key": "\(.*\)",*$/\1/p' "$1" | head -n 1; }
# GitHub hides a masked value wherever it would show in the job log.
mask() { if [ -n "${GITHUB_ACTIONS:-}" ] && [ -n "$1" ]; then echo "::add-mask::$1"; fi; }

applied_with_id() { # file command...
  file=$1
  shift
  "$@" -f "$file" >out.log 2>&1 && grep -q "^id: " "$file" && grep -q "^# Written by" "$file" || { tail -n 20 out.log | sed 's/^/        /'; return 1; }
}
# What get writes is what a repository keeps; diff has to find it unchanged.
round_trips() { # get-command... -- diff-command...
  get=()
  while [ "$1" != -- ]; do get+=("$1"); shift; done
  shift
  "${get[@]}" >out.log 2>&1 && "$@" >>out.log 2>&1 || { tail -n 20 out.log | sed 's/^/        /'; return 1; }
}
schedule_enabled() { steadybit schedule list --experiment "$FROM_TEMPLATE" --jq '.[0].enabled' 2>/dev/null; }

if [ -n "$ADMIN_TOKEN" ]; then
  cat >template.yml <<EOF
# Written by e2e/platform.sh; only waits.
templateTitle: $MARK-$RUN-template
templateDescription: Written by the CLI's platform test.
placeholders:
  - key: DURATION
    name: Duration
    description: How long to wait
tags: []
lanes:
  - steps:
      - type: wait
        ignoreFailure: false
        parameters:
          duration: '[[DURATION]]'
EOF
  check "template apply creates the template and writes its id into the file" applied_with_id template.yml admin template apply
  TEMPLATE=$(id_of template.yml)
  check "template diff finds no drift right after apply" exits_with 0 admin template diff -f template.yml
  sed -i.bak 's/^templateDescription: .*/templateDescription: Changed by the test./' template.yml && rm -f template.yml.bak
  check "template diff exits with 2 once the file changed" exits_with 2 admin template diff -f template.yml
  check "template apply updates the template" exits_with 0 admin template apply -f template.yml
  check "template get --placeholders lists what to fill in" prints "DURATION:" admin template get -i "$TEMPLATE" --placeholders

  # The team runs what an admin provides, which is how templates reach pipelines.
  check "the team runs an experiment from the template" exits_with 0 \
    steadybit experiment run --template "$TEMPLATE" --team "$TEAM" --environment "$ENVIRONMENT" -p DURATION=3s \
    --external-id "$MARK-$RUN-template" --yes --allowParallel --report template.json
  FROM_TEMPLATE=$(json_key template.json)
  FROM_TEMPLATE_RUN=$(json_id template.json)
  check "the experiment from the template has the placeholder filled in" prints "duration: 3s" steadybit experiment get -k "$FROM_TEMPLATE"

  cat >property.yml <<EOF
key: $PROPERTY
label: $MARK property
dataType: STRING_LIST
EOF
  cat >association.yml <<EOF
key: $PROPERTY
associationType: EXPERIMENT
experimentKey: $FROM_TEMPLATE
editableInExecution: true
required: false
EOF
  check "an admin defines a list property" exits_with 0 admin property definition apply -f property.yml
  check "an admin gives it to that experiment's runs only" exits_with 0 admin property association apply -f association.yml
  check "the team sets the list property with one value" exits_with 0 steadybit execution property set -i "$FROM_TEMPLATE_RUN" -k "$PROPERTY" --value one
  check "the team adds a value to it" exits_with 0 steadybit execution property add -i "$FROM_TEMPLATE_RUN" -k "$PROPERTY" --value two
  check "the run has both values" prints '["one","two"]' steadybit execution get -i "$FROM_TEMPLATE_RUN" --jq ".properties.$PROPERTY | tojson"
  check "execution artifact list lists a run's artifacts" exits_with 0 steadybit execution artifact list -i "$FROM_TEMPLATE_RUN" -t json

  # Far in the future, and enabled only for a moment: it never runs.
  check "schedule create adds a disabled schedule" exits_with 0 \
    steadybit schedule create -k "$FROM_TEMPLATE" --start-at 2099-01-01T00:00:00Z --disabled
  SCHEDULE=$(steadybit schedule list --experiment "$FROM_TEMPLATE" --jq '.[0].id' 2>/dev/null)
  check "schedule get and diff find no drift" round_trips steadybit schedule get -i "$SCHEDULE" -f schedule.yml -- steadybit schedule diff -f schedule.yml
  check "schedule update changes it" exits_with 0 steadybit schedule update -i "$SCHEDULE" --cron "0 0 0 1 1 ? 2099" --timezone Europe/Berlin
  check "schedule diff exits with 2 after the change" exits_with 2 steadybit schedule diff -f schedule.yml
  check "schedule apply puts the file back" exits_with 0 steadybit schedule apply -f schedule.yml
  steadybit schedule enable -i "$SCHEDULE" >/dev/null 2>&1
  check "schedule enable enables it" test "$(schedule_enabled)" = true
  steadybit schedule disable -i "$SCHEDULE" >/dev/null 2>&1
  check "schedule disable disables it" test "$(schedule_enabled)" = false
  check "schedule delete deletes it" exits_with 0 steadybit schedule delete -i "$SCHEDULE"

  cat >profile.yml <<EOF
# Written by e2e/platform.sh.
name: $MARK-$RUN-profile
templates:
  - category: $MARK
    templateIds:
      - $TEMPLATE
EOF
  check "service-profile apply creates a profile with the template" applied_with_id profile.yml admin service-profile apply
  check "service-profile diff finds no drift" exits_with 0 admin service-profile diff -f profile.yml
  # Its query matches nothing, and its validation is never run.
  cat >service.yml <<EOF
# Written by e2e/platform.sh; targets nothing.
name: $MARK-$RUN-service
environment: $ENVIRONMENT
team: $TEAM
serviceProfile: $MARK-$RUN-profile
query: target.type="cli-e2e.nothing"
validations:
  - type: action
    actionType: com.steadybit.extension_http.check.periodically
    ignoreFailure: false
    parameters:
      url: https://example.com
      method: GET
      headers: []
      statusCode: 200-299
      successRate: 100
      maxConcurrent: 1
      requestsPerSecond: 1
      responseTime: 500ms
      responseTimeMode: NO_VERIFICATION
      responseTimeMeasurement: TIME_TO_FIRST_BYTE
      readTimeout: 5s
      connectTimeout: 5s
      followRedirects: false
    radius:
      targetType: com.steadybit.extension_http.client-location
      maximum: 1
EOF
  check "the team creates a service from a file" applied_with_id service.yml steadybit service apply
  SERVICE=$(id_of service.yml)
  check "service diff finds no drift" exits_with 0 steadybit service diff -f service.yml
  check "the service offers the profile's template" prints "$TEMPLATE" steadybit service experiment list -i "$SERVICE"
  check "service experiment link adds the experiment" exits_with 0 steadybit service experiment link -i "$SERVICE" -k "$FROM_TEMPLATE" --category "$MARK"
  check "service risk scores the service" prints "Risk of service" steadybit service risk -i "$SERVICE"
  steadybit service variable set -i "$SERVICE" cliE2e=one >/dev/null 2>&1
  check "service variable set and get" prints "cliE2e: one" steadybit service variable get -i "$SERVICE"
  check "service experiment unlink removes the experiment" exits_with 0 steadybit service experiment unlink -i "$SERVICE" -k "$FROM_TEMPLATE"
  check "service delete deletes it" exits_with 0 steadybit service delete -i "$SERVICE"
  check "service-profile delete deletes it" exits_with 0 admin service-profile delete -i "$(id_of profile.yml)"

  # --- Environments, teams, access tokens, hubs, integrations, audit log ----------------

  cat >environment.yml <<EOF
# Written by e2e/platform.sh; holds no targets.
name: $MARK-$RUN-environment
predicate:
  query: target.type="cli-e2e.nothing"
EOF
  check "environment apply creates the environment" applied_with_id environment.yml admin environment apply
  ENV_ID=$(id_of environment.yml)
  check "environment get and diff find no drift" round_trips admin environment get -i "$ENV_ID" -f environment.yml -- admin environment diff -f environment.yml
  admin environment variable set -i "$ENV_ID" cliE2e=one >/dev/null 2>&1
  check "environment variable set and get" prints "cliE2e: one" admin environment variable get -i "$ENV_ID"

  cat >team.yml <<EOF
# Written by e2e/platform.sh.
key: $SCOPE_TEAM
name: $MARK-$RUN-team
allowedActions:
  - wait
allowedEnvironments:
  - $ENVIRONMENT
EOF
  check "team apply creates the team" exits_with 0 admin team apply -f team.yml
  check "team get and diff find no drift" round_trips admin team get -k "$SCOPE_TEAM" -f team.yml -- admin team diff -f team.yml
  check "team environment add gives it the new environment" prints "$MARK-$RUN-environment" \
    admin team environment add -k "$SCOPE_TEAM" --environment "$MARK-$RUN-environment"
  # An address no one has: without --validate it is skipped, with it the command fails.
  check "team member add skips an unknown user" prints "0 member(s)" admin team member add -k "$SCOPE_TEAM" --email nobody@cli-e2e.invalid
  check "team member add --validate fails on an unknown user" exits_with 1 admin team member add -k "$SCOPE_TEAM" --email nobody@cli-e2e.invalid --validate

  # A token of the new team, which expires tomorrow at the latest.
  tomorrow=$(date -u -d tomorrow +%F 2>/dev/null || date -u -v+1d +%F)
  admin access-token create --name "$MARK-$RUN-token" --type TEAM --team "$SCOPE_TEAM" --expires-at "$tomorrow" -t json >token.json 2>out.log
  mask "$(token_of token.json)"
  check "access-token create makes a team token" test -n "$(token_of token.json)"
  check "the new token works" prints "$SCOPE_TEAM" with_token token.json steadybit team get -k "$SCOPE_TEAM" --jq .key
  admin access-token recreate -i "$(json_id token.json)" --yes -t json >recreated.json 2>out.log
  mask "$(token_of recreated.json)"
  check "access-token recreate makes a new token" test -n "$(token_of recreated.json)"
  check "the replaced token stops working" exits_with 1 with_token token.json steadybit team get -k "$SCOPE_TEAM"
  check "the recreated token works" prints "$SCOPE_TEAM" with_token recreated.json steadybit team get -k "$SCOPE_TEAM" --jq .key
  check "access-token delete deletes it" exits_with 0 admin access-token delete -i "$(json_id recreated.json)" --yes
  rm -f token.json recreated.json

  # A second entry for the public hub; nothing is imported from it.
  cat >hub.yml <<EOF
# Written by e2e/platform.sh; nothing is imported from it.
hubName: $MARK-$RUN-hub
hubLink: https://hub.steadybit.com
repositoryUrl: https://raw.githubusercontent.com/steadybit/reliability-hub-db/main/index.json
EOF
  check "hub apply --synchronize adds the hub and fetches its templates" prints "synchronized" admin hub apply -f hub.yml --synchronize
  check "hub diff finds no drift" exits_with 0 admin hub diff -f hub.yml
  check "hub resync fetches it again" prints "synchronized" admin hub resync -i "$(id_of hub.yml)"
  check "hub delete deletes it" exits_with 0 admin hub delete -i "$(id_of hub.yml)" --yes

  # For the new team only, which runs nothing, so none is ever called.
  cat >webhook.yml <<EOF
# Written by e2e/platform.sh; for a team that runs nothing.
scope: TEAM
team: $SCOPE_TEAM
name: $MARK-$RUN-webhook
url: https://cli-e2e.invalid/webhook
events:
  - experiment.execution.created
targetAttributeIncludes:
  - '*'
headers: {}
EOF
  cat >preflight.yml <<EOF
# Written by e2e/platform.sh; for a team that runs nothing.
scope: TEAM
team: $SCOPE_TEAM
name: $MARK-$RUN-preflight
url: https://cli-e2e.invalid/preflight
events: []
targetAttributeIncludes:
  - '*'
headers: {}
EOF
  cat >preflight-action.yml <<EOF
# Written by e2e/platform.sh; for a team that runs nothing.
scope: TEAM
team: $SCOPE_TEAM
name: $MARK-$RUN-preflight-action
preflightActionId: com.steadybit.extension_loadtest.preflight.github-action
inflightInterval: 10s
inflightTimeout: 5s
EOF
  for kind in webhook preflight preflight-action; do
    check "integration $kind apply creates it" applied_with_id $kind.yml admin integration $kind apply
    check "integration $kind diff finds no drift" exits_with 0 admin integration $kind diff -f $kind.yml
    check "integration $kind delete deletes it" exits_with 0 admin integration $kind delete -i "$(id_of $kind.yml)"
  done

  check "team delete deletes the team" exits_with 0 admin team delete -k "$SCOPE_TEAM" --purge-experiments --yes
  check "environment delete deletes the environment" exits_with 0 admin environment delete -i "$ENV_ID" --yes

  check "audit-log reads today's entries" exits_with 0 admin audit-log --from "$(date -u +%F)" -t json
  check "killswitch status reads it; nothing here turns it on" prints "kill switch is" admin killswitch status
  check "reports read the platform" exits_with 0 admin report teams -t json
  check "license show reads the license" exits_with 0 admin license show
else
  echo "  skip  the checks that need STEADYBIT_E2E_ADMIN_TOKEN"
fi

echo
if [ "$failures" -eq 0 ]; then
  echo "all platform checks passed"
else
  echo "$failures platform check(s) failed"
fi
exit "$failures"

#!/usr/bin/env bash
# SPDX-License-Identifier: MIT
# SPDX-FileCopyrightText: 2026 Steadybit GmbH

# Uses the CLI against a real platform, the way a pipeline does: experiments applied from
# files, run in parallel with a report, checked for drift, canceled by the SIGTERM a
# canceled job sends, refused while another runs, and found again by a gate. Everything
# only waits, belongs to one test team, and is deleted afterwards, also when a check fails.
#
# Needs STEADYBIT_TOKEN (a team token) and STEADYBIT_URL, and `steadybit` on the PATH.
# STEADYBIT_E2E_TEAM and STEADYBIT_E2E_ENVIRONMENT name the team and its environment.

set -uo pipefail

: "${STEADYBIT_TOKEN:?a team access token is needed}"
TEAM=${STEADYBIT_E2E_TEAM:-CLI}
ENVIRONMENT=${STEADYBIT_E2E_ENVIRONMENT:-Global}
# Every experiment of this suite has a name starting with this, which is how a later run
# finds and deletes what a crashed one left behind.
MARK=cli-e2e-ci
RUN=${GITHUB_RUN_ID:-local-$$}

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
  rm -rf "$work"
}
trap cleanup EXIT

echo "steadybit $(steadybit --version) against ${STEADYBIT_URL:-https://platform.steadybit.com}, team $TEAM"

experiment experiments/a.yml a 5s
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
  steadybit experiment run -f experiments --yes --parallel 2 --report report.xml >run.log 2>&1 || { tail -n 20 run.log; exit 1; }
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
# --no-wait has to notice instead of passing.
check "a run started with --no-wait begins" exits_with 0 steadybit experiment run -k "$LONG" --yes --no-wait --allowParallel
until_run_is "$LONG" RUNNING
check "--no-wait fails on a run the platform refused" exits_with 1 steadybit experiment run -k "$A" --yes --no-wait
check "execution list --fail-on-match gates on that canceled run" exits_with 1 \
  steadybit execution list --key "$A" --state CANCELED --ended-from "$(date -u +%F)" --limit 1 --fail-on-match
check "execution list prints the platform's runs as JSON" sh -c "
  [ \"\$(steadybit execution list --team $TEAM --limit 2 --jq length 2>/dev/null)\" -ge 1 ]
"

echo
if [ "$failures" -eq 0 ]; then
  echo "all platform checks passed"
else
  echo "$failures platform check(s) failed"
fi
exit "$failures"

#!/bin/sh
# SPDX-License-Identifier: MIT
# SPDX-FileCopyrightText: 2026 Steadybit GmbH

# Smoke tests for the packaged CLI. Everything here needs a real process: an exit
# status, a terminal, or the spawn of a subcommand. Anything that can be asserted
# in-process belongs in the vitest suite instead, so these stay at the level of
# "did it exit correctly" rather than checking output in detail.

set -u
failures=0

# Driving the prompts needs a pty allocator, which the shipped image has no reason to
# carry. Installed here rather than baked into a second image: a derived image has to
# name its base by tag, and a tag can silently resolve to something from a registry
# instead of the build under test.
if ! command -v expect >/dev/null 2>&1; then
  apk add --no-cache expect >/dev/null 2>&1 || {
    echo "cannot install expect, which the interactive checks need"
    exit 1
  }
fi

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

exits_with() {
  expected=$1
  shift
  "$@" >/dev/null 2>&1
  actual=$?
  [ "$actual" -eq "$expected" ] || {
    echo "        expected exit $expected, got $actual"
    return 1
  }
}

echo "steadybit CLI container smoke tests"

check "--version succeeds" exits_with 0 steadybit --version
check "--help succeeds" exits_with 0 steadybit --help
check "a subcommand is spawned and runs" exits_with 0 steadybit experiment --help
check "an unknown command fails" exits_with 1 steadybit definitely-not-a-command

# The access token is resolved before anything else, so this is the guard on every
# platform-touching command.
check "a missing access token fails" exits_with 1 env STEADYBIT_TOKEN= steadybit experiment get -k ADM-1
check "an unreachable platform fails" exits_with 1 \
  env STEADYBIT_TOKEN=t STEADYBIT_URL=http://127.0.0.1:1 steadybit experiment get -k ADM-1

# Colour is gated on stdout being a terminal. Both halves are checked here rather than
# depending on how the container was started: the pipe below is genuinely not a
# terminal, and expect genuinely provides one.
check "output is clean when piped" sh -c \
  '! env STEADYBIT_TOKEN= steadybit experiment get -k ADM-1 2>&1 | grep -q "$(printf "\033")"'
check "output is coloured on a terminal" expect /e2e/colour-on-tty.exp

# The interactive flow, driven through a pty. Writes into the container's own home.
check "profile add stores what was typed" sh -c '
  expect /e2e/add-profile.exp >/dev/null 2>&1 || exit 1
  grep -q "\"name\": \"e2e\"" "$HOME/.steadybit/profiles.json" || exit 1
  steadybit config profile list | grep -q "e2e"
'

check "ctrl-c during a prompt exits 130" sh -c '
  expect /e2e/cancel-profile.exp >/dev/null 2>&1
  [ $? -eq 130 ]
'

check "ctrl-c leaves no stack trace and no profile" sh -c '
  rm -rf "$HOME/.steadybit"
  out=$(expect /e2e/cancel-profile.exp 2>&1)
  echo "$out" | grep -q "ExitPromptError" && exit 1
  [ ! -f "$HOME/.steadybit/profiles.json" ]
'

echo
if [ "$failures" -eq 0 ]; then
  echo "all container smoke tests passed"
else
  echo "$failures container smoke test(s) failed"
fi
exit "$failures"

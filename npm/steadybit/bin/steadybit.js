#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Runs the steadybit binary for this platform, which npm installed as one of the
// optional dependencies. The CLI itself is a Go binary; this only finds and starts it.
'use strict';

const { spawn } = require('node:child_process');
const path = require('node:path');

const platforms = {
  'darwin arm64': '@steadybit/cli-darwin-arm64',
  'darwin x64': '@steadybit/cli-darwin-x64',
  'linux arm64': '@steadybit/cli-linux-arm64',
  'linux x64': '@steadybit/cli-linux-x64',
  'win32 arm64': '@steadybit/cli-win32-arm64',
  'win32 x64': '@steadybit/cli-win32-x64',
};

function binaryPath() {
  const pkg = platforms[`${process.platform} ${process.arch}`];
  if (!pkg) {
    fail(`The Steadybit CLI is not available for ${process.platform} ${process.arch}.`);
  }
  const executable = process.platform === 'win32' ? 'steadybit.exe' : 'steadybit';
  try {
    return require.resolve(path.posix.join(pkg, 'bin', executable));
  } catch {
    fail(
      `The package ${pkg}, which holds the CLI for this platform, is not installed. ` +
        'It is an optional dependency; reinstall without --no-optional / --omit=optional:\n\n    npm install -g steadybit'
    );
  }
}

function fail(message) {
  console.error(message);
  process.exit(1);
}

const child = spawn(binaryPath(), process.argv.slice(2), { stdio: 'inherit' });

// Ctrl-C reaches the CLI directly, as they share the terminal; this process only waits
// for it, so that the CLI decides how to end, and its exit status (130 for Ctrl-C)
// becomes ours. A signal sent to this process alone, as a CI runner does, is passed on.
for (const signal of ['SIGINT', 'SIGTERM', 'SIGHUP']) {
  process.on(signal, () => {
    if (signal !== 'SIGINT') {
      child.kill(signal);
    }
  });
}

child.on('error', error => fail(`Failed to start the Steadybit CLI: ${error.message}`));
child.on('exit', (code, signal) => {
  if (signal) {
    process.kill(process.pid, signal);
    return;
  }
  process.exit(code ?? 1);
});

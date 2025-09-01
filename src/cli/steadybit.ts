#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Command, Option } from 'commander';
import colors from 'colors/safe';
import { satisfies } from 'semver';
import { enableRequestLogging } from '../api/http';

// eslint-disable-next-line @typescript-eslint/no-require-imports
const requiredNodejsVersion = require('../../package.json').engines.node;
const actualNodejsVersion = process.version;

if (!satisfies(actualNodejsVersion, requiredNodejsVersion)) {
  const help = `
Node.js version ${actualNodejsVersion} is not supported. The Steadybit CLI
requires a Node.js version that satisfies the following version range:

                       ${colors.bold(requiredNodejsVersion)}

We recommend to install Node.js via a version manager. For example,
using the Node Version Manager (NVM):

               ${colors.bold('https://github.com/nvm-sh/nvm#readme')}
`;
  console.error(colors.red(help.trim()));
  process.exit(1);
}

new Command()
  // Prefer to load at runtime directly from the package.json to simplify
  // the TypeScript build. Without this, we would have to make the build
  // more complicated to adapt the root dir accordingly.
  // eslint-disable-next-line
  .version(require('../../package.json').version)
  .addOption(new Option('-v, --verbose', 'Enable verbose logging').default(false))
  .hook('preSubcommand', thisCommand => {
    if (thisCommand.opts().verbose) {
      enableRequestLogging();
    }
  })
  .command('advice', 'Show/verify advice status.')
  .command('config', 'Show/modify the CLI configuration and authentication profiles.')
  .command('experiment', 'Check and run experiments.')
  .parseAsync(process.argv);

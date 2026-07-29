#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Command, Option } from 'commander';
import colors from '../colors.ts';
import { satisfies } from 'semver';
import { enableRequestLogging } from '../api/http.ts';
import { packageJson } from '../packageJson.ts';

const requiredNodejsVersion = packageJson.engines.node;
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
  .version(packageJson.version)
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

#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Command } from 'commander';

import { setVersion } from '../defRepo/setVersion';
import { check } from '../defRepo/check';

const program = new Command();

program
  .command('set-version')
  .description('Set the versions of policies and tasks referenced within task and policy definition files in this repository.')
  .requiredOption('-v, --version <version>', 'The version to set.')
  .option('-d, --directory <directory>', 'The directory to search for task and policy files.', process.cwd())
  .action(setVersion);

program
  .command('check')
  .description('Checks that the tasks and policies checked into the GitHub repository are valid. Assumes that the tasks and policies are checked in.')
  .option('-v, --version <version>', 'The version to check. Must resolve to a git ref existing in the repository.')
  .option('-d, --directory <directory>', 'The directory to search for task and policy files.', process.cwd())
  .action(check);

program.parseAsync(process.argv);

#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Command } from 'commander';
import { execute } from '../experiment/exec';
import { requirePlatformAccess } from './requirements';

const program = new Command();

program
  .command('exec')
  .description('Executes an experiment run.')
  .requiredOption('-k, --key <key>', 'The experiment key.')
  .option('--no-wait', 'Do not wait for experiment execution to finish.')
  .option(
    '--yes',
    'Skip the prompt asking for experiment execution confirmation. Not necessary when no TTY is attached.',
    false
  )
  .action(requirePlatformAccess(execute));

program.parseAsync(process.argv);

#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Command } from 'commander';
import { execute } from '../experiment/exec';
import { requirePlatformAccess } from './requirements';

const program = new Command();

program
  .command('run')
  .alias('exec')
  .description('Executes an experiment run.')
  .requiredOption('-k, --key <key>', 'The experiment key.')
  .option('--no-wait', 'Do not wait for experiment run to finish.')
  .option(
    '--yes',
    'Skip the prompt asking for experiment run confirmation. Not necessary when no TTY is attached.',
    false
  )
  .action(requirePlatformAccess(execute));

program.parseAsync(process.argv);

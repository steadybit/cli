#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Command, Option } from 'commander';
import { executeExperiment } from '../experiment/exec';
import { getExperiment } from '../experiment/get';
import { applyExperiment } from '../experiment/apply';
import { deleteExperiment } from '../experiment/delete';
import { requirePlatformAccess } from './requirements';

const program = new Command();

program
  .command('run')
  .alias('exec')
  .description('Executes an experiment run. If a file is specified the experiment is saved before execution.')
  .addOption(new Option('-k, --key <key>', 'The experiment key.').conflicts('file'))
  .addOption(new Option('-f, --file <file>', 'The path to the experiment file.').conflicts('key'))
  .addOption(new Option('--no-wait', 'Do not wait for experiment run to finish.'))
  .addOption(new Option('--yes', 'Skip the prompt asking for experiment run confirmation. Not necessary when no TTY is attached.').default(false))
  .action(requirePlatformAccess(executeExperiment));

program
  .command('get')
  .description('Get an experiment from Steadybit. Output is written to file or stdout.')
  .addOption(new Option('-k, --key <key>', 'The experiment key.').makeOptionMandatory(true))
  .addOption(new Option('-f, --file <file>', 'The path to the experiment file.'))
  .action(requirePlatformAccess(getExperiment));

program
  .command('apply')
  .description('Upload an experiment to Steadybit. If a key is provided, an update is performed. Otherwise, the externalId from the yaml is used to create or update the experiment.')
  .addOption(new Option('-k, --key <key>', 'The experiment key.'))
  .addOption(new Option('-f, --file <file>', 'The path to the experiment file.').makeOptionMandatory(true))
  .action(requirePlatformAccess(applyExperiment));

program
  .command('delete')
  .description('Delete an experiment from Steadybit.')
  .addOption(new Option('-k, --key <key>', 'The experiment key.').makeOptionMandatory(true))
  .action(requirePlatformAccess(deleteExperiment));


program.parseAsync(process.argv);

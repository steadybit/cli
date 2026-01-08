#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Command, Option } from 'commander';
import { executeExperiments } from '../experiment/exec';
import { getExperiment } from '../experiment/get';
import { dump } from '../experiment/dump';
import { applyExperiments } from '../experiment/apply';
import { deleteExperiment } from '../experiment/delete';
import { requirePlatformAccess } from './requirements';

const program = new Command();

program
  .command('run')
  .alias('exec')
  .description('Executes an experiment run. If a file is specified the experiment is saved before execution.')
  .addOption(new Option('-k, --key <key>', 'The experiment key.').conflicts('file'))
  .addOption(
    new Option(
      '-f, --file <files...>',
      'The path to the experiment file or a directory containing multiple files.'
    ).conflicts('key')
  )
  .addOption(
    new Option(
      '-R, --recursive',
      'Process the directory used in -f, --file recursively. Useful when you want to manage related experiments organized within the same directory.'
    ).default(false)
  )
  .addOption(new Option('--no-wait', 'Do not wait for experiment run to finish.'))
  .addOption(
    new Option(
      '--yes',
      'Skip the prompt asking for experiment run confirmation. Not necessary when no TTY is attached.'
    ).default(false)
  )
  .addOption(
    new Option(
      '--allowParallel',
      'Skip the prompt warning about another experiment running and allow always parallel execution.'
    ).default(false)
  )
  .action(requirePlatformAccess(executeExperiments));

program
  .command('get')
  .description('Get an experiment from Steadybit. Output is written to file or stdout.')
  .addOption(new Option('-k, --key <key>', 'The experiment key.').makeOptionMandatory(true))
  .addOption(new Option('-f, --file <file>', 'The path to the experiment file.'))
  .addOption(
    new Option(
      '-t, --type <type>',
      'The output format of the experiment ("json" or "yaml"). (default: if a file with ".json"-suffix is given: "json", "yaml" otherwise.)' // intentionally documented here and not using .default(flags, description) as otherwise the file-extension-logic isn't working
    )
  )
  .action(requirePlatformAccess(getExperiment));

program
  .command('apply')
  .description(
    'Upload an experiment to Steadybit. If a key is provided, an update is performed. Otherwise, the externalId from the file is used to create or update the experiment.'
  )
  .addOption(new Option('-k, --key <key>', 'The experiment key.'))
  .addOption(
    new Option(
      '-f, --file <files...>',
      'The path to the experiment file or a directory containing multiple files'
    ).makeOptionMandatory(true)
  )
  .addOption(
    new Option(
      '-R, --recursive',
      'Process the directory used in -f, --file recursively. Useful when you want to manage related experiments organized within the same directory.'
    ).default(false)
  )
  .action(requirePlatformAccess(applyExperiments));

program
  .command('delete')
  .description('Delete an experiment from Steadybit.')
  .addOption(new Option('-k, --key <key>', 'The experiment key.').makeOptionMandatory(true))
  .action(requirePlatformAccess(deleteExperiment));

program
  .command('dump')
  .description('Dump all experiments and executions from all teams in Steadybit.')
  .addOption(new Option('-d, --directory <dir>', 'The path to dump all the experiments to').default('.'))
  .addOption(new Option('-t, --type <type>', 'The output format of the experiment ("json" or "yaml").').default('yaml'))
  .action(requirePlatformAccess(dump));

program.parseAsync(process.argv);

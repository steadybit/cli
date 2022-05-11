#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Command } from 'commander';

import { deleteServiceDefinition } from '../serviceDefinition/delete';
import { validateSemverRangeCommanderArgument } from '../semver';
import { verify } from '../serviceDefinition/verify';
import { apply } from '../serviceDefinition/apply';
import { open } from '../serviceDefinition/open';
import { init } from '../serviceDefinition/init';
import { exec } from '../serviceDefinition/exec';
import { show } from '../serviceDefinition/show';
import { requirePlatformAccess } from './requirements';

const program = new Command();

program
  .command('apply')
  .option('-f, --file <file>', 'Path to the service definition file.', '.steadybit.yml')
  .description('Upload a service definition.')
  .action(requirePlatformAccess(apply));
program
  .command('delete')
  .option('-i, --id <id>', 'Id to be deleted.')
  .option('-f, --file <file>', 'Path to the service definition file.', '.steadybit.yml')
  .description('Deletes a service definition.')
  .action(requirePlatformAccess(deleteServiceDefinition));
program
  .command('init')
  .description('Initialize a service definition file.')
  .action(requirePlatformAccess(init));
program
  .command('open')
  .option('-f, --file <file>', 'Path to the service definition file.', '.steadybit.yml')
  .description('Open the service in the Steadybit UI.')
  .action(requirePlatformAccess(open));
program
  .command('verify')
  .option('-f, --file <file>', 'Path to the service definition file.', '.steadybit.yml')
  .option('-pp, --print-parameters', 'Print task parameters when listing tasks.')
  .option('-pm, --print-matrix-context', 'Print the matrix execution context information when listing tasks.')
  .description('Read the current service definition and state.')
  .action(requirePlatformAccess(verify));
program
  .command('exec')
  .option('--no-wait', 'Do not wait for task executions to finish.')
  .option('--yes', 'Skip the prompt asking for experiment execution confirmation. Not necessary when no TTY is attached.', false)
  .option('--no-error-on-empty-task-set', 'Whether to end the process with an error when an execution is triggered for an empty set of tasks.')
  .option('--no-error-on-task-failure', 'Whether to end the process with an error when an executed task fails.')
  .option('-f, --file <file>', 'Path to the service definition file.', '.steadybit.yml')
  .option('-t, --task <tasks...>', 'Optional filter to limit execution to those tasks matching the given task name. Can appear multiple times')
  .description('Execute tasks defined for a service.')
  .action(requirePlatformAccess(exec));
program
  .command('show')
  .option('-f, --file <file>', 'Path to the service definition file.', '.steadybit.yml')
  .option('-n, --name <name>', 'Optional task name to filter the result list.')
  .option(
    '-v, --version <version>',
    'Optional task version to filter the result list.',
    validateSemverRangeCommanderArgument
  )
  .description('Show a list of tasks and policies referenced by this service definition.')
  .action(requirePlatformAccess(show));

program.parseAsync(process.argv);

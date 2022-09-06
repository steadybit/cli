#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Command } from 'commander';
import { apply } from '../policyBinding/apply';
import { deletePolicyBinding } from '../policyBinding/delete';
import { exec } from '../policyBinding/exec';
import { init } from '../policyBinding/init';
import { open } from '../policyBinding/open';
import { requirePlatformAccess } from './requirements';
import { show } from '../policyBinding/show';
import { validateSemverRangeCommanderArgument } from '../semver';
import { verify } from '../policyBinding/verify';

const program = new Command();

program
  .command('apply')
  .option('-f, --file <file>', 'Path to the policy binding file.', '.steadybit.yml')
  .description('Upload a policy binding.')
  .action(requirePlatformAccess(apply));
program
  .command('delete')
  .option('-i, --id <id>', 'Id to be deleted.')
  .option('-f, --file <file>', 'Path to the policy binding file.', '.steadybit.yml')
  .description('Deletes a policy binding.')
  .action(requirePlatformAccess(deletePolicyBinding));
program.command('init').description('Initialize a policy binding file.').action(requirePlatformAccess(init));
program
  .command('open')
  .option('-f, --file <file>', 'Path to the policy binding file.', '.steadybit.yml')
  .description('Open the policy binding in the Steadybit UI.')
  .action(requirePlatformAccess(open));
program
  .command('verify')
  .option('-f, --file <file>', 'Path to the policy binding file.', '.steadybit.yml')
  .option('-pp, --print-parameters', 'Print task parameters when listing tasks.')
  .option('-pm, --print-matrix-context', 'Print the matrix execution context information when listing tasks.')
  .description('Read the current policy binding and state.')
  .action(requirePlatformAccess(verify));
program
  .command('run')
  .alias('exec')
  .option('--no-wait', 'Do not wait for task runs to finish.')
  .option(
    '--yes',
    'Skip the prompt asking for experiment run confirmation. Not necessary when no TTY is attached.',
    false
  )
  .option(
    '--no-error-on-empty-task-set',
    'Whether to end the process with an error when a run is triggered for an empty set of tasks.'
  )
  .option('--no-error-on-task-failure', 'Whether to end the process with an error when a task fails.')
  .option('-f, --file <file>', 'Path to the policy binding file.', '.steadybit.yml')
  .option(
    '-t, --task <tasks...>',
    'Optional filter to limit runs to those tasks matching the given task name. Can appear multiple times.'
  )
  .option('-pp, --print-parameters', 'Print task parameters when listing tasks.')
  .option('-pm, --print-matrix-context', 'Print the matrix execution context information when listing tasks.')
  .description('Run tasks defined through a policy binding.')
  .action(requirePlatformAccess(exec));
program
  .command('show')
  .option('-f, --file <file>', 'Path to the policy binding file.', '.steadybit.yml')
  .option('-n, --name <name>', 'Optional task name to filter the result list.')
  .option(
    '-v, --version <version>',
    'Optional task version to filter the result list.',
    validateSemverRangeCommanderArgument
  )
  .description('Show a list of tasks and policies referenced by this policy binding.')
  .action(requirePlatformAccess(show));

program.parseAsync(process.argv);

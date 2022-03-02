#!/usr/bin/env node
/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import { Command } from 'commander';

import { deleteServiceDefinition } from '../serviceDefinition/delete';
import { validateSemverRangeCommanderArgument } from '../semver';
import { verify } from '../serviceDefinition/verify';
import { apply } from '../serviceDefinition/apply';
import { open } from '../serviceDefinition/open';
import { init } from '../serviceDefinition/init';
import { show } from '../serviceDefinition/show';

const program = new Command();

program
  .command('apply')
  .option('-f, --file <file>', 'Path to the service definition file.', '.steadybit.yml')
  .description('Upload a service definition.')
  .action(apply);
program
  .command('delete')
  .option('-i, --id <id>', 'Id to be deleted.')
  .option('-f, --file <file>', 'Path to the service definition file.', '.steadybit.yml')
  .description('Deletes a service definition.')
  .action(deleteServiceDefinition);
program
  .command('init')
  .description('Initialize a service definition file.')
  .action(init);
program
  .command('open')
  .option('-f, --file <file>', 'Path to the service definition file.', '.steadybit.yml')
  .description('Open the service in the Steadybit UI.')
  .action(open);
program
  .command('verify')
  .option('-f, --file <file>', 'Path to the service definition file.', '.steadybit.yml')
  .option('-pp, --print-parameters', 'Print task parameters when listing tasks.')
  .option('-pm, --print-matrix-context', 'Print the matrix execution context information when listing tasks.')
  .description('Read the current service definition and state.')
  .action(verify);
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
  .action(show);

program.parseAsync(process.argv);

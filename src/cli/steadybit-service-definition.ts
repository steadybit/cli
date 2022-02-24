#!/usr/bin/env node
/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import { Command } from 'commander';

import { deleteServiceDefinition } from '../serviceDefinition/delete';
import { verify } from '../serviceDefinition/verify';
import { apply } from '../serviceDefinition/apply';
import { open } from '../serviceDefinition/open';
import { init } from '../serviceDefinition/init';

const program = new Command();

program
  .command('apply <path-to-service-definition-file>')
  .description('Upload a service definition.')
  .action(apply);
program
  .command('delete <id-or-path-to-service-definition-file>')
  .description('Deletes a service definition.')
  .action(deleteServiceDefinition);
program
  .command('init')
  .description('Initialize a service definition file.')
  .action(init);
program
  .command('open <path-to-service-definition-file>')
  .description('Open the service in the Steadybit UI.')
  .action(open);
program
  .command('verify <path-to-service-definition-file>')
  .option('-pp, --print-parameters', 'Print task parameters when listing tasks.')
  .option('-pm, --print-matrix-context', 'Print the matrix execution context information when listing tasks.')
  .description('Read the current service definition and state.')
  .action(verify);

program.parseAsync(process.argv);

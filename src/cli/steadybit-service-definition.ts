#!/usr/bin/env node

import { Command } from 'commander';

import { deleteServiceDefinition } from '../serviceDefinition/delete';
import { apply } from '../serviceDefinition/apply';
import { init } from '../serviceDefinition/init';

const program = new Command();

program
  .command('init')
  .description('Initialize a service definition file.')
  .action(init);
program
  .command('apply <path-to-service-definition-file>')
  .description('Upload a service definition.')
  .action(apply);
program
  .command('delete <id-or-path-to-service-definition-file>')
  .description('Deletes a service definition.')
  .action(deleteServiceDefinition);

program.parseAsync(process.argv);

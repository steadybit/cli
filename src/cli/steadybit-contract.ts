#!/usr/bin/env node

import { Command } from 'commander';

import { establish } from '../contract/establish';
import { init } from '../contract/init';

const program = new Command();

program
  .command('establish <contract-path>')
  .description('Upload a contract definition to the Steadybit server.')
  .action(establish);
program
  .command('init')
  .description('Initialize a contract definition file.')
  .action(init);

program.parseAsync(process.argv);

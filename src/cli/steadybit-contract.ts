#!/usr/bin/env node

import { Command } from 'commander';

import { establish } from '../contract/establish';

const program = new Command();

program.command('establish <contract-path>').description('yo!').action(establish);

program.parseAsync(process.argv);

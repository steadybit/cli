#!/usr/bin/env node

import { Command } from 'commander';

import { establish } from '../policy/establish';

const program = new Command();

program.command('establish <policy-path>').description('yo!').action(establish);

program.parseAsync(process.argv);

#!/usr/bin/env node
/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import { Command } from 'commander';

import { remove } from '../config/profile/remove';
import { select } from '../config/profile/select';
import { list } from '../config/profile/list';
import { add } from '../config/profile/add';

const program = new Command();

program
  .command('add')
  .description('Interactively configure a new profile.')
  .action(add);
program
  .command('list')
  .description('List all configured profiles.')
  .action(list);
program
  .command('ls')
  .description('Alias for list.')
  .action(list);
  program
  .command('remove')
  .description('Interactively remove an existing profile.')
  .action(remove);
program
  .command('select')
  .description('Interactively change the currently active profile.')
  .action(select);


program.parseAsync(process.argv);

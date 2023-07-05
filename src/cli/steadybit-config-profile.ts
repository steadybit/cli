#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH
import { Command, Option } from 'commander';
import { remove } from '../config/profile/remove';
import { select } from '../config/profile/select';
import { list } from '../config/profile/list';
import { add } from '../config/profile/add';
import { defaultBaseUrl } from '../config';

const program = new Command();

program
  .command('add')
  .description('Configure a new profile (interactively or via options).')
  .addOption(new Option('-n, --name <name>', 'Name of the profile'))
  .addOption(new Option('-b, --baseUrl <url>', 'Base URL to be used').default(defaultBaseUrl))
  .addOption(new Option('-t, --token <token>', 'Team API token'))
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

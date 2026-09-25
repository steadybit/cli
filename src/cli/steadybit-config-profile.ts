#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH
import { Command, Option } from 'commander';
import { remove } from '../config/profile/remove.ts';
import { select } from '../config/profile/select.ts';
import { list } from '../config/profile/list.ts';
import { add } from '../config/profile/add.ts';
import { defaultBaseUrl } from '../config/index.ts';
import { withExamples } from './help.ts';

const program = new Command();

withExamples(
  program
    .command('add')
    .description('Configure a new profile (interactively or via options).')
    .addOption(new Option('-n, --name <name>', 'Name of the profile'))
    .addOption(new Option('-b, --baseUrl <url>', 'Base URL to be used').default(defaultBaseUrl))
    .addOption(new Option('-t, --token <token>', 'Team API token'))
    .action(add),
  ['steadybit config profile add', 'steadybit config profile add -n prod -t "$STEADYBIT_TOKEN"']
);
withExamples(program.command('list').description('List all configured profiles.').action(list), [
  'steadybit config profile list',
]);
withExamples(program.command('ls').description('Alias for list.').action(list), ['steadybit config profile ls']);
withExamples(program.command('remove').description('Interactively remove an existing profile.').action(remove), [
  'steadybit config profile remove',
]);
withExamples(
  program.command('select').description('Interactively change the currently active profile.').action(select),
  ['steadybit config profile select']
);

program.parseAsync(process.argv);

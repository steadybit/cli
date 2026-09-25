#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { Command, Option } from 'commander';
import { requirePlatformAccess } from './requirements.ts';
import { withExamples } from './help.ts';
import { getTemplate, listTemplates } from '../template/commands.ts';

const program = new Command();

withExamples(
  program
    .command('list')
    .description('List experiment templates. Filters of the same kind match any of the given values.')
    .addOption(new Option('--tag <tags...>', 'Only list templates with one of these tags.'))
    .addOption(new Option('--target-type <types...>', 'Only list templates targeting one of these target types.'))
    .addOption(new Option('--action <actions...>', 'Only list templates using one of these actions.'))
    .addOption(new Option('--search <phrases...>', 'Only list templates whose title or description match.'))
    .action(requirePlatformAccess(listTemplates)),
  [
    'steadybit template list',
    'steadybit template list --search kubernetes --action com.steadybit.extension_host.stress-cpu',
  ]
);

withExamples(
  program
    .command('get')
    .description('Get an experiment template. Output is written to file or stdout.')
    .addOption(new Option('-i, --id <id>', 'The experiment template id.').makeOptionMandatory(true))
    .addOption(new Option('-f, --file <file>', 'The path to write the template to.'))
    .addOption(
      new Option(
        '-t, --type <type>',
        'The output format ("json" or "yaml"). (default: "json" if the file ends in ".json", "yaml" otherwise.)'
      )
    )
    .addOption(
      new Option(
        '--placeholders',
        'Only output the template placeholders, as a file to fill in and pass to --placeholders.'
      )
    )
    .action(requirePlatformAccess(getTemplate)),
  [
    'steadybit template get -i d7e65100-1d20-4980-be87-c351704910b8',
    'steadybit template get -i d7e65100-1d20-4980-be87-c351704910b8 --placeholders -f values.yml',
  ]
);

program.parseAsync(process.argv);

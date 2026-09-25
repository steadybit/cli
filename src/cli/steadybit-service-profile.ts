#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { Command, Option } from 'commander';
import { requirePlatformAccess } from './requirements.ts';
import { withExamples } from './help.ts';
import {
  applyServiceProfiles,
  deleteServiceProfile,
  getServiceProfile,
  listServiceProfiles,
} from '../serviceProfile/commands.ts';

const program = new Command();

const PROFILE_ID = '019eacd7-fb2c-733a-bed5-99a935323db5';

function idOption() {
  return new Option('-i, --id <id>', 'The service profile id.').makeOptionMandatory(true);
}

withExamples(
  program
    .command('list')
    .description('List service profiles.')
    .addOption(new Option('--name <name>', 'Only list profiles whose name contains this.'))
    .addOption(new Option('--origin <origins...>', 'Only list "provided" or "custom" profiles.'))
    .addOption(new Option('--default', 'Only list the default profile.'))
    .action(requirePlatformAccess(listServiceProfiles)),
  ['steadybit service-profile list', 'steadybit service-profile list --origin custom']
);

withExamples(
  program
    .command('get')
    .description('Get a service profile. Output is written to file or stdout.')
    .addOption(idOption())
    .addOption(new Option('-f, --file <file>', 'The path to write the service profile to.'))
    .addOption(
      new Option(
        '-t, --type <type>',
        'The output format ("json" or "yaml"). (default: "json" if the file ends in ".json", "yaml" otherwise.)'
      )
    )
    .action(requirePlatformAccess(getServiceProfile)),
  [`steadybit service-profile get -i ${PROFILE_ID} -f profile.yml`]
);

withExamples(
  program
    .command('apply')
    .description(
      'Create or update service profiles from files. A file without an id creates a profile, and the new id is written back to it.'
    )
    .addOption(
      new Option(
        '-f, --file <files...>',
        'The path to the service profile file or a directory containing multiple files.'
      ).makeOptionMandatory(true)
    )
    .addOption(new Option('-R, --recursive', 'Process the directory used in -f, --file recursively.').default(false))
    .addOption(
      new Option(
        '--delete-experiments',
        'Delete the provided experiments of services that use templates removed from the profile.'
      ).default(false)
    )
    .action(requirePlatformAccess(applyServiceProfiles)),
  ['steadybit service-profile apply -f profile.yml']
);

withExamples(
  program
    .command('delete')
    .description('Delete a custom service profile.')
    .addOption(idOption())
    .action(requirePlatformAccess(deleteServiceProfile)),
  [`steadybit service-profile delete -i ${PROFILE_ID}`]
);

program.parseAsync(process.argv);

#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { Command, Option } from 'commander';
import { requirePlatformAccess } from './requirements.ts';
import { withExamples } from './help.ts';
import { collectKeyValue, parseDecimal } from './options.ts';
import {
  applyServices,
  deleteService,
  getService,
  getServiceVariables,
  linkExperiment,
  listServiceExperiments,
  listServices,
  provideExperiment,
  setServiceVariables,
  showServiceRisk,
  unlinkExperiment,
} from '../service/commands.ts';

const program = new Command();

const SERVICE_ID = '019cd80d-a4c9-775b-bdf8-2672a280ce7c';

function idOption() {
  return new Option('-i, --id <id>', 'The service id.').makeOptionMandatory(true);
}

function typeOption() {
  return new Option(
    '-t, --type <type>',
    'The output format ("json" or "yaml"). (default: "json" if the file ends in ".json", "yaml" otherwise.)'
  );
}

withExamples(
  program
    .command('list')
    .description('List services. Filters of the same kind match any of the given values.')
    .addOption(new Option('--team <keys...>', 'Only list services of these teams, by team key.'))
    .addOption(new Option('--environment <names...>', 'Only list services in these environments.'))
    .addOption(new Option('--experiment <keys...>', 'Only list services these experiments are linked to.'))
    .action(requirePlatformAccess(listServices)),
  ['steadybit service list', 'steadybit service list --team ADM --environment Global']
);

withExamples(
  program
    .command('get')
    .description('Get a service. Output is written to file or stdout.')
    .addOption(idOption())
    .addOption(new Option('-f, --file <file>', 'The path to write the service to.'))
    .addOption(typeOption())
    .action(requirePlatformAccess(getService)),
  [`steadybit service get -i ${SERVICE_ID} -f service.yml`]
);

withExamples(
  program
    .command('apply')
    .description(
      'Create or update services from files. A file without an id creates a service, and the new id is written back to it.'
    )
    .addOption(
      new Option(
        '-f, --file <files...>',
        'The path to the service file or a directory containing multiple files.'
      ).makeOptionMandatory(true)
    )
    .addOption(new Option('-R, --recursive', 'Process the directory used in -f, --file recursively.').default(false))
    .addOption(
      new Option(
        '--delete-experiments',
        'When the service profile changes, delete provided experiments whose templates the new profile does not contain. Without it, such a change is refused.'
      ).default(false)
    )
    .action(requirePlatformAccess(applyServices)),
  ['steadybit service apply -f service.yml', 'steadybit service apply -f ./services -R']
);

withExamples(
  program
    .command('delete')
    .description('Delete a service.')
    .addOption(idOption())
    .action(requirePlatformAccess(deleteService)),
  [`steadybit service delete -i ${SERVICE_ID}`]
);

withExamples(
  program
    .command('risk')
    .description('Show the risk score of a service, overall, per category and per experiment.')
    .addOption(idOption())
    .addOption(new Option('-t, --type <type>', 'Print the raw risk as "json" or "yaml" instead of tables.'))
    .addOption(
      new Option(
        '--fail-above <score>',
        'Exit with a non-zero status when the overall risk is above this score.'
      ).argParser(parseDecimal)
    )
    .action(requirePlatformAccess(showServiceRisk)),
  [`steadybit service risk -i ${SERVICE_ID}`, `steadybit service risk -i ${SERVICE_ID} --fail-above 50`]
);

const experiment = program.command('experiment').description('Manage the experiments of a service.');

withExamples(
  experiment
    .command('list')
    .description(
      'List the experiments of a service: those provided by its service profile, created or not, and custom ones linked to it.'
    )
    .addOption(idOption())
    .addOption(new Option('--category <categories...>', 'Only list experiments in these categories.'))
    .addOption(new Option('--type <types...>', 'Only list "provided" or "custom" experiments.'))
    .action(requirePlatformAccess(listServiceExperiments)),
  [
    `steadybit service experiment list -i ${SERVICE_ID}`,
    `steadybit service experiment list -i ${SERVICE_ID} --type custom`,
  ]
);

withExamples(
  experiment
    .command('provide')
    .description("Create or update a provided experiment of a service from one of its service profile's templates.")
    .addOption(idOption())
    .addOption(
      new Option('--template <id>', 'The template, which must be part of the service profile.').makeOptionMandatory(
        true
      )
    )
    .addOption(
      new Option('-k, --experiment <key>', 'Update this existing provided experiment instead of creating one.')
    )
    .addOption(
      new Option('-p, --placeholder <KEY=VALUE>', 'A placeholder value. Repeat for more.').argParser(collectKeyValue)
    )
    .addOption(
      new Option('--placeholders <file>', 'A YAML/JSON file mapping placeholder keys to values. -p overrides entries.')
    )
    .addOption(
      new Option(
        '--no-reset-properties',
        'Keep the properties of an existing experiment instead of resetting them to the template.'
      )
    )
    .action(requirePlatformAccess(provideExperiment)),
  [
    `steadybit service experiment provide -i ${SERVICE_ID} --template d7e65100-1d20-4980-be87-c351704910b8 -p REPLICAS=3`,
  ]
);

withExamples(
  experiment
    .command('link')
    .description('Link an existing experiment to a service as a custom experiment.')
    .addOption(idOption())
    .addOption(new Option('-k, --experiment <key>', 'The experiment to link.').makeOptionMandatory(true))
    .addOption(new Option('--category <category>', 'The category to link it in.').makeOptionMandatory(true))
    .action(requirePlatformAccess(linkExperiment)),
  [`steadybit service experiment link -i ${SERVICE_ID} -k ADM-1 --category Redundancy`]
);

withExamples(
  experiment
    .command('unlink')
    .description('Remove a custom experiment from a service. The experiment itself is kept.')
    .addOption(idOption())
    .addOption(new Option('-k, --experiment <key>', 'The experiment to unlink.').makeOptionMandatory(true))
    .action(requirePlatformAccess(unlinkExperiment)),
  [`steadybit service experiment unlink -i ${SERVICE_ID} -k ADM-1`]
);

const variable = program.command('variable').description('Manage the variables of a service.');

withExamples(
  variable
    .command('get')
    .description('Print the variables of a service.')
    .addOption(idOption())
    .addOption(new Option('-t, --type <type>', 'The output format ("json" or "yaml").').default('yaml'))
    .action(requirePlatformAccess(getServiceVariables)),
  [`steadybit service variable get -i ${SERVICE_ID}`]
);

withExamples(
  variable
    .command('set')
    .description(
      'Set variables of a service, keeping the others. With --replace, the given variables become the only ones.'
    )
    .argument('[KEY=VALUE...]', 'Variables to set.')
    .addOption(idOption())
    .addOption(
      new Option(
        '-f, --file <file>',
        'A YAML/JSON file mapping variable names to values, which may be lists or select expressions.'
      )
    )
    .addOption(new Option('--replace', 'Remove every variable not given.').default(false))
    .action(requirePlatformAccess(setServiceVariables)),
  [
    `steadybit service variable set -i ${SERVICE_ID} endpoint=http://shop.internal region=eu`,
    `steadybit service variable set -i ${SERVICE_ID} -f variables.yml --replace`,
  ]
);

program.parseAsync(process.argv);

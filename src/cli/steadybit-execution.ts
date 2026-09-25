#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { Command, Option } from 'commander';
import { requirePlatformAccess } from './requirements.ts';
import { withExamples } from './help.ts';
import { parseDecimal } from './options.ts';
import { getExecution } from '../execution/get.ts';
import { cancel } from '../execution/cancel.ts';
import { addProperty, setProperty } from '../execution/property.ts';
import { downloadArtifacts, listArtifacts } from '../execution/artifacts.ts';

const program = new Command();

function idOption() {
  return new Option('-i, --id <id>', 'The experiment run id.').makeOptionMandatory(true).argParser(parseDecimal);
}

withExamples(
  program
    .command('get')
    .description(
      'Get an experiment run, including its steps and target executions. Output is written to file or stdout.'
    )
    .addOption(idOption())
    .addOption(new Option('-f, --file <file>', 'The path to write the experiment run to.'))
    .addOption(
      new Option(
        '-t, --type <type>',
        'The output format ("json" or "yaml"). (default: "json" if the file ends in ".json", "yaml" otherwise.)'
      )
    )
    .action(requirePlatformAccess(getExecution)),
  ['steadybit execution get -i 1234', 'steadybit execution get -i 1234 -t json | jq .state']
);

withExamples(
  program
    .command('cancel')
    .description('Cancel a running experiment run. The run stops as soon as its agents have been told.')
    .addOption(idOption())
    .action(requirePlatformAccess(cancel)),
  ['steadybit execution cancel -i 1234']
);

const property = program.command('property').description('Change the properties of an experiment run.');

withExamples(
  property
    .command('set')
    .description(
      'Set the value of a property of an experiment run. Only properties editable in a run can be changed. Several --value set a list property.'
    )
    .addOption(idOption())
    .addOption(new Option('-k, --key <key>', 'The property key.').makeOptionMandatory(true))
    .addOption(new Option('--value <values...>', 'The value to set.').makeOptionMandatory(true))
    .addOption(new Option('--json', 'Parse each value as JSON, to send a number or an object.').default(false))
    .action(requirePlatformAccess(setProperty)),
  [
    'steadybit execution property set -i 1234 -k approvedBy --value "Jane Doe"',
    'steadybit execution property set -i 1234 -k tickets --value SHOP-1 SHOP-2',
    'steadybit execution property set -i 1234 -k score --value 7 --json',
  ]
);

withExamples(
  property
    .command('add')
    .description('Add a value to a list property of an experiment run.')
    .addOption(idOption())
    .addOption(new Option('-k, --key <key>', 'The property key.').makeOptionMandatory(true))
    .addOption(new Option('--value <value>', 'The value to add.').makeOptionMandatory(true).argParser(v => [v]))
    .addOption(new Option('--json', 'Parse the value as JSON, to send a number or an object.').default(false))
    .action(requirePlatformAccess(addProperty)),
  ['steadybit execution property add -i 1234 -k tickets --value SHOP-3']
);

const artifact = program.command('artifact').description('List and download the artifacts of an experiment run.');

withExamples(
  artifact
    .command('list')
    .description('List the artifacts that the actions of an experiment run attached.')
    .addOption(idOption())
    .action(requirePlatformAccess(listArtifacts)),
  ['steadybit execution artifact list -i 1234']
);

withExamples(
  artifact
    .command('download')
    .description(
      'Download the artifacts of an experiment run into <directory>/<target execution>/<artifact>. Without filters, all of them are downloaded.'
    )
    .addOption(idOption())
    .addOption(new Option('-a, --artifact <artifact>', 'Only download artifacts with this id, usually the file name.'))
    .addOption(new Option('--target-execution <id>', 'Only download artifacts of this target execution.'))
    .addOption(new Option('-d, --directory <dir>', 'The directory to download into.').default('.'))
    .addOption(
      new Option(
        '-o, --output <file>',
        'Write the artifact to this file instead. Requires exactly one match.'
      ).conflicts('directory')
    )
    .action(requirePlatformAccess(downloadArtifacts)),
  [
    'steadybit execution artifact download -i 1234 -d ./artifacts',
    'steadybit execution artifact download -i 1234 -a jmeter-report.zip -o report.zip',
  ]
);

program.parseAsync(process.argv);

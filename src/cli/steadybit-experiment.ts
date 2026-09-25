#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Command, Option } from 'commander';
import { executeExperiments } from '../experiment/exec.ts';
import { getExperiment } from '../experiment/get.ts';
import { dump } from '../experiment/dump.ts';
import { applyExperiments } from '../experiment/apply.ts';
import { deleteExperiment } from '../experiment/delete.ts';
import { requirePlatformAccess } from './requirements.ts';
import { withExamples } from './help.ts';
import { collectKeyValue, parseDecimal } from './options.ts';

const program = new Command();

// Shared by `run` and `apply`, which both render an experiment from a template when
// --template is given instead of reading it from a file.
function addTemplateOptions(command: Command): Command {
  return command
    .addOption(
      new Option('--template <id>', 'Create the experiment from the experiment template with this id.').conflicts(
        'file'
      )
    )
    .addOption(new Option('--team <key>', 'With --template: the key of the team owning the experiment.'))
    .addOption(new Option('--environment <name>', 'With --template: the environment the experiment runs in.'))
    .addOption(
      new Option(
        '--external-id <id>',
        'With --template: an identifier of your own. Using the same one again updates the experiment it created before.'
      )
    )
    .addOption(
      new Option('-p, --placeholder <KEY=VALUE>', 'With --template: a placeholder value. Repeat for more.').argParser(
        collectKeyValue
      )
    )
    .addOption(
      new Option(
        '--placeholders <file>',
        'With --template: a YAML/JSON file mapping placeholder keys to values. -p overrides entries.'
      )
    )
    .addOption(
      new Option(
        '--variable <KEY=VALUE>',
        'With --template: an experiment variable to add to the experiment. Repeat for more.'
      ).argParser(collectKeyValue)
    )
    .addOption(
      new Option(
        '--no-reset-properties',
        'With --template: keep the properties of an existing experiment instead of resetting them to the template.'
      )
    );
}

const run = program
  .command('run')
  .alias('exec')
  .description('Executes an experiment run. If a file is specified the experiment is saved before execution.')
  .addOption(new Option('-k, --key <key>', 'The experiment key.').conflicts('file'))
  .addOption(
    new Option(
      '-f, --file <files...>',
      'The path to the experiment file or a directory containing multiple files.'
    ).conflicts('key')
  )
  .addOption(
    new Option(
      '-R, --recursive',
      'Process the directory used in -f, --file recursively. Useful when you want to manage related experiments organized within the same directory.'
    ).default(false)
  )
  .addOption(new Option('--no-wait', 'Do not wait for experiment run to finish.'))
  .addOption(
    new Option(
      '--yes',
      'Skip the prompt asking for experiment run confirmation. Not necessary when no TTY is attached.'
    ).default(false)
  )
  .addOption(
    new Option(
      '--allowParallel',
      'Skip the prompt warning about another experiment running and allow always parallel execution.'
    ).default(false)
  )
  .addOption(
    new Option(
      '--retries <number>',
      'Number of retries when the experiment fails validation (e.g., missing targets). 0 means no retry.'
    )
      .default(0)
      .argParser(parseDecimal)
  )
  .addOption(
    new Option('--retryInterval <seconds>', 'Interval in seconds between retries.').default(10).argParser(parseDecimal)
  )
  .addOption(
    new Option(
      '--execution-variable <KEY=VALUE>',
      'With --template: a variable for this run only, overriding experiment and environment variables. Repeat for more.'
    ).argParser(collectKeyValue)
  );
addTemplateOptions(run).action(requirePlatformAccess(executeExperiments));
withExamples(run, [
  'steadybit experiment run -k ADM-1',
  'steadybit experiment run -f experiment.yml --no-wait',
  'steadybit experiment run -f ./experiments -R --yes',
  'steadybit experiment run --template d7e65100-1d20-4980-be87-c351704910b8 --team ADM --environment Global -p CLUSTER=prod',
]);

const get = program
  .command('get')
  .description('Get an experiment from Steadybit. Output is written to file or stdout.')
  .addOption(new Option('-k, --key <key>', 'The experiment key.').makeOptionMandatory(true))
  .addOption(new Option('-f, --file <file>', 'The path to the experiment file.'))
  .addOption(
    new Option(
      '-t, --type <type>',
      'The output format of the experiment ("json" or "yaml"). (default: if a file with ".json"-suffix is given: "json", "yaml" otherwise.)' // intentionally documented here and not using .default(flags, description) as otherwise the file-extension-logic isn't working
    )
  )
  .action(requirePlatformAccess(getExperiment));
withExamples(get, ['steadybit experiment get -k ADM-1', 'steadybit experiment get -k ADM-1 -f experiment.json']);

const apply = program
  .command('apply')
  .description(
    'Upload an experiment to Steadybit. If a key is provided, an update is performed. Otherwise, the externalId from the file is used to create or update the experiment. With --template, the experiment is created from an experiment template instead of a file.'
  )
  .addOption(new Option('-k, --key <key>', 'The experiment key.'))
  .addOption(
    new Option(
      '-f, --file <files...>',
      'The path to the experiment file or a directory containing multiple files'
    ).conflicts('template')
  )
  .addOption(
    new Option(
      '-R, --recursive',
      'Process the directory used in -f, --file recursively. Useful when you want to manage related experiments organized within the same directory.'
    ).default(false)
  );
addTemplateOptions(apply).action(requirePlatformAccess(applyExperiments));
withExamples(apply, [
  'steadybit experiment apply -f experiment.yml',
  'steadybit experiment apply -f ./experiments -R',
  'steadybit experiment apply --template d7e65100-1d20-4980-be87-c351704910b8 --team ADM --external-id shop-latency -p CLUSTER=prod',
  'steadybit experiment apply --template d7e65100-1d20-4980-be87-c351704910b8 -k ADM-12 --placeholders values.yml',
]);

const del = program
  .command('delete')
  .description('Delete an experiment from Steadybit.')
  .addOption(new Option('-k, --key <key>', 'The experiment key.').makeOptionMandatory(true))
  .action(requirePlatformAccess(deleteExperiment));
withExamples(del, ['steadybit experiment delete -k ADM-1']);

const dumpCommand = program
  .command('dump')
  .description('Dump all experiments and executions from all teams in Steadybit.')
  .addOption(new Option('-d, --directory <dir>', 'The path to dump all the experiments to').default('.'))
  .addOption(new Option('-t, --type <type>', 'The output format of the experiment ("json" or "yaml").').default('yaml'))
  .addOption(
    new Option('--team <keys...>', 'Only dump the given teams, by team key. Defaults to every accessible team.')
  )
  .action(requirePlatformAccess(dump));
withExamples(dumpCommand, [
  'steadybit experiment dump -d ./dump',
  'steadybit experiment dump -d ./dump -t json --team ADM WEBHOOK',
]);

program.parseAsync(process.argv);

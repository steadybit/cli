#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { Command, Option } from 'commander';
import { requirePlatformAccess } from './requirements.ts';
import { withExamples } from './help.ts';
import { collectKeyValue } from './options.ts';
import {
  applySchedules,
  createSchedule,
  deleteSchedule,
  disableSchedule,
  enableSchedule,
  getSchedule,
  listSchedules,
  updateSchedule,
} from '../schedule/commands.ts';

const program = new Command();

function idOption() {
  return new Option('-i, --id <id>', 'The experiment schedule id.').makeOptionMandatory(true);
}

// Shared by `create` and `update`, which take the same fields as flags.
function addScheduleFields(command: Command): Command {
  return command
    .addOption(
      new Option(
        '--cron <expression>',
        'Run repeatedly on this Quartz cron expression (seconds first), e.g. "0 0 9 ? * MON-FRI".'
      )
    )
    .addOption(new Option('--start-at <timestamp>', 'Run once at this ISO 8601 time, e.g. 2026-10-01T09:00:00Z.'))
    .addOption(new Option('--timezone <zone>', 'The timezone of the cron expression, e.g. Europe/Berlin.'))
    .addOption(new Option('--allow-parallel', 'Run even when another experiment is running.'))
    .addOption(new Option('--no-allow-parallel', 'Skip the run when another experiment is running.'))
    .addOption(
      new Option(
        '--variable <KEY=VALUE>',
        'A variable for the scheduled runs, overriding experiment and environment variables. Repeat for more.'
      ).argParser(collectKeyValue)
    );
}

withExamples(
  program
    .command('list')
    .description('List experiment schedules.')
    .addOption(new Option('--team <keys...>', 'Only list schedules of these teams, by team key.'))
    .addOption(new Option('--experiment <keys...>', 'Only list schedules of these experiments, by experiment key.'))
    .action(requirePlatformAccess(listSchedules)),
  ['steadybit schedule list', 'steadybit schedule list --team ADM --experiment ADM-1 ADM-2']
);

withExamples(
  program
    .command('get')
    .description('Get an experiment schedule. Output is written to file or stdout.')
    .addOption(idOption())
    .addOption(new Option('-f, --file <file>', 'The path to write the schedule to.'))
    .addOption(
      new Option(
        '-t, --type <type>',
        'The output format ("json" or "yaml"). (default: "json" if the file ends in ".json", "yaml" otherwise.)'
      )
    )
    .action(requirePlatformAccess(getSchedule)),
  ['steadybit schedule get -i 01951394-727f-76a0-8675-c7519ebd0ff5 -f schedule.yml']
);

withExamples(
  program
    .command('apply')
    .description(
      'Create or update experiment schedules from files. A file without an id creates a schedule, and the new id is written back to it.'
    )
    .addOption(
      new Option(
        '-f, --file <files...>',
        'The path to the schedule file or a directory containing multiple files.'
      ).makeOptionMandatory(true)
    )
    .addOption(new Option('-R, --recursive', 'Process the directory used in -f, --file recursively.').default(false))
    .action(requirePlatformAccess(applySchedules)),
  ['steadybit schedule apply -f schedule.yml', 'steadybit schedule apply -f ./schedules -R']
);

withExamples(
  addScheduleFields(
    program
      .command('create')
      .description('Schedule an experiment, either repeatedly with --cron or once with --start-at.')
      .addOption(
        new Option('-k, --experiment <key>', 'The key of the experiment to schedule.').makeOptionMandatory(true)
      )
      .addOption(new Option('--disabled', 'Create the schedule disabled.'))
  ).action(requirePlatformAccess(createSchedule)),
  [
    'steadybit schedule create -k ADM-1 --cron "0 0 9 ? * MON-FRI" --timezone Europe/Berlin',
    'steadybit schedule create -k ADM-1 --start-at 2026-10-01T09:00:00Z --no-allow-parallel',
  ]
);

withExamples(
  addScheduleFields(
    program
      .command('update')
      .description('Change an experiment schedule. Only the given fields are changed.')
      .addOption(idOption())
  ).action(requirePlatformAccess(updateSchedule)),
  ['steadybit schedule update -i 01951394-727f-76a0-8675-c7519ebd0ff5 --cron "0 30 8 ? * *"']
);

withExamples(
  program
    .command('enable')
    .description('Enable an experiment schedule.')
    .addOption(idOption())
    .action(requirePlatformAccess(enableSchedule)),
  ['steadybit schedule enable -i 01951394-727f-76a0-8675-c7519ebd0ff5']
);

withExamples(
  program
    .command('disable')
    .description('Disable an experiment schedule without deleting it.')
    .addOption(idOption())
    .action(requirePlatformAccess(disableSchedule)),
  ['steadybit schedule disable -i 01951394-727f-76a0-8675-c7519ebd0ff5']
);

withExamples(
  program
    .command('delete')
    .description('Delete an experiment schedule.')
    .addOption(idOption())
    .action(requirePlatformAccess(deleteSchedule)),
  ['steadybit schedule delete -i 01951394-727f-76a0-8675-c7519ebd0ff5']
);

program.parseAsync(process.argv);

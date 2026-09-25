// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import fs from 'node:fs/promises';
import { Table } from 'console-table-printer';
import { abortExecution } from '../errors.ts';
import { format, output, readStructuredFile } from '../structuredFiles.ts';
import { resolveExperimentFiles } from '../experiment/files.ts';
import {
  fetchSchedule,
  fetchSchedules,
  patchSchedule,
  type PatchSchedule,
  removeSchedule,
  type Schedule,
  type UpsertSchedule,
  upsertSchedule,
} from './api.ts';

export interface ListOptions {
  team?: string[];
  experiment?: string[];
}

export async function listSchedules(options: ListOptions) {
  const schedules = await fetchSchedules(options);
  if (schedules.length === 0) {
    console.log('No experiment schedules found.');
    return;
  }
  const table = new Table({
    columns: [
      { name: 'id', title: 'Id', alignment: 'left' },
      { name: 'experiment', title: 'Experiment', alignment: 'left' },
      { name: 'when', title: 'When', alignment: 'left' },
      { name: 'enabled', title: 'Enabled', alignment: 'left' },
      { name: 'allowParallel', title: 'Parallel', alignment: 'left' },
    ],
  });
  table.addRows(
    schedules.map(s => ({
      id: s.id,
      experiment: s.experimentKey,
      when: s.cron ? `${s.cron}${s.timezone ? ` (${s.timezone})` : ''}` : (s.startAt ?? ''),
      enabled: String(s.enabled ?? true),
      allowParallel: String(s.allowParallel ?? true),
    }))
  );
  table.printTable();
}

export interface GetOptions {
  id: string;
  file?: string;
  type?: string;
}

export async function getSchedule(options: GetOptions) {
  await output(toFileContent(await fetchSchedule(options.id)), options);
  if (options.file) {
    console.log('Experiment schedule %s written to %s.', options.id, options.file);
  }
}

// What the platform reports about the last edit is not part of what can be sent, so it
// is left out of files, keeping `get` followed by `apply` a round trip.
function toFileContent(schedule: Schedule): UpsertSchedule {
  const content: Partial<Schedule> = { ...schedule };
  delete content.editedBy;
  delete content.lastUpdated;
  return content as UpsertSchedule;
}

export interface ApplyOptions {
  file: string[];
  recursive: boolean;
}

export async function applySchedules(options: ApplyOptions) {
  const files = await resolveExperimentFiles(options.file, options.recursive);
  for (const file of files) {
    const { content, datatype } = await readStructuredFile<Schedule>(file, 'schedule');
    if (!content?.experimentKey) {
      throw abortExecution("Schedule file '%s' does not name an experimentKey.", file);
    }
    const result = await upsertSchedule(toFileContent(content));
    if (!content.id) {
      // As `experiment apply` does with the key: the next apply of this file then
      // updates the schedule instead of creating a second one.
      // The id goes first in the file, where people look, even if the file had an empty one.
      const withId = Object.assign({ id: result.schedule.id }, content, { id: result.schedule.id });
      await fs.writeFile(file, format(withId, datatype), { encoding: 'utf8' });
    }
    console.log(
      'Experiment schedule %s for %s %s.',
      result.schedule.id,
      result.schedule.experimentKey,
      result.created ? 'created' : 'updated'
    );
  }
}

export interface ScheduleFields {
  cron?: string;
  startAt?: string;
  timezone?: string;
  allowParallel?: boolean;
  variable?: Record<string, string>;
}

export interface CreateOptions extends ScheduleFields {
  experiment: string;
  disabled?: boolean;
}

export async function createSchedule(options: CreateOptions) {
  rejectCronWithStartAt(options);
  requireOneOfCronOrStartAt(options);
  const { schedule } = await upsertSchedule({
    experimentKey: options.experiment,
    cron: options.cron,
    startAt: options.startAt,
    timezone: options.timezone,
    allowParallel: options.allowParallel,
    enabled: !options.disabled,
    variables: options.variable,
  });
  console.log('Experiment schedule %s for %s created.', schedule.id, schedule.experimentKey);
}

function requireOneOfCronOrStartAt(options: ScheduleFields) {
  if (!options.cron && !options.startAt) {
    throw abortExecution('Either --cron or --start-at must be specified.');
  }
}

// A schedule either repeats or runs once. The platform rejects both together on create,
// and on update would clear whichever it applied first.
function rejectCronWithStartAt(options: ScheduleFields) {
  if (options.cron && options.startAt) {
    throw abortExecution('--cron and --start-at cannot be combined.');
  }
}

export interface UpdateOptions extends ScheduleFields {
  id: string;
}

export async function updateSchedule(options: UpdateOptions) {
  rejectCronWithStartAt(options);
  const patch: PatchSchedule = {
    cron: options.cron,
    startAt: options.startAt,
    timezone: options.timezone,
    allowParallel: options.allowParallel,
    variables: options.variable,
  };
  if (Object.values(patch).every(value => value === undefined)) {
    throw abortExecution('Nothing to update. Pass at least one of the options, see --help.');
  }
  await patchAndReport(options.id, patch, 'updated');
}

export interface IdOptions {
  id: string;
}

export async function enableSchedule(options: IdOptions) {
  await patchAndReport(options.id, { enabled: true }, 'enabled');
}

export async function disableSchedule(options: IdOptions) {
  await patchAndReport(options.id, { enabled: false }, 'disabled');
}

async function patchAndReport(id: string, patch: PatchSchedule, outcome: string) {
  const schedule = await patchSchedule(id, patch);
  console.log('Experiment schedule %s for %s %s.', id, schedule.experimentKey, outcome);
}

export async function deleteSchedule(options: IdOptions) {
  await removeSchedule(options.id);
  console.log('Experiment schedule %s deleted.', options.id);
}

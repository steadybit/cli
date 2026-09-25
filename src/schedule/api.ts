// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import type { Schemas } from '../api/schemas.ts';
import { ApiError } from '../api/error.ts';
import { executeApiCall } from '../api/http.ts';
import { abortExecution, abortExecutionWithError } from '../errors.ts';

export type Schedule = Schemas['ExperimentScheduleAO'];
export type UpsertSchedule = Schemas['UpsertExperimentScheduleAO'];
export type PatchSchedule = Schemas['PatchExperimentScheduleAO'];

function notFoundOr(e: unknown, id: string, msg: string): Error {
  if (e instanceof ApiError && e.status === 404) {
    return abortExecution('Experiment schedule %s not found.', id);
  }
  return abortExecutionWithError(e, msg, id);
}

export async function fetchSchedules(filter: { team?: string[]; experiment?: string[] }): Promise<Schedule[]> {
  try {
    const response = await executeApiCall({
      method: 'GET',
      path: '/api/experiments/schedules/v2',
      queryParameters: { team: filter.team, experiment: filter.experiment },
    });
    return (await response.json()) as Schedule[];
  } catch (e) {
    throw abortExecutionWithError(e, 'Failed to get the experiment schedules');
  }
}

export async function fetchSchedule(id: string): Promise<Schedule> {
  try {
    const response = await executeApiCall({
      method: 'GET',
      path: `/api/experiments/schedules/${encodeURIComponent(id)}`,
    });
    return (await response.json()) as Schedule;
  } catch (e) {
    throw notFoundOr(e, id, 'Failed to get experiment schedule %s');
  }
}

export async function upsertSchedule(schedule: UpsertSchedule): Promise<{ created: boolean; schedule: Schedule }> {
  try {
    const response = await executeApiCall({ method: 'POST', path: '/api/experiments/schedules', body: schedule });
    return { created: response.status === 201, schedule: (await response.json()) as Schedule };
  } catch (e) {
    throw abortExecutionWithError(e, 'Failed to save the experiment schedule for %s', schedule.experimentKey);
  }
}

export async function patchSchedule(id: string, patch: PatchSchedule): Promise<Schedule> {
  try {
    const response = await executeApiCall({
      method: 'PATCH',
      path: `/api/experiments/schedules/${encodeURIComponent(id)}`,
      body: patch,
    });
    return (await response.json()) as Schedule;
  } catch (e) {
    throw notFoundOr(e, id, 'Failed to update experiment schedule %s');
  }
}

export async function removeSchedule(id: string): Promise<void> {
  try {
    await executeApiCall({ method: 'DELETE', path: `/api/experiments/schedules/${encodeURIComponent(id)}` });
  } catch (e) {
    throw notFoundOr(e, id, 'Failed to delete experiment schedule %s');
  }
}

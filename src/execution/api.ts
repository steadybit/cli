// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import type { Schemas } from '../api/schemas.ts';
import { ApiError } from '../api/error.ts';
import { executeApiCall } from '../api/http.ts';
import { abortExecution, abortExecutionWithError } from '../errors.ts';

export type Execution = Schemas['ExperimentExecutionAO'];

function notFoundOr(e: unknown, id: number, msg: string): Error {
  if (e instanceof ApiError && e.status === 404) {
    return abortExecution('Experiment run %s not found.', id);
  }
  return abortExecutionWithError(e, msg, id);
}

// Steps, and with them the target executions and their artifacts, are only included
// when asked for. Without them a run looks the same as one that attached nothing.
export async function fetchExecution(id: number): Promise<Execution> {
  try {
    const response = await executeApiCall({
      method: 'GET',
      path: `/api/experiments/executions/${id}`,
      queryParameters: { fields: 'steps' },
    });
    return (await response.json()) as Execution;
  } catch (e) {
    throw notFoundOr(e, id, 'Failed to get experiment run %s');
  }
}

// The platform accepts a cancel with 202 and hands it on to the agents, so the run is
// still stopping when this returns. A 200 means there was nothing left to cancel.
export async function cancelExecution(id: number): Promise<'accepted' | 'already-ended'> {
  try {
    const response = await executeApiCall({ method: 'POST', path: `/api/experiments/executions/${id}/cancel` });
    return response.status === 202 ? 'accepted' : 'already-ended';
  } catch (e) {
    throw notFoundOr(e, id, 'Failed to cancel experiment run %s');
  }
}

export type PropertyOperation = 'set' | 'add';

export async function changeExecutionProperty(
  id: number,
  key: string,
  operation: PropertyOperation,
  value: unknown
): Promise<void> {
  try {
    await executeApiCall({
      method: 'POST',
      path: `/api/experiments/executions/${id}/properties/${encodeURIComponent(key)}/${operation}`,
      body: value,
    });
  } catch (e) {
    throw notFoundOr(e, id, `Failed to ${operation} property ${key} of experiment run %s`);
  }
}

export async function downloadArtifact(id: number, targetExecutionId: string, artifactId: string): Promise<Buffer> {
  try {
    const response = await executeApiCall({
      method: 'GET',
      path: `/api/experiments/executions/${id}/artifacts/${encodeURIComponent(targetExecutionId)}/${encodeURIComponent(artifactId)}`,
      // Artifacts are reports and log archives, which can take far longer than an API
      // response to arrive.
      timeout: 300000,
    });
    return Buffer.from(await response.arrayBuffer());
  } catch (e) {
    if (e instanceof ApiError && e.status === 404) {
      throw abortExecution('Artifact %s of experiment run %s not found.', artifactId, id);
    }
    throw abortExecutionWithError(e, `Failed to download artifact ${artifactId} of experiment run %s`, id);
  }
}

// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { ExecuteResult, ExecutionResult } from './types';

import { abortExecutionWithError } from '../errors';
import { executeApiCall } from '../api/http';

export async function executeExperiment(key: string): Promise<ExecuteResult> {
  try {
    const response = await executeApiCall({
      method: 'POST',
      path: `/api/experiments/${encodeURIComponent(key)}/execute`,
    });
    return { location: response.headers.get('Location') ?? '' };
  } catch (e) {
    const error = await abortExecutionWithError(e, 'Failed to execute task %s (%s@%s)', key);
    throw error;
  }
}

export async function getExperimentExecution(path: string): Promise<ExecutionResult> {
  try {
    const response = await executeApiCall({
      method: 'GET',
      path,
      fullyQualifiedUrl: true,
    });
    return (await response.json()) as ExecutionResult;
  } catch (e) {
    const error = await abortExecutionWithError(e, 'Failed to get experiment execution ');
    throw error;
  }
}

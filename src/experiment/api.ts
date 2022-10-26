// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { ExecuteResult, ExecutionResult, Experiment, UpsertAndExecuteResult, UpsertResult } from './types';

import { abortExecution, abortExecutionWithError } from '../errors';
import { executeApiCall } from '../api/http';

export async function executeExperiment(key: string): Promise<ExecuteResult> {
  try {
    const response = await executeApiCall({
      method: 'POST',
      path: `/api/experiments/${encodeURIComponent(key)}/execute`
    });
    return { location: response.headers.get('Location') ?? '' };
  } catch (e) {
    throw await abortExecutionWithError(e, 'Failed to run task %s (%s@%s)', key);
  }
}

export async function getExperimentExecution(path: string): Promise<ExecutionResult> {
  try {
    const response = await executeApiCall({
      method: 'GET',
      path,
      fullyQualifiedUrl: true
    });
    return (await response.json()) as ExecutionResult;
  } catch (e) {
    throw await abortExecutionWithError(e, 'Failed to get experiment run ');
  }
}


export async function fetchExperiment(key: string): Promise<Experiment> {
  try {
    const response = await executeApiCall({
      method: 'GET',
      path: `/api/experiments/${encodeURIComponent(key)}`
    });
    const experiment = await response.json();
    delete experiment.version; // We remove the version (as this makes things complicated to use). Will be removed from API in the future.
    return experiment;
  } catch (e: any) {
    if (e.response?.status === 404) {
      throw abortExecution('Experiment %s not found.', key);
    } else {
      throw await abortExecutionWithError(e, 'Failed to get the experiment. HTTP request failed.');
    }
  }
}

export async function updateExperiment(key: string, experiment: Experiment): Promise<void> {
  try {
    await executeApiCall({
      method: 'POST',
      path: `/api/experiments/${encodeURIComponent(key)}`,
      body: experiment
    });
  } catch (e: any) {
    if (e.response.status === 404) {
      throw await abortExecutionWithError('Experiment %s not found.', key);
    } else {
      throw await abortExecutionWithError(e, 'Failed to save the experiment. HTTP request failed.');
    }
  }
}

export async function removeExperiment(key: string): Promise<void> {
  try {
    await executeApiCall({
      method: 'DELETE',
      path: `/api/experiments/${encodeURIComponent(key)}`
    });
  } catch (e: any) {
    if (e.response.status === 404) {
      throw abortExecution('Experiment %s not found.', key);
    } else {
      throw await abortExecutionWithError(e, 'Failed to delete the experiment. HTTP request failed.');
    }
  }
}

export async function upsertExperiment(experiment: Experiment): Promise<UpsertResult> {
  try {
    const response = await executeApiCall({
      method: 'POST',
      path: '/api/experiments',
      body: experiment
    });
    const location = response.headers.get('Location');
    const key = location?.substring(location.lastIndexOf('/') + 1);
    return { created: response.status === 201, key };
  } catch (e: any) {
    throw await abortExecutionWithError(e, 'Failed to save the experiment. HTTP request failed.');
  }
}

export async function upsertAndExecuteExperiment(experiment: Experiment): Promise<UpsertAndExecuteResult> {
  try {
    const response = await executeApiCall({
      method: 'POST',
      path: '/api/experiments/execute',
      body: experiment
    });
    const body = await response.json();
    return {
      key: body.key,
      location: response.headers.get('Location') ?? `/api/experiments/executions/${body.executionId}`
    };
  } catch (e: any) {
    throw await abortExecutionWithError(e, 'Failed to save and run the experiment. HTTP request failed.');
  }
}

// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import {
  ExecuteResult,
  ExecutionError,
  ExecutionList,
  ExecutionResult,
  Experiment,
  ExperimentList,
  UpsertAndExecuteResult,
  UpsertResult,
} from './types';

import { abortExecution, abortExecutionWithError, getExecutionErrorBody } from '../errors';
import { executeApiCall } from '../api/http';
import { confirm } from '../prompt/confirm';

export async function executeExperiment(
  key: string,
  yes: boolean,
  allowParallelExecutions: boolean = false
): Promise<ExecuteResult> {
  try {
    const response = await executeApiCall({
      method: 'POST',
      path: `/api/experiments/${encodeURIComponent(key)}/execute?allowParallel=${String(allowParallelExecutions)}`,
    });

    let uiLocation = 'please update your platform to get the UI location';
    const body = await response.text();
    if (body && body.length > 0) {
      const json = JSON.parse(body);
      uiLocation = json.uiLocation;
    }
    return { location: response.headers.get('Location') ?? '', uiLocation };
  } catch (e) {
    if (
      !allowParallelExecutions &&
      (await getExecutionErrorBody<ExecutionError>(e))?.type ===
        'https://steadybit.com/problems/another-experiment-running-exception' &&
      (yes ||
        (await confirm(`There is already an experiment running. Do you want to start ${key} in parallel?`, {
          defaultYes: false,
          defaultWhenNonInteractive: false,
        })))
    ) {
      // try again, but run in parallel
      return executeExperiment(key, yes, true);
    }
    throw await abortExecutionWithError(e, 'Failed to run experiment (%s)', key);
  }
}

export async function upsertAndExecuteExperiment(
  experiment: Experiment,
  allowParallelExecutions: boolean = false
): Promise<UpsertAndExecuteResult> {
  try {
    const response = await executeApiCall({
      method: 'POST',
      path: `/api/experiments/execute?allowParallel=${String(allowParallelExecutions)}`,
      body: experiment,
    });

    let uiLocation = 'please update your platform to get the UI location';
    let key = experiment.key;
    let executionId = undefined;
    const body = await response.text();
    if (body && body.length > 0) {
      const json = JSON.parse(body);
      uiLocation = json.uiLocation;
      key = json.key;
      executionId = json.executionId;
    }
    return {
      key,
      location: response.headers.get('Location') ?? `/api/experiments/executions/${executionId}`,
      uiLocation,
    };
  } catch (e: any) {
    if (
      !allowParallelExecutions &&
      (await getExecutionErrorBody<ExecutionError>(e))?.type ===
        'https://steadybit.com/problems/another-experiment-running-exception' &&
      (await confirm(`There is already an experiment running. Do you want to start ${experiment.key} in parallel?`, {
        defaultYes: false,
        defaultWhenNonInteractive: false,
      }))
    ) {
      // try again, but run in parallel
      return upsertAndExecuteExperiment(experiment, true);
    }
    throw await abortExecutionWithError(e, 'Failed to save and run the experiment. HTTP request failed.');
  }
}

export async function getExperimentExecutionUsingUrl(url: string): Promise<ExecutionResult> {
  try {
    const response = await executeApiCall({
      method: 'GET',
      path: url,
      fullyQualifiedUrl: true,
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
      path: `/api/experiments/${encodeURIComponent(key)}`,
    });
    const experiment = (await response.json()) as Experiment;
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
      body: experiment,
    });
  } catch (e: any) {
    if (e.response.status === 404) {
      throw abortExecution('Experiment %s not found.', key);
    } else {
      throw await abortExecutionWithError(e, 'Failed to save the experiment. HTTP request failed.');
    }
  }
}

export async function removeExperiment(key: string): Promise<void> {
  try {
    await executeApiCall({
      method: 'DELETE',
      path: `/api/experiments/${encodeURIComponent(key)}`,
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
      body: experiment,
    });
    const location = response.headers.get('Location');
    const key = location?.substring(location.lastIndexOf('/') + 1);
    return { created: response.status === 201, key };
  } catch (e: any) {
    throw await abortExecutionWithError(e, 'Failed to save the experiment. HTTP request failed.');
  }
}

export async function fetchExperiments(teamKey: string): Promise<ExperimentList> {
  try {
    const response = await executeApiCall({
      method: 'GET',
      path: '/api/experiments',
      queryParameters: {
        team: teamKey,
      },
    });
    return (await response.json()) as ExperimentList;
  } catch (e: any) {
    throw await abortExecutionWithError(e, 'Failed to get the experiments. HTTP request failed.');
  }
}

export async function fetchExecutionsForExperiment(key: string): Promise<ExecutionList> {
  try {
    const response = await executeApiCall({
      method: 'GET',
      path: `/api/experiments/${encodeURIComponent(key)}/executions`,
    });
    return (await response.json()) as ExecutionList;
  } catch (e: any) {
    throw await abortExecutionWithError(e, 'Failed to get the executions. HTTP request failed.');
  }
}

export async function getExperimentExecution(id: number, abortOnError = true): Promise<ExecutionResult> {
  try {
    const response = await executeApiCall({
      method: 'GET',
      path: `/api/experiments/executions/${id}`,
    });
    return (await response.json()) as ExecutionResult;
  } catch (e) {
    if (abortOnError) {
      throw await abortExecutionWithError(e, 'Failed to get experiment run ');
    } else {
      throw e;
    }
  }
}

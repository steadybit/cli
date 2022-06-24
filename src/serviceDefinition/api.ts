// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { ServiceDefinition, ServiceState, Task, TaskHistoryEntry } from './types';

import { abortExecutionWithError } from '../errors';
import { executeApiCall } from '../api/http';

export async function getState(serviceDefinition: ServiceDefinition): Promise<ServiceState> {
  try {
    const response = await executeApiCall({
      method: 'GET',
      path: `/api/service-states/${encodeURIComponent(serviceDefinition.id)}`,
    });
    return (await response.json()) as ServiceState;
  } catch (e) {
    const error = await abortExecutionWithError(e, 'Failed to read state for service %s', serviceDefinition.id);
    throw error;
  }
}

export async function executeTask(serviceState: ServiceState, task: Task): Promise<TaskHistoryEntry> {
  try {
    const response = await executeApiCall({
      method: 'POST',
      path: `/api/service-states/${encodeURIComponent(serviceState.id)}/task/${encodeURIComponent(task.id)}/execute`,
    });
    return (await response.json()) as TaskHistoryEntry;
  } catch (e) {
    const error = await abortExecutionWithError(
      e,
      'Failed to run task %s (%s@%s)',
      task.id,
      task.definition.name,
      task.definition.version
    );
    throw error;
  }
}

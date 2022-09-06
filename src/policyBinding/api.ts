// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { PolicyBinding, PolicyBindingState, Task, TaskHistoryEntry } from './types';

import { abortExecutionWithError } from '../errors';
import { executeApiCall } from '../api/http';

export async function getState(policyBinding: PolicyBinding): Promise<PolicyBindingState> {
  try {
    const response = await executeApiCall({
      method: 'GET',
      path: `/api/policy-bindings/state/${encodeURIComponent(policyBinding.id)}`,
    });
    return (await response.json()) as PolicyBindingState;
  } catch (e) {
    const error = await abortExecutionWithError(e, 'Failed to read state for policy binding %s', policyBinding.id);
    throw error;
  }
}

export async function executeTask(policyBindingState: PolicyBindingState, task: Task): Promise<TaskHistoryEntry> {
  try {
    const response = await executeApiCall({
      method: 'POST',
      path: `/api/policy-bindings/state/${encodeURIComponent(policyBindingState.id)}/task/${encodeURIComponent(task.id)}/execute`,
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

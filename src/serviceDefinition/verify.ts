import { red, bold, green, blue } from 'colors/safe';

import { TaskState, ResilienceScoreServiceState } from './types';
import { abortExecutionWithError } from '../errors';
import { loadServiceDefinition } from './loading';
import { executeApiCall } from '../api';

const taskSuffix: Record<TaskState, string> = {
  'SUCCESS': 'ok',
  'FAILURE': 'failure',
  'PENDING': 'pending'
};

const taskColor: Record<TaskState, typeof red> = {
  'SUCCESS': green,
  'FAILURE': red,
  'PENDING': blue
};

const taskOrder: Record<TaskState, number> = {
  'SUCCESS': 0,
  'PENDING': 1,
  'FAILURE': 2
};

export async function verify(serviceDefinitionPath: string) {
  const serviceDefinition = await loadServiceDefinition(serviceDefinitionPath);

  try {
    const response = await executeApiCall({
      method: 'GET',
      path: `/api/resilience-score/service/${serviceDefinition.id}`,
    });
    const state = (await response.json()) as ResilienceScoreServiceState;

    printResilienceLevel(state);
    console.log();
    printTaskList(state);

    if (!hasAchievedDesiredLevel(state)) {
      process.exit(1);
    }
  } catch (e) {
    throw abortExecutionWithError(e, 'Failed to read state for service %s', serviceDefinition.id);
  }
}

function printResilienceLevel(state: ResilienceScoreServiceState) {
  let color: typeof red = s => s;
  if (!hasAchievedDesiredLevel(state)) {
    color = red;
  }

  console.log(color(bold(`Resilience Level`)));
  console.log(color(bold('================')));
  console.log('Desired: %s', state.desiredResilienceLevel);
  console.log(color('Actual:  %s'), state.actualResilienceLevel);
}

function printTaskList(state: ResilienceScoreServiceState) {
  console.log(bold(`Tasks`));
  console.log(bold('====='));

  const sortedTasks = state.tasks
    .slice()
    .sort((a, b) => taskOrder[a.state] - taskOrder[b.state]);

  const countByType: Record<TaskState, number> = {
    'SUCCESS': 0,
    'FAILURE': 0,
    'PENDING': 0
  };

  for (const task of sortedTasks) {
    countByType[task.state]++;
    console.log(taskColor[task.state](` - %s (%s)`), task.name, taskSuffix[task.state]);
  }

  console.log();
  console.log(taskColor.SUCCESS(`Ok:      %d`), countByType.SUCCESS);
  console.log(taskColor.PENDING(`Pending: %d`), countByType.PENDING);
  console.log(taskColor.FAILURE(`Failure: %d`), countByType.FAILURE);
}

function hasAchievedDesiredLevel(state: ResilienceScoreServiceState) {
  return state.desiredResilienceLevel === state.actualResilienceLevel;
}

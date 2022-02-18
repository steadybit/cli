import { red, bold, green, blue, gray } from 'colors/safe';

import { TaskState, ServiceState } from './types';
import { abortExecutionWithError } from '../errors';
import { loadServiceDefinition } from './loading';
import { executeApiCall } from '../api';

const taskSuffix: Record<TaskState, string> = {
  SUCCESS: 'ok',
  FAILURE: 'failure',
  PENDING: 'pending',
};

const taskColor: Record<TaskState, typeof red> = {
  SUCCESS: green,
  FAILURE: red,
  PENDING: blue,
};

const taskOrder: Record<TaskState, number> = {
  SUCCESS: 0,
  PENDING: 1,
  FAILURE: 2,
};

interface Options {
  printParameters?: boolean;
  printMatrixContext?: boolean;
}

export async function verify(serviceDefinitionPath: string, options: Options) {
  const serviceDefinition = await loadServiceDefinition(serviceDefinitionPath);

  try {
    const response = await executeApiCall({
      method: 'GET',
      path: `/api/service-states/${serviceDefinition.id}`,
    });
    const state = (await response.json()) as ServiceState;

    printTaskList(options, state);
  } catch (e) {
    const error = await abortExecutionWithError(e, 'Failed to read state for service %s', serviceDefinition.id);
    throw error;
  }
}

function printTaskList(options: Options, state: ServiceState) {
  console.log(bold(`Tasks`));
  console.log(bold('====='));

  const sortedTasks = state.tasks.slice().sort((a, b) => {
    let result = taskOrder[a.state] - taskOrder[b.state];
    if (result === 0) {
      result = a.definition.name.localeCompare(b.definition.name);
    }
    return result;
  });

  const countByType: Record<TaskState, number> = {
    SUCCESS: 0,
    FAILURE: 0,
    PENDING: 0,
  };

  for (const task of sortedTasks) {
    countByType[task.state]++;
    console.log(
      taskColor[task.state](` - %s (%s)`),
      `${task.definition.name}@${task.definition.version}`,
      taskSuffix[task.state]
    );

    if (options.printParameters) {
      printJson('Parameters', task.parameters);
    }

    if (options.printMatrixContext || Object.keys(task.matrixContext).length > 0) {
      printJson('Matrix Context', task.matrixContext);
    }
  }

  console.log();
  console.log(taskColor.SUCCESS(`Ok:      %d`), countByType.SUCCESS);
  console.log(taskColor.PENDING(`Pending: %d`), countByType.PENDING);
  console.log(taskColor.FAILURE(`Failure: %d`), countByType.FAILURE);
}

function printJson(title: string, obj: never): void {
  const content = `${title}:\n${JSON.stringify(obj, undefined, 2)}\n`
    .split('\n')
    .map(s => '   ' + s)
    .join('\n');

  console.log(gray(content));
}

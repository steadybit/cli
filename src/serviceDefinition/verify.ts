// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { blue, bold, gray, green, red } from 'colors/safe';
import { ServiceState, Task, TaskState } from './types';

import { loadServiceDefinition } from './files';
import { getState } from './api';

const taskSuffix: Record<TaskState, string> = {
  SUCCESS: 'ok',
  FAILURE: 'failure',
  PENDING: 'not checked',
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
  file: string;
  printParameters?: boolean;
  printMatrixContext?: boolean;
}

const openHelp = gray(
  `
Do you need more information? You can see more details within Steadybit's
user interface. Try executing the following to open it:

                   ${bold('steadybit service open')}
`.trim()
);

export async function verify(options: Options) {
  const serviceDefinition = await loadServiceDefinition(options.file);
  const state = await getState(serviceDefinition);

  printTaskList(options, state);

  console.log();
  console.log(openHelp);

  const countByType = getTaskCountByType(state);
  if (countByType.FAILURE > 0 || countByType.PENDING > 0) {
    process.exit(1);
  }
}

function getTaskCountByType(state: ServiceState): Record<TaskState, number> {
  const countByType: Record<TaskState, number> = {
    SUCCESS: 0,
    FAILURE: 0,
    PENDING: 0,
  };

  for (const task of state.tasks) {
    countByType[task.state]++;
  }

  return countByType;
}

function printTaskList(options: Options, state: ServiceState) {
  console.log(bold('Tasks'));
  console.log(bold('====='));

  const sortedTasks = state.tasks.slice().sort((a, b) => {
    let result = taskOrder[a.state] - taskOrder[b.state];
    if (result === 0) {
      result = a.definition.name.localeCompare(b.definition.name);
    }
    return result;
  });

  const countByType = getTaskCountByType(state);

  const taskCountByCoordinate: Record<string, number> = {};
  for (const task of sortedTasks) {
    const key = getCoordinateKey(task);
    taskCountByCoordinate[key] = (taskCountByCoordinate[key] ?? -1) + 1;
  }

  for (const task of sortedTasks) {
    const key = getCoordinateKey(task);

    console.log(taskColor[task.state](' - %s (%s)'), key, taskSuffix[task.state]);

    if (options.printParameters) {
      printJson('Parameters', task.parameters);
    }

    if (options.printMatrixContext || taskCountByCoordinate[key] > 1) {
      printJson('Matrix Context', task.matrixContext);
    }
  }

  console.log();
  console.log(taskColor.SUCCESS('Ok:          %d'), countByType.SUCCESS);
  console.log(taskColor.PENDING('Not Checked: %d'), countByType.PENDING);
  console.log(taskColor.FAILURE('Failure:     %d'), countByType.FAILURE);
}

function getCoordinateKey(task: Task): string {
  return `${task.definition.name}@${task.definition.version}`;
}

function printJson(title: string, obj: unknown): void {
  const content = `${title}:\n${JSON.stringify(obj, undefined, 2)}\n`
    .split('\n')
    .map(s => '   ' + s)
    .join('\n');

  console.log(gray(content));
}

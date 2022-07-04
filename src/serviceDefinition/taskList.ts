// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { blue, gray, green, red } from 'colors/safe';
import { printTitle } from '../print/title';
import { Task, TaskState } from './types';
import { indent } from '../print/indent';
import yaml from 'js-yaml';

const taskSuffix: Record<TaskState, string> = {
  SUCCESS: 'ok',
  FAILURE: 'failure',
  PENDING: 'unchecked',
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

export interface PrintTaskListOptions {
  printParameters?: boolean;
  printMatrixContext?: boolean;
  colorsByTaskState?: boolean;
  renderTaskState?: boolean;
  title?: string;
}

export function printTaskList(options: PrintTaskListOptions, tasks: Task[]) {
  printTitle(options.title ?? 'Tasks');

  const sortedTasks = tasks.slice().sort((a, b) => {
    let result = taskOrder[a.state] - taskOrder[b.state];
    if (result === 0) {
      result = a.definition.name.localeCompare(b.definition.name);
    }
    return result;
  });

  const taskCountByCoordinate = getTaskCountByCoordinate(tasks);

  for (const task of sortedTasks) {
    const key = getCoordinateKey(task);
    const addColor = options.colorsByTaskState ? taskColor[task.state] : (s: string) => s;
    console.log(addColor(' - %s%s'), key, options.renderTaskState ? ` (${taskSuffix[task.state]})` : '');

    if (options.printParameters) {
      printObj('Parameters', task.parameters);
    }

    if (options.printMatrixContext || taskCountByCoordinate[key] > 1) {
      printObj('Matrix Context', task.matrixContext);
    }
  }
}

export function getTaskCountByType(tasks: Task[]): Record<TaskState, number> {
  const countByType: Record<TaskState, number> = {
    SUCCESS: 0,
    FAILURE: 0,
    PENDING: 0,
  };

  for (const task of tasks) {
    countByType[task.state]++;
  }

  return countByType;
}

export function printTaskStateByType(countByType: Record<TaskState, number>): void {
  console.log(taskColor.SUCCESS('Ok:          %d'), countByType.SUCCESS);
  if (countByType.PENDING >= 0) {
    console.log(taskColor.PENDING('Not Checked: %d'), countByType.PENDING);
  }
  console.log(taskColor.FAILURE('Failure:     %d'), countByType.FAILURE);
}

export function getCoordinateKey(task: Task): string {
  return `${task.definition.name}@${task.definition.version}`;
}

function printObj(title: string, obj: unknown): void {
  const content = indent(`${title}:`, 3) + '\n' + indent(yaml.dump(obj), 5) + '\n';
  console.log(gray(content));
}

export function getTaskCountByCoordinate(tasks: Task[]): Record<string, number> {
  const taskCountByCoordinate: Record<string, number> = {};
  for (const task of tasks) {
    const key = getCoordinateKey(task);
    taskCountByCoordinate[key] = (taskCountByCoordinate[key] ?? -1) + 1;
  }
  return taskCountByCoordinate;
}

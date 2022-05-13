// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { ServiceDefinition, ServiceState, Task, TaskHistoryEntry, TaskState } from '../types';
import { from, switchMap, map, takeWhile, lastValueFrom, distinct } from 'rxjs';
import { getCoordinateKey, getTaskCountByCoordinate, printTaskStateByType } from '../taskList';
import { executeTask as sendExecuteTaskHttpRequest } from '../api';
import { toTwoColumns } from '../../print/columns';
import { executeApiCall } from '../../api/sse';
import { gray, green, red } from 'colors/safe';
import { printTitle } from '../../print/title';
import { indent } from '../../print/indent';
import { toUrl } from '../../api/common';
import { Options } from './types';
import yaml from 'js-yaml';
import ora from 'ora';

const columnLength = 20;

export async function execute(
  serviceDefinition: ServiceDefinition,
  serviceState: ServiceState,
  tasks: Task[],
  options: Options
): Promise<boolean | undefined> {
  const taskCountByCoordinate = getTaskCountByCoordinate(tasks);

  const countByType: Record<TaskState, number> = {
    SUCCESS: 0,
    FAILURE: 0,
    PENDING: -1,
  };

  let allSuccessful: boolean | undefined = true;
  let i = 0;
  for (const task of tasks) {
    const start = Date.now();

    printTitle(`Executing task ${++i} of ${tasks.length}`);

    console.log(toTwoColumns('Name', columnLength, task.definition.name));
    console.log(toTwoColumns('Version', columnLength, task.definition.version));
    console.log(toTwoColumns('Type', columnLength, task.type));
    if (options.printParameters) {
      console.log('Parameters:');
      console.log(indent(yaml.dump(task.parameters), 2));
    }
    if (options.printMatrixContext || taskCountByCoordinate[getCoordinateKey(task)] > 1) {
      console.log('Matrix Context:');
      console.log(indent(yaml.dump(task.matrixContext), 2));
    }

    const result = await executeSingle(serviceState, task, options);

    if (result !== undefined) {
      const durationMillis = Date.now() - start;
      let message: string;
      if (result === true) {
        message = green('Execution completed successfully');
        countByType['SUCCESS']++;
      } else {
        message = red('Execution completed with failures');
        countByType['FAILURE']++;
      }
      console.log(gray('\n%s (took %s ms)'), message, durationMillis);
    }
    console.log();

    // If at least one result is undefined, we don't know if all tasks were successful or failed.
    // We therefore need to retain that state.
    if (result === undefined) {
      allSuccessful = undefined;
    } else if (allSuccessful !== undefined) {
      allSuccessful = allSuccessful && result;
    }
  }

  printTitle('Summary');
  printTaskStateByType(countByType);

  return allSuccessful;
}

async function executeSingle(serviceState: ServiceState, task: Task, options: Options): Promise<boolean | undefined> {
  const spinner = ora('Executing task...').start();
  try {
    if (!options.wait) {
      const taskHistoryEntry = await sendExecuteTaskHttpRequest(serviceState, task);
      await printExecutionStart(taskHistoryEntry);
      return undefined;
    }

    return await lastValueFrom(
      executeApiCall<TaskHistoryEntry>({
        path: `/api/service-states/${serviceState.id}/task/${task.id}/execute/wait`,
      })
        .pipe(distinct(e => e.type))
        .pipe(switchMap(e => from(informAboutHistoryEntry(e, spinner)).pipe(map(() => e))))
        .pipe(takeWhile(e => e.type !== 'EXECUTION_COMPLETE', true))
        .pipe(map(e => e.meta.executionState === 'SUCCESS'))
    );
  } finally {
    spinner.stop();
  }
}

async function informAboutHistoryEntry(entry: TaskHistoryEntry, spinner: ora.Ora): Promise<void> {
  spinner.stop();
  if (entry.type === 'EXECUTION_REQUESTED') {
    await printExecutionStart(entry);
  }
  spinner.start();
}

async function printExecutionStart(entry: TaskHistoryEntry): Promise<void> {
  if (entry.meta.executionId) {
    console.log(toTwoColumns('Execution ID', columnLength, entry.meta.executionId));
  }
  if (entry.meta.url) {
    const url = await toUrl(entry.meta.url);
    console.log(toTwoColumns('URL', columnLength, url));
  }
}

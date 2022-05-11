// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Options } from './types';
import { abortExecution } from '../../errors';
import colors from 'colors/safe';
import { confirm } from '../../prompt/confirm';
import { execute } from './execution';
import { getState } from '../api';
import { getTasks } from './taskIdentification';
import { loadServiceDefinition } from '../files';
import { printTaskList } from '../taskList';

export async function exec(options: Options): Promise<void> {
  const serviceDefinition = await loadServiceDefinition(options.file);

  const serviceState = await getState(serviceDefinition);
  const tasks = await getTasks(options, serviceState);

  if (tasks.length === 0) {
    console.log(colors.red(`No tasks found for service ${colors.bold(serviceDefinition.name)}. Nothing to execute.`));
    process.exit(options.errorOnEmptyTaskSet ? 1 : 0);
  } else if (!options.wait && tasks.length > 1) {
    throw abortExecution(
      'Cannot execute multiple tasks with --no-wait. Please remove --no-wait or execute tasks one by one.'
    );
  }

  printTaskList(
    {
      ...options,
      title: 'Tasks to execute',
    },
    tasks
  );
  console.log();

  if (!options.yes) {
    const confirmation = await confirm('Are you sure you want to execute these tasks?', {
      defaultYes: false,
      defaultWhenNonInteractive: true,
    });
    if (!confirmation) {
      process.exit(0);
    }
  }

  const allSuccessful = await execute(serviceDefinition, serviceState, tasks, options);
  // beware: allSuccessful may be undefined when using --no-wait. We therefore have to check
  // for === false explicitly.
  process.exit(allSuccessful === false && options.errorOnTaskFailure ? 1 : 0);
}

// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { loadServiceDefinition } from '../files';
import { getTasks } from './taskIdentification';
import { printSummary } from './taskSummary';
import { Options } from './types';
import colors from 'colors/safe';
import { confirm } from '../../prompt/confirm';
import { execute } from './execution';

export async function exec(options: Options): Promise<void> {
  console.log({ options });

  const serviceDefinition = await loadServiceDefinition(options.file);

  const tasks = await getTasks(options, serviceDefinition);
  if (tasks.length === 0) {
    console.log(colors.red(`No tasks found for service ${colors.bold(serviceDefinition.name)}. Nothing to execute.`));
    process.exit(options.errorOnEmptyTaskSet ? 1 : 0);
  }

  printSummary(serviceDefinition, tasks);

  if (!options.yes) {
    const confirmation = await confirm('Are you sure you want to execute these tasks?', { defaultYes: false, defaultWhenNonInteractive: true });
    if (!confirmation) {
      process.exit(0);
    }
  }

  const allSuccessful = await execute(serviceDefinition, tasks, options);
  process.exit(!allSuccessful && options.errorOnTaskFailure ? 1 : 0);
}

// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Options } from './types';
import { abortExecution } from '../../errors';
import colors from 'colors/safe';
import { confirm } from '../../prompt/confirm';
import { execute } from './execution';
import { getState } from '../api';
import { getTasks } from './taskIdentification';
import { loadPolicyBinding } from '../files';
import { printTaskList } from '../taskList';

export async function exec(options: Options): Promise<void> {
  const policyBinding = await loadPolicyBinding(options.file);

  const policyBindingState = await getState(policyBinding);
  const tasks = await getTasks(options, policyBindingState);

  if (tasks.length === 0) {
    console.log(colors.red(`No tasks found for policy-binding ${colors.bold(policyBinding.name)}. Nothing to run.`));
    process.exit(options.errorOnEmptyTaskSet ? 1 : 0);
  } else if (!options.wait && tasks.length > 1) {
    throw abortExecution(
      'Cannot run multiple tasks with --no-wait. Please remove --no-wait or run tasks one by one.'
    );
  }

  printTaskList(
    {
      ...options,
      title: 'Tasks to run',
    },
    tasks
  );
  console.log();

  if (!options.yes) {
    const confirmation = await confirm('Are you sure you want to run these tasks?', {
      defaultYes: false,
      defaultWhenNonInteractive: true,
    });
    if (!confirmation) {
      process.exit(0);
    }
  }

  const allSuccessful = await execute(policyBinding, policyBindingState, tasks, options);
  // beware: allSuccessful may be undefined when using --no-wait. We therefore have to check
  // for === false explicitly.
  process.exit(allSuccessful === false && options.errorOnTaskFailure ? 1 : 0);
}

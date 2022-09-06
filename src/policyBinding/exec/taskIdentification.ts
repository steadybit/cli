// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { PolicyBindingState, Task } from '../types';
import { Options } from './types';

export function getTasks(options: Options, policyBindingState: PolicyBindingState): Task[] {
  let tasks = policyBindingState.tasks;

  const taskFilter = options.task;
  if (taskFilter != null) {
    tasks = tasks.filter(task => taskFilter.includes(task.definition.name));
  }

  return tasks;
}

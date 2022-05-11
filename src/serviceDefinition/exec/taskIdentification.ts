// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { ServiceState, Task } from '../types';
import { Options } from './types';

export function getTasks(options: Options, serviceState: ServiceState): Task[] {
  let tasks = serviceState.tasks;

  const taskFilter = options.task;
  if (taskFilter != null) {
    tasks = tasks.filter(task => taskFilter.includes(task.definition.name));
  }

  return tasks;
}

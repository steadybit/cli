// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { ServiceDefinition, Task } from '../types';
import { Options } from './types';
import { getState } from '../api';

export async function getTasks(options: Options, serviceDefinition: ServiceDefinition): Promise<Task[]> {
  const state = await getState(serviceDefinition);
  let tasks = state.tasks;

  const taskFilter = options.task;
  if (taskFilter != null) {
    tasks = tasks.filter(task => taskFilter.includes(task.definition.name));
  }

  return tasks;
}

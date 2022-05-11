// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { ServiceDefinition, Task } from '../types';
import colors from 'colors/safe';

export function printSummary(serviceDefinition: ServiceDefinition, tasks: Task[]): void {
  console.log(getSummary(serviceDefinition, tasks));
}

export function getSummary(serviceDefinition: ServiceDefinition, tasks: Task[]): string {
  // TODO
  return `${tasks.length} tasks found for ${colors.bold(serviceDefinition.name)}.`;
}

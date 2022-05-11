// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { ServiceDefinition, Task } from '../types';
import { Options } from './types';

export async function execute(serviceDefinition: ServiceDefinition, tasks: Task[], options: Options): Promise<boolean> {
  return true;
}

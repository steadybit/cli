// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { ServiceDefinition, Task } from '../types';
import { getSummary } from './taskSummary';

describe('taskSummary', () => {
  it('should print a summary', () => {
    const serviceDefinition = {
      id: 'test',
      name: 'gateway'
    } as ServiceDefinition;

    const tasks = [
      {
        definition: {
          name: 'task1',
          version: '1.0.0',
        },
      } as Task,
      {
        definition: {
          name: 'task2',
          version: '2.0.0',
        },
      } as Task,
    ];

    expect(getSummary(serviceDefinition, tasks)).toMatchInlineSnapshot(`"2 tasks found for [1mgateway[22m."`);
  });
});

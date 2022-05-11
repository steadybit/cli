// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { ServiceDefinition, Task } from '../types';
import { getTasks } from './taskIdentification';
import { getState } from '../api';
import { Options } from './types';

const getStateMock = getState as jest.Mock;

jest.mock('../api', () => ({
  getState: jest.fn(),
}));

const serviceDefinition = {
  id: 'test',
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

describe('taskIdentification', () => {
  beforeEach(() => {
    getStateMock.mockReset();
  });

  describe('getTasks', () => {
    it('must not filter tasks', async () => {
      // Given
      getStateMock.mockResolvedValue({
        tasks,
      });

      // When
      const result = await getTasks({} as Options, serviceDefinition);

      // Then
      expect(result).toEqual(tasks);
    });

    it('must filter tasks when an empty array is passed', async () => {
      // Given
      getStateMock.mockResolvedValue({
        tasks,
      });

      // When
      const result = await getTasks(
        {
          task: [],
        } as Options,
        serviceDefinition
      );

      // Then
      expect(result).toEqual([]);
    });

    it('must filter tasks by name', async () => {
      // Given
      getStateMock.mockResolvedValue({
        tasks,
      });

      // When
      const result = await getTasks(
        {
          task: ['task2'],
        } as Options,
        serviceDefinition
      );

      // Then
      expect(result).toEqual([tasks[1]]);
    });
  });
});

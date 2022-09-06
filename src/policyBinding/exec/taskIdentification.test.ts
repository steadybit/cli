// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { PolicyBindingState, Task } from '../types';
import { getTasks } from './taskIdentification';
import { getState } from '../api';
import { Options } from './types';

const getStateMock = getState as jest.Mock;

jest.mock('../api', () => ({
  getState: jest.fn(),
}));

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

      // When
      const result = getTasks({} as Options, { tasks } as PolicyBindingState);

      // Then
      expect(result).toEqual(tasks);
    });

    it('must filter tasks when an empty array is passed', async () => {
      // Given

      // When
      const result = getTasks(
        {
          task: [],
        } as Options,
        { tasks } as PolicyBindingState
      );

      // Then
      expect(result).toEqual([]);
    });

    it('must filter tasks by name', async () => {
      // Given

      // When
      const result = getTasks(
        {
          task: ['task2'],
        } as Options,
        { tasks } as PolicyBindingState
      );

      // Then
      expect(result).toEqual([tasks[1]]);
    });
  });
});

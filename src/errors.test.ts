// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { describe, expect, it } from 'vitest';
import { executeApiCall } from './api/http.ts';
import { getExecutionErrorBody } from './errors.ts';

interface ExecutionProblem {
  type: string;
}

describe('errors', () => {
  describe('getExecutionErrorBody', () => {
    // executeApiCall consumes the response body to build its error message, so the body
    // has to be carried on the error itself for this to return anything at all.
    it('should expose the problem body of a failed API call', async () => {
      const error = await executeApiCall({ method: 'GET', path: '/api/problem' }).catch(e => e);

      expect(getExecutionErrorBody<ExecutionProblem>(error)).toEqual({
        type: 'https://steadybit.com/problems/another-experiment-running-exception',
      });
    });

    it('should return undefined when there is no body to parse', async () => {
      expect(getExecutionErrorBody(new Error('boom'))).toBeUndefined();
    });
  });
});

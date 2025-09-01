// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { executeApiCall, options } from './http';

describe('http', () => {
  beforeAll(() => {
    options.defaultWaitTime = 10;
  });

  afterAll(() => {
    options.defaultWaitTime = 1000;
  });

  describe('too many requests', () => {
    it('should not handle codes besides Too Many Requests', async () => {
      await expect(() =>
        executeApiCall({
          method: 'GET',
          path: `/api/status`,
          queryParameters: {
            code: '500',
            body: 'Internal Server Error',
          },
        })
      ).rejects.toThrow('responded with unexpected status code: 500 - Internal Server Error');
    });

    it('should retry on too many requests response', async () => {
      const response = await executeApiCall({
        method: 'GET',
        path: `/api/status`,
        queryParameters: {
          code: '429',
          times: '3',
        },
      });
      expect(response.status).toEqual(200);
    });

    it('should throw last error if max retries is exceeded', async () => {
      await expect(() =>
        executeApiCall({
          method: 'GET',
          path: `/api/status`,
          queryParameters: {
            code: '429',
            times: '4',
          },
        })
      ).rejects.toThrow('responded with unexpected status code: 429');
    });

    it('should wait for reset time', async () => {
      const start = new Date();
      const response = await executeApiCall({
        method: 'GET',
        path: `/api/status`,
        queryParameters: {
          code: '429',
          times: '1',
          reset: '2',
        },
      });
      const duration = new Date().getTime() - start.getTime();
      expect(response.status).toEqual(200);
      expect(duration).toBeGreaterThan(2000);
    });
  });
});

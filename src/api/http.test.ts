// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { executeApiCall, options } from './http.ts';

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

  // Absolute URLs arrive from platform responses, most notably the Location header of a
  // started run, and every request carries the API access token.
  describe('absolute urls', () => {
    it('should send a foreign origin to the configured platform instead', async () => {
      // The message carries the URL that was actually requested, so it shows both that
      // the path survived and that nothing was sent to the other host.
      await expect(() =>
        executeApiCall({
          method: 'GET',
          path: 'https://attacker.example/api/status?code=500&body=served%20by%20the%20mock',
        })
      ).rejects.toThrow(
        'Steadybit API at GET http://example.com/api/status?code=500&body=served%20by%20the%20mock responded with unexpected status code: 500 - served by the mock'
      );
    });

    it('should accept an absolute url on the configured origin', async () => {
      const response = await executeApiCall({
        method: 'GET',
        path: 'http://example.com/api/status',
      });

      expect(response.status).toEqual(200);
    });

    it('should keep query parameters when the path is absolute', async () => {
      await expect(() =>
        executeApiCall({
          method: 'GET',
          path: 'http://example.com/api/status',
          queryParameters: {
            code: '500',
            body: 'Internal Server Error',
          },
        })
      ).rejects.toThrow('responded with unexpected status code: 500 - Internal Server Error');
    });
  });
});

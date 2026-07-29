// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { afterAll, beforeAll, describe, expect, it } from 'vitest';
import { delay, http, HttpResponse } from 'msw';
import { server } from '../mocks/server.ts';
import { executeApiCall, options } from './http.ts';

describe('http', () => {
  beforeAll(() => {
    options.defaultWaitTime = 10;
    options.rateLimitBudget = 5000;
  });

  afterAll(() => {
    options.defaultWaitTime = 1000;
    options.rateLimitBudget = 120000;
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

    it('should keep retrying past the old four-attempt cap', async () => {
      const response = await executeApiCall({
        method: 'GET',
        path: `/api/status`,
        queryParameters: {
          code: '429',
          times: '20',
        },
      });
      expect(response.status).toEqual(200);
    });

    // The budget counts time spent waiting on 429s, not elapsed time. Measuring elapsed
    // time meant the pacing before a request could spend the whole budget, leaving none
    // for the rate limit the pacing exists to survive.
    it('should not let a slow response consume the budget', async () => {
      const budget = options.rateLimitBudget;
      options.rateLimitBudget = 100;
      let attempts = 0;
      server.use(
        http.get('http://example.com/api/slow-429', async () => {
          attempts++;
          // Each response takes far longer than the whole budget.
          await delay(300);
          return attempts <= 2 ? new HttpResponse(null, { status: 429 }) : HttpResponse.text('recovered');
        })
      );

      try {
        const response = await executeApiCall({ method: 'GET', path: '/api/slow-429' });
        expect(await response.text()).toEqual('recovered');
        expect(attempts).toEqual(3);
      } finally {
        options.rateLimitBudget = budget;
      }
    });

    it('should surface the rate limit once the budget is spent', async () => {
      const budget = options.rateLimitBudget;
      options.rateLimitBudget = 25;
      try {
        await expect(() =>
          executeApiCall({
            method: 'GET',
            path: `/api/status`,
            queryParameters: {
              code: '429',
              times: '10000',
            },
          })
        ).rejects.toThrow('responded with unexpected status code: 429');
      } finally {
        options.rateLimitBudget = budget;
      }
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

  // A dump makes hundreds of requests, so a single flaky DNS lookup or reset connection
  // must not end the command.
  describe('transport failures', () => {
    it('should retry an idempotent request that fails in transit', async () => {
      let attempts = 0;
      server.use(
        http.get('http://example.com/api/flaky', () => {
          attempts++;
          return attempts === 1 ? HttpResponse.error() : HttpResponse.text('recovered');
        })
      );

      const response = await executeApiCall({ method: 'GET', path: '/api/flaky' });

      expect(await response.text()).toEqual('recovered');
      expect(attempts).toEqual(2);
    });

    it('should not repeat a POST, which may already have started a run', async () => {
      let attempts = 0;
      server.use(
        http.post('http://example.com/api/flaky', () => {
          attempts++;
          return HttpResponse.error();
        })
      );

      await expect(() => executeApiCall({ method: 'POST', path: '/api/flaky' })).rejects.toThrow(
        'Failed to call Steadybit API at POST'
      );
      expect(attempts).toEqual(1);
    });

    it('should give up on an idempotent request once the attempts are used', async () => {
      let attempts = 0;
      server.use(
        http.get('http://example.com/api/always-broken', () => {
          attempts++;
          return HttpResponse.error();
        })
      );

      await expect(() => executeApiCall({ method: 'GET', path: '/api/always-broken' })).rejects.toThrow(
        'Failed to call Steadybit API at GET'
      );
      expect(attempts).toEqual(options.maxRetries + 2);
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

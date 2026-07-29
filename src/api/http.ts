// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { setTimeout as sleep } from 'node:timers/promises';
import { getHeaders, toUrl } from './common.ts';
import { ApiError } from './error.ts';
import { errorMessage } from '../errors.ts';
import { rateLimiter } from './rateLimit.ts';
import { ensurePlatformAccessConfigurationIsAvailable } from '../config/requirePlatformAccess.ts';

const TOO_MANY_REQUESTS = 429;

export const options = {
  maxRetries: 2,
  defaultWaitTime: 1000,
  rateLimitBudget: 120000,
};

export interface ApiCallArguments {
  path: string;
  method: string;
  queryParameters?: Record<string, string>;
  body?: unknown;
  timeout?: number; // defaults to 30000
}

export function enableRequestLogging() {
  process.env.REQUEST_LOGGING_ENABLED = 'true';
}

async function doFetch(
  url: string,
  method: string,
  headers: Record<string, string>,
  body: undefined | string,
  signal: AbortSignal
) {
  if (process.env.REQUEST_LOGGING_ENABLED === 'true') {
    console.log(`> HTTP ${method} ${url}`);
    for (const [key, value] of Object.entries(headers)) {
      console.log(`> ${key}: ${maskSensitiveHeader(key, value)}`);
    }
    console.log(`> `);
    if (body) {
      console.log(body);
    }
    console.log('');
  }

  const response = await fetch(url, {
    method,
    headers,
    body,
    signal,
    redirect: 'error',
  });

  if (process.env.REQUEST_LOGGING_ENABLED === 'true') {
    console.log(`< HTTP ${response.status} ${response.statusText}`);
    for (const [key, value] of response.headers) {
      console.log(`< ${key}: ${maskSensitiveHeader(key, value)}`);
    }
    console.log(`< `);
    try {
      const text = await response.clone().text();
      if (text) {
        console.log(text);
      }
    } catch {
      // ignore
    }
    console.log('');
  }
  return response;
}

export async function executeApiCall({
  method,
  path,
  queryParameters,
  body,
  timeout = 30000,
}: ApiCallArguments): Promise<Response> {
  await ensurePlatformAccessConfigurationIsAvailable();
  const url = await toUrl(path, queryParameters);
  const headers = await getHeaders();

  const response = await doWithRetry(method, async () => {
    await rateLimiter.acquire();
    // The deadline stays attached to the response, so it bounds reading the body as
    // well and not just the wait for the status line. Clearing it once the headers
    // arrive would leave a stalled body download running forever.
    const signal = AbortSignal.timeout(timeout);
    try {
      return await doFetch(url, method, headers, body ? JSON.stringify(body) : undefined, signal);
    } catch (e) {
      throw new Error(`Failed to call Steadybit API at ${method} ${url}: ${describeFetchError(e)}`, {
        cause: e,
      });
    }
  });

  if (!response.ok) {
    let body = '';
    try {
      body = await response.text();
    } catch {
      // ignore
    }
    throw new ApiError(
      `Steadybit API at ${method} ${url} responded with unexpected status code: ${response.status} - ${body || '<no body>'}`,
      response,
      body
    );
  }

  return response;
}

// Repeating a request that may already have been applied is only safe for methods
// defined to be idempotent: a POST that failed in transit might still have started an
// experiment run, so it is reported rather than retried.
const IDEMPOTENT_METHODS = new Set(['GET', 'HEAD', 'OPTIONS', 'PUT', 'DELETE']);

async function doWithRetry(method: string, fn: () => Promise<Response>): Promise<Response> {
  // A rate limit says the request was rejected, not applied, so waiting it out is safe
  // whatever the method. What it needs is time, not a number of tries: an attempt count
  // multiplied by the reset interval gave up after roughly eighteen seconds, which a
  // fan-out like `experiment dump` exceeds routinely because every rejected request
  // retries into the same window. A transport failure is still capped by attempts, since
  // repeating it is the part that carries risk.
  //
  // The budget counts time actually spent waiting on 429s, not elapsed time. A deadline
  // taken at the start would have been spent by the pacing in rateLimiter.acquire(),
  // which under a large dump can hold a request back for minutes — leaving no budget for
  // the rate limit the pacing exists to survive.
  const maxTransportAttempts = options.maxRetries + 2;
  const mayRepeat = IDEMPOTENT_METHODS.has(method.toUpperCase());
  let rateLimitWait = 0;
  let transportAttempt = 1;

  for (;;) {
    let response: Response;
    try {
      response = await fn();
    } catch (e) {
      // Transport failures used to end the command outright, which meant one flaky DNS
      // lookup out of the hundreds a dump makes discarded the whole run.
      if (!mayRepeat || transportAttempt >= maxTransportAttempts) {
        throw e;
      }
      await sleep(withJitter(options.defaultWaitTime * transportAttempt));
      transportAttempt++;
      continue;
    }

    if (response.status !== TOO_MANY_REQUESTS) {
      return response;
    }
    const resetHeader = response.headers.get('RateLimit-Reset') || response.headers.get('Retry-After');
    const retryInMillis = (resetHeader && Number.parseInt(resetHeader) * 1000) || options.defaultWaitTime;
    if (rateLimitWait + retryInMillis > options.rateLimitBudget) {
      return response;
    }
    await sleep(retryInMillis);
    rateLimitWait += retryInMillis;
  }
}

// Without this, requests that failed together — a dump shares one DNS resolver and one
// connection pool — would come back together and fail together again.
function withJitter(millis: number): number {
  return millis / 2 + Math.random() * (millis / 2);
}

const SENSITIVE_HEADERS = new Set(['authorization', 'cookie', 'set-cookie', 'proxy-authorization']);

function maskSensitiveHeader(name: string, value: string): string {
  return SENSITIVE_HEADERS.has(name.toLowerCase()) ? '<redacted>' : value;
}

// The global fetch reports every transport failure as "fetch failed" and carries the
// actual reason (DNS, TLS, ECONNREFUSED) on the cause chain.
function describeFetchError(e: unknown): string {
  const message = errorMessage(e);
  const causeMessage = (e as { cause?: Error })?.cause?.message;
  return causeMessage && causeMessage !== message ? `${message}: ${causeMessage}` : message;
}

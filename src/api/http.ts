// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { setTimeout as sleep } from 'node:timers/promises';
import { getHeaders, toUrl } from './common.ts';
import { ApiError } from './error.ts';
import { errorMessage } from '../errors.ts';
import { ensurePlatformAccessConfigurationIsAvailable } from '../config/requirePlatformAccess.ts';

const TOO_MANY_REQUESTS = 429;

export const options = {
  maxRetries: 2,
  defaultWaitTime: 1000,
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

  const response = await doWithRetry(async () => {
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

async function doWithRetry(fn: () => Promise<Response>): Promise<Response> {
  // maxRetries has always allowed one more retry than its name suggests. The count is
  // kept as-is to avoid shortening how long the CLI rides out a rate limit.
  const maxAttempts = options.maxRetries + 2;
  for (let attempt = 1; ; attempt++) {
    const response = await fn();
    if (response.status !== TOO_MANY_REQUESTS || attempt >= maxAttempts) {
      return response;
    }
    const resetHeader = response.headers.get('RateLimit-Reset') || response.headers.get('Retry-After');
    const retryInMillis = (resetHeader && Number.parseInt(resetHeader) * 1000) || options.defaultWaitTime;
    await sleep(retryInMillis);
  }
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

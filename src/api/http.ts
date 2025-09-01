// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import fetch, { Response } from 'node-fetch';
import { getHeaders, toUrl } from './common';

/*
 * This is required because we are supporting Node.js v14
 */
import { ensurePlatformAccessConfigurationIsAvailable } from '../config/requirePlatformAccess';
import http from 'http';
import https from 'https';

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
  expect2xx?: boolean; // defaults to true
  // https://github.com/node-fetch/node-fetch#manual-redirect
  redirect?: 'manual' | 'error';
  fullyQualifiedUrl?: boolean;
}

const httpAgentOptions = {
  keepAlive: true,
  keepAliveMsecs: 60000,
  timeout: 60000,
  maxSockets: 64,
};

const httpAgent = new http.Agent(httpAgentOptions);
const httpsAgent = new https.Agent(httpAgentOptions);

export async function executeApiCall({
  method,
  path,
  queryParameters,
  body,
  timeout = 30000,
  expect2xx = true,
  redirect = 'error',
  fullyQualifiedUrl = false,
}: ApiCallArguments): Promise<Response> {
  await ensurePlatformAccessConfigurationIsAvailable();
  const url = fullyQualifiedUrl ? path : await toUrl(path, queryParameters);
  const headers = await getHeaders();

  const controller = new AbortController();
  const response = await doWithRetry(async () => {
    const timeoutHandle = setTimeout(() => controller.abort(), timeout);
    try {
      return await fetch(url, {
        method,
        headers,
        body: body ? JSON.stringify(body) : undefined,
        signal: controller.signal,
        agent: getHttpAgent,
        redirect,
      });
    } catch (e) {
      throw new Error(
        `Failed to call Steadybit API at ${method} ${url}: ${(e as Error)?.message ?? 'Unknown Cause'}`,
        // @ts-expect-error TypeScript doesn't know about error causes yet
        { cause: e }
      );
    } finally {
      clearTimeout(timeoutHandle);
    }
  });

  if (expect2xx && !response.ok) {
    let body = '';
    try {
      body = await response.text();
    } catch {
      // ignore
    }
    const error: any = new Error(
      `Steadybit API at ${method} ${url} responded with unexpected status code: ${response.status} - ${body ?? '<no body>'}`
    );
    error.response = response;
    throw error;
  }

  return response;
}

async function doWithRetry(fn: () => Promise<fetch.Response>): Promise<fetch.Response> {
  let response;
  let retry = 0;
  do {
    response = await fn();
    if (response.status !== TOO_MANY_REQUESTS) {
      break;
    }
    const resetHeader = response.headers.get('RateLimit-Reset') || response.headers.get('Retry-After');
    const retryInSeconds = (resetHeader && Number.parseInt(resetHeader) * 1000) || options.defaultWaitTime;
    await sleep(retryInSeconds);
  } while (retry++ <= options.maxRetries);
  return response;
}

function sleep(millis: number): Promise<void> {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve();
    }, millis);
  });
}

function getHttpAgent(parsedUrl: URL) {
  if (parsedUrl.protocol == 'http:') {
    return httpAgent;
  } else {
    return httpsAgent;
  }
}

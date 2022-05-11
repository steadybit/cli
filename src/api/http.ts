// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

/*
 * This is required because we are supporting Node.js v14
 */
import { AbortController } from 'node-abort-controller';
import fetch, { Response } from 'node-fetch';
import https from 'https';
import http from 'http';

import { ensurePlatformAccessConfigurationIsAvailable } from '../config/requirePlatformAccess';
import { getHeaders, toUrl } from './common';

export interface ApiCallArguments {
  path: string;
  method: string;
  queryParameters?: Record<string, string>;
  body?: unknown;
  timeout?: number; // defaults to 30000
  expect2xx?: boolean; // defaults to true
  // https://github.com/node-fetch/node-fetch#manual-redirect
  redirect?: 'manual' | 'error';
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
}: ApiCallArguments): Promise<Response> {
  await ensurePlatformAccessConfigurationIsAvailable();
  const url = await toUrl(path, queryParameters);
  const headers = await getHeaders();

  const controller = new AbortController();
  const timeoutHandle = setTimeout(() => controller.abort(), timeout);

  let response: Response;
  try {
    response = await fetch(url, {
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

  if (expect2xx && !response.ok) {
    const error: any = new Error(
      `Steadybit API at ${method} ${url} responded with unexpected status code: ${response.status}`
    );
    error.response = response;
    throw error;
  }

  return response;
}

function getHttpAgent(parsedUrl: URL) {
  if (parsedUrl.protocol == 'http:') {
    return httpAgent;
  } else {
    return httpsAgent;
  }
}

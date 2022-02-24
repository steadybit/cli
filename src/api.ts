/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import fetch, { Response } from 'node-fetch';
import colors from 'colors/safe';
import https from 'https';
import http from 'http';

import { Configuration } from './config/types';
import { getConfiguration } from './config';
import { abortExecution } from './errors';

export interface ApiCallArguments {
  path: string;
  method: string;
  queryParameters?: Record<string, string>;
  body?: unknown;
  timeout?: number; // defaults to 30000
  expect2xx?: boolean; // defaults to true
  // https://github.com/node-fetch/node-fetch#manual-redirect
  redirect?: 'manual' | 'error'
}

// Prefer to load at runtime directly from the package.json to simplify
// the TypeScript build. Without this, we would have to make the build
// more complicated to adapt the root dir accordingly.
// eslint-disable-next-line
const packageJson = require('../package.json');

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
  redirect = 'error'
}: ApiCallArguments): Promise<Response> {
  const config = await getConfiguration();
  await checkPrerequisites(config);

  const controller = new AbortController();
  const timeoutHandle = setTimeout(() => controller.abort(), timeout);

  const url = `${config.baseUrl}${path}?${new URLSearchParams(queryParameters ?? {}).toString()}`;

  let response: Response;
  try {
    response = await fetch(url, {
      method,
      headers: {
        Authorization: `accessToken ${config.apiAccessToken}`,
        'Content-Type': 'application/json',
        Accept: 'application/json, */*',
        // eslint-disable-next-line @typescript-eslint/no-var-requires
        'User-Agent': `${packageJson.name}@${packageJson.version}`,
      },
      body: body ? JSON.stringify(body) : undefined,
      // @ts-expect-error Signal missing in the type definitions
      signal: controller.signal,
      agent: getHttpAgent,
      redirect
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
    const error: any = new Error(`Steadybit API at ${method} ${url} responded with unexpected status code: ${response.status}`);
    error.response = response;
    throw error;
  }

  return response;
}

async function checkPrerequisites(config: Configuration) {
  if (!config.apiAccessToken) {
    throw abortExecution(
      `API access token not defined. You can define the access token via profiles (${colors.bold(
        'steadybit config profile add'
      )} or the ${colors.bold('STEADYBIT_TOKEN')} environment variable.`
    );
  }
}

function getHttpAgent(parsedUrl: URL) {
  if (parsedUrl.protocol == 'http:') {
    return httpAgent;
  } else {
    return httpsAgent;
  }
}

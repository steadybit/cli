import fetch, { Response } from 'node-fetch';

import { abortExecution } from './errors';

const baseUrl = process.env.STEADYBIT_URL || 'https://platform.steadybit.io';
const accessToken = process.env.STEADYBIT_TOKEN;

export interface ApiCallArguments {
  path: string;
  method: string;
  queryParameters?: Record<string, string>;
  body?: never;
  timeout?: number; // defaults to 30000
  expect2xx?: boolean; // defaults to true
}

export async function executeApiCall({ method, path, queryParameters, body, timeout = 30000, expect2xx = true }: ApiCallArguments): Promise<Response> {
  await checkPrerequisites();

  const controller = new AbortController();
  const timeoutHandle = setTimeout(() => controller.abort(), timeout);

  const url = `${baseUrl}${path}?${new URLSearchParams(
    queryParameters ?? {}
  ).toString()}`;

  let response: Response;
  try {
    response = await fetch(url, {
      method: method,
      headers: {
        Authorization: `accessToken ${accessToken}`,
      },
      body: body ? JSON.stringify(body) : undefined,
      // @ts-expect-error Signal missing in the type definitions
      signal: controller.signal,
    });
  } catch (e) {
    throw new Error(
      `Failed to call Steadybit API at ${method} ${url}: ${
        (e as Error)?.message ?? 'Unknown Cause'
      }`,
      // @ts-expect-error TypeScript doesn't know about error causes yet
      { cause: e }
    );
  } finally {
    clearTimeout(timeoutHandle);
  }

  if (expect2xx && !response.ok) {
    throw new Error(
      `Steadybit API at ${method} ${url} responded with unexpected status code: ${response.status}`
    );
  }

  return response;
}

async function checkPrerequisites() {
  if (!accessToken) {
    throw abortExecution(
      'API access token not defined. You can define the access token via the "STEADYBIT_TOKEN" environment variable.'
    );
  }
}

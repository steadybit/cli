// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { getConfiguration } from '../config/index.ts';
import { abortExecution } from '../errors.ts';
import { packageJson } from '../packageJson.ts';

export async function toUrl(path: string, queryParameters?: Record<string, string>): Promise<string> {
  const config = await getConfiguration();
  let url = isAbsoluteUrl(path) ? onConfiguredOrigin(path, config.baseUrl) : `${config.baseUrl}${path}`;
  if (queryParameters) {
    url = `${url}${url.includes('?') ? '&' : '?'}${new URLSearchParams(queryParameters).toString()}`;
  }
  return url;
}

function isAbsoluteUrl(path: string): boolean {
  return /^https?:\/\//i.test(path);
}

// Absolute URLs reach us from platform responses, most notably the Location header a
// run returns, and every request carries the API access token. A platform behind a
// proxy legitimately names its public host there, which need not be the host the CLI
// was configured with, so only the origin is replaced and the path is kept verbatim.
// That keeps such deployments working while the token never leaves the configured host.
function onConfiguredOrigin(url: string, baseUrl: string): string {
  let target: URL;
  let base: URL;
  try {
    target = new URL(url);
    base = new URL(baseUrl);
  } catch {
    throw abortExecution("Cannot request '%s' relative to the configured platform at '%s'.", url, baseUrl);
  }

  return target.origin === base.origin ? url : `${base.origin}${target.pathname}${target.search}`;
}

export async function getHeaders(): Promise<Record<string, string>> {
  const config = await getConfiguration();
  return {
    Authorization: `accessToken ${config.apiAccessToken}`,
    'Content-Type': 'application/json',
    Accept: 'application/json, */*',
    'User-Agent': `${packageJson.name}@${packageJson.version}`,
  };
}

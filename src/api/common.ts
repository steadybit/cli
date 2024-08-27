// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { getConfiguration } from '../config';

// Prefer to load at runtime directly from the package.json to simplify
// the TypeScript build. Without this, we would have to make the build
// more complicated to adapt the root dir accordingly.
// eslint-disable-next-line
const packageJson = require('../../package.json');

export async function toUrl(path: string, queryParameters?: Record<string, string>): Promise<string> {
  const config = await getConfiguration();
  let url = `${config.baseUrl}${path}`;
  if (queryParameters) {
    url = `${url}?${new URLSearchParams(queryParameters).toString()}`;
  }
  return url;
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

/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import { getActiveProfile } from './profile/service';
import { Configuration } from './types';

export const defaultBaseUrl = 'https://platform.steadybit.io';

export async function getConfiguration(): Promise<Configuration> {
  let apiAccessToken: string | undefined;
  let baseUrl = defaultBaseUrl;

  const profile = await getActiveProfile();
  if (profile) {
    apiAccessToken = profile.apiAccessToken;
    baseUrl = profile.baseUrl ?? baseUrl;
  }

  // Environment arguments take precedence over the global system configuration.
  apiAccessToken = process.env.STEADYBIT_TOKEN ?? apiAccessToken;
  baseUrl = process.env.STEADYBIT_URL ?? baseUrl;

  // A typically error case is that the baseUrl carries a trailing slash. This is fine
  // in our persisted config files, but we don't want our CLI to internally work with those.
  if (baseUrl.endsWith('/')) {
    baseUrl = baseUrl.substring(0, baseUrl.length - 1);
  }

  return {
    apiAccessToken,
    baseUrl,
  };
}

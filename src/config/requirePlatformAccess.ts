// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import colors from 'colors/safe';

import { abortExecutionWithOpts } from '../errors';
import { getConfiguration } from './index';

const platformAccessConfigurationMissingHelp = `
No API access token configuration was found for Steadybit platform access.
You can configure API access tokens through configuration profiles or
environment variables (${colors.bold('STEADYBIT_TOKEN')}). We recommend configuration profiles
for local CLI usage. You can add a configuration profile via

                  ${colors.bold('steadybit config profile add')}
`.trim();

export async function ensurePlatformAccessConfigurationIsAvailable() {
  const config = await getConfiguration();

  if (!config.apiAccessToken) {
    throw abortExecutionWithOpts({colorize: false}, platformAccessConfigurationMissingHelp);
  }
}

// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH
import colors from '../../colors.ts';
import input from '@inquirer/input';
import password from '@inquirer/password';

import { cancelable } from '../../prompt/cancellation.ts';
import { validateNotBlank, validateHttpUrl } from '../../prompt/validation.ts';
import { addProfile } from './service.ts';
import { defaultBaseUrl } from '../index.ts';
import type { Profile } from './types.ts';

const startHelp = `
Configuration profiles enable you to use the CLI without repeatedly providing
passwords or having to remember environment variables. Configuration profiles
are stored in ~/.steadybit
`.trim();

const finishHelp = `
${colors.green('Done!')} You can now start using the CLI. For example, you could start
to run your first experiment via:

                   ${colors.bold('steadybit experiment run -k <your-key>')}
`.trim();

interface Options {
  name: string;
  baseUrl?: string;
  token: string;
}

export async function add(options: Options): Promise<void> {
  const profile: Profile =
    options?.name && options?.token
      ? { name: options.name, baseUrl: options.baseUrl, apiAccessToken: options.token }
      : await ask();
  await addProfile(profile);

  console.log();
  console.log(finishHelp);
}

async function ask(): Promise<Profile> {
  console.clear();
  console.log(startHelp);
  console.log();

  const name = await cancelable(
    input({
      message: 'Profile name:',
      validate: validateNotBlank,
    })
  );

  const baseUrl = await cancelable(
    input({
      message: 'Base URL of the Steadybit server:',
      default: defaultBaseUrl,
      validate: validateHttpUrl,
    })
  );

  console.log(`
The CLI will need an API access token of ${colors.bold('type team')} to communicate with
the Steadybit servers. You can generate one through the following URL:

          ${baseUrl.replace(/\/$/, '')}/settings/api-tokens
`);

  const apiAccessToken = await cancelable(
    password({
      message: 'API access token:',
      validate: validateNotBlank,
    })
  );

  return { name, baseUrl, apiAccessToken };
}

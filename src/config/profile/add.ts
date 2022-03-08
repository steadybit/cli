/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import colors from 'colors/safe';
import inquirer from 'inquirer';

import { validateNotBlank, validateHttpUrl } from '../../prompt/validation';
import { addProfile } from './service';
import { defaultBaseUrl } from '..';
import { Profile } from './types';

export async function add(): Promise<void> {
  const profile = await ask();
  await addProfile(profile);
  console.log(colors.green('Done!'));
}

async function ask(): Promise<Profile> {
  const answers1 = await inquirer.prompt([
    {
      type: 'input',
      name: 'name',
      message: 'Profile name:',
      validate: validateNotBlank,
    },
    {
      type: 'input',
      name: 'baseUrl',
      message: 'Base URL of the Steadybit server:',
      default: defaultBaseUrl,
      validate: validateHttpUrl,
    },
  ]);

  console.log(`
The CLI will need an API access token of ${colors.bold('type team')} to communicate with
the Steadybit servers. You can generate one through the following URL:

          ${answers1.baseUrl.replace(/\/$/, '')}/settings/api-tokens
`);

  const answers2 = await inquirer.prompt([
    {
      type: 'password',
      name: 'apiAccessToken',
      message: 'API access token:',
      validate: validateNotBlank,
    },
  ]);

  return {
    ...answers1,
    ...answers2
  };
}

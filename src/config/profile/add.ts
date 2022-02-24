/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import inquirer from 'inquirer';

import { validateNotBlank, validateHttpUrl } from '../../prompt/validation';
import { addProfile } from './service';
import { defaultBaseUrl } from '..';
import { Profile } from './types';

export async function add(): Promise<void> {
  const profile = await ask();
  await addProfile(profile);
}

async function ask(): Promise<Profile> {
  return await inquirer.prompt([
    {
      type: 'input',
      name: 'name',
      message: 'Profile name:',
      validate: validateNotBlank,
    },
    {
      type: 'password',
      name: 'apiAccessToken',
      message: 'API access token:',
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
}

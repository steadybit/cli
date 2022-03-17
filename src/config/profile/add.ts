/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import colors from 'colors/safe';
import inquirer from 'inquirer';

import { validateNotBlank, validateHttpUrl } from '../../prompt/validation';
import { addProfile } from './service';
import { defaultBaseUrl } from '..';
import { Profile } from './types';

const startHelp = `
Configuration profiles enable you to use the CLI without repeatedly providing
passwords or having to remember environment variables. Configuration profiles
are stored in ~/.steadybit
`.trim();

const finishHelp = `
${colors.green('Done!')} You can now start using the CLI. For example, you could start
to define your first service via:

                   ${colors.bold('steadybit service init')}
`.trim();

export async function add(): Promise<void> {
  console.clear();
  console.log(startHelp);
  console.log();

  const profile = await ask();
  await addProfile(profile);

  console.log();
  console.log(finishHelp);
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

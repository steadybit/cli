// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import inquirer from 'inquirer';

import { setActiveProfile, getProfiles } from './service';
import { abortExecution } from '../../errors';

export async function select(): Promise<void> {
  const activeProfileName = await promptProfileSelection('Choose the new active profile:');
  await setActiveProfile(activeProfileName);
}

export async function promptProfileSelection(message: string): Promise<string> {
  const profiles = await getProfiles();
  if (profiles.length === 0) {
    throw abortExecution('No profiles configured.');
  }

  const { activeProfileName } = await inquirer.prompt([
    {
      type: 'list',
      name: 'activeProfileName',
      message,
      choices: profiles.map(p => p.name),
    },
  ]);

  return activeProfileName;
}

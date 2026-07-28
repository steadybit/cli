// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import selectPrompt from '@inquirer/select';
import { cancelable } from '../../prompt/cancellation.ts';

import { setActiveProfile, getProfiles } from './service.ts';
import { abortExecution } from '../../errors.ts';

export async function select(): Promise<void> {
  const activeProfileName = await promptProfileSelection('Choose the new active profile:');
  await setActiveProfile(activeProfileName);
}

export async function promptProfileSelection(message: string): Promise<string> {
  const profiles = await getProfiles();
  if (profiles.length === 0) {
    throw abortExecution('No profiles configured.');
  }

  return await cancelable(
    selectPrompt({
      message,
      choices: profiles.map(p => p.name),
    })
  );
}

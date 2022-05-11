// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import colors from 'colors/safe';

import { getActiveProfile, getProfiles } from './service';

export async function list(): Promise<void> {
  const profiles = await getProfiles();
  const activeProfile = await getActiveProfile();

  profiles
    .slice()
    .sort((a, b) => a.name.localeCompare(b.name))
    .forEach(p => {
      const isActive = p.name === activeProfile?.name;

      if (isActive) {
        console.log('* %s', colors.green(p.name));
      } else {
        console.log('  %s', p.name);
      }
    });
}

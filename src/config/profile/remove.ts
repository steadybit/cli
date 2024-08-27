// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { promptProfileSelection } from './select';
import { removeProfile } from './service';

export async function remove(): Promise<void> {
  const profileName = await promptProfileSelection('Choose profile to delete:');
  await removeProfile(profileName);
}

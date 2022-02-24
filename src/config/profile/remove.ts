/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import { promptProfileSelection } from './select';
import { removeProfile } from './service';

export async function remove(): Promise<void> {
  const profileName = await promptProfileSelection('Choose profile to delete:');
  await removeProfile(profileName);
}


// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { homedir } from 'node:os';
import fs from 'node:fs/promises';
import path from 'node:path';

import { abortExecution, errorMessage } from '../../errors.ts';
import type { Profile } from './types.ts';

const configDir = path.join(homedir(), '.steadybit');
const profilesFile = path.join(configDir, 'profiles.json');
const activeProfileFile = path.join(configDir, 'activeProfile');

export async function addProfile(profile: Profile): Promise<void> {
  const profiles = await getProfiles();

  const updatedProfiles = profiles.filter(p => p.name !== profile.name).concat(profile);

  await writeProfiles(updatedProfiles);
}

export async function removeProfile(profileName: string): Promise<void> {
  const profiles = await getProfiles();

  const updatedProfiles = profiles.filter(p => p.name !== profileName);

  await writeProfiles(updatedProfiles);
}

export async function getProfiles(): Promise<Profile[]> {
  await ensureConfigDirectoryExists();

  let fileContent: string;
  try {
    fileContent = await fs.readFile(profilesFile, { encoding: 'utf8' });
  } catch (e) {
    if ((e as any)?.code === 'ENOENT') {
      return [];
    }

    throw abortExecution("Failed to read file '%s': %s", profilesFile, errorMessage(e));
  }

  try {
    return JSON.parse(fileContent);
  } catch (e) {
    throw abortExecution("Failed to parse file '%s' as JSON: %s", profilesFile, errorMessage(e));
  }
}

async function writeProfiles(profiles: Profile[]): Promise<void> {
  await ensureConfigDirectoryExists();

  try {
    await fs.writeFile(profilesFile, JSON.stringify(profiles, undefined, 2));
  } catch (e) {
    throw abortExecution("Failed to write to file '%s': %s", profilesFile, errorMessage(e));
  }
}

async function ensureConfigDirectoryExists() {
  await fs.mkdir(configDir, { recursive: true });
}

export async function getActiveProfile(): Promise<Profile | undefined> {
  await ensureConfigDirectoryExists();

  let activeProfileName: string | undefined;
  try {
    activeProfileName = await fs.readFile(activeProfileFile, { encoding: 'utf8' });
    // Users opening and saving the file might end up adding a trailing new line character.
    activeProfileName = activeProfileName.trim();
  } catch (e) {
    if ((e as any)?.code !== 'ENOENT') {
      throw abortExecution("Failed to read file '%s': %s", profilesFile, errorMessage(e));
    }
  }

  const profiles = await getProfiles();
  const activeProfile: Profile | undefined = profiles.find(p => p.name === activeProfileName?.trim()) ?? profiles[0];
  return activeProfile;
}

export async function setActiveProfile(profileName: string): Promise<void> {
  await ensureConfigDirectoryExists();

  try {
    await fs.writeFile(activeProfileFile, profileName);
  } catch (e) {
    throw abortExecution("Failed to write to file '%s': %s", activeProfileFile, errorMessage(e));
  }
}

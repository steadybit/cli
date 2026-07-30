// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { homedir } from 'node:os';
import fs from 'node:fs/promises';
import path from 'node:path';

import { abortExecution, errorMessage } from '../../errors.ts';
import type { Profile } from './types.ts';

// Resolved per call rather than at module load. HOME is what decides where the profile
// store lives, and a test — or anything else setting it after this module is imported —
// would otherwise be talking to the developer's real ~/.steadybit.
const configDir = () => path.join(homedir(), '.steadybit');
const profilesFile = () => path.join(configDir(), 'profiles.json');
const activeProfileFile = () => path.join(configDir(), 'activeProfile');

// The profile files are read on every API call, three times per call via
// getConfiguration(). Doing the work once turns thousands of redundant syscalls into a
// handful during commands like `experiment dump`.
//
// Remembered against the config directory rather than outright, because that directory
// follows HOME: a single cached value would answer for whichever home happened to be
// current first, and for the directory-creation entry that meant leaving a later
// directory unmade while every write into it failed. Failures are not remembered, so an
// unreadable or unwritable home stays retryable, and the writers below forget the entry
// for the home they wrote to.
function oncePerConfigDirectory<T>(work: () => Promise<T>): (() => Promise<T>) & { forget: () => void } {
  const done = new Map<string, Promise<T>>();
  const runOnce = () => {
    const directory = configDir();
    let result = done.get(directory);
    if (!result) {
      result = work().catch(e => {
        done.delete(directory);
        throw e;
      });
      done.set(directory, result);
    }
    return result;
  };
  runOnce.forget = () => {
    done.delete(configDir());
  };
  return runOnce;
}

const ensureConfigDirectoryExists = oncePerConfigDirectory(async () => {
  await fs.mkdir(configDir(), { recursive: true });
});

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

const readProfiles = oncePerConfigDirectory(async (): Promise<Profile[]> => {
  await ensureConfigDirectoryExists();

  let fileContent: string;
  try {
    fileContent = await fs.readFile(profilesFile(), { encoding: 'utf8' });
  } catch (e) {
    if ((e as any)?.code === 'ENOENT') {
      return [];
    }

    throw abortExecution("Failed to read file '%s': %s", profilesFile(), errorMessage(e));
  }

  try {
    return JSON.parse(fileContent);
  } catch (e) {
    throw abortExecution("Failed to parse file '%s' as JSON: %s", profilesFile(), errorMessage(e));
  }
});

export function getProfiles(): Promise<Profile[]> {
  return readProfiles();
}

async function writeProfiles(profiles: Profile[]): Promise<void> {
  await ensureConfigDirectoryExists();

  try {
    await fs.writeFile(profilesFile(), JSON.stringify(profiles, undefined, 2));
  } catch (e) {
    throw abortExecution("Failed to write to file '%s': %s", profilesFile(), errorMessage(e));
  }
  readProfiles.forget();
}

const readActiveProfileName = oncePerConfigDirectory(async (): Promise<string | undefined> => {
  await ensureConfigDirectoryExists();

  try {
    // Users opening and saving the file might end up adding a trailing new line character.
    return (await fs.readFile(activeProfileFile(), { encoding: 'utf8' })).trim();
  } catch (e) {
    if ((e as any)?.code !== 'ENOENT') {
      throw abortExecution("Failed to read file '%s': %s", activeProfileFile(), errorMessage(e));
    }
    return undefined;
  }
});

export async function getActiveProfile(): Promise<Profile | undefined> {
  const activeProfileName = await readActiveProfileName();
  const profiles = await getProfiles();
  const activeProfile: Profile | undefined = profiles.find(p => p.name === activeProfileName) ?? profiles[0];
  return activeProfile;
}

export async function setActiveProfile(profileName: string): Promise<void> {
  await ensureConfigDirectoryExists();

  try {
    await fs.writeFile(activeProfileFile(), profileName);
  } catch (e) {
    throw abortExecution("Failed to write to file '%s': %s", activeProfileFile(), errorMessage(e));
  }
  readActiveProfileName.forget();
}

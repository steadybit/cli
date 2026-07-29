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
// getConfiguration(). Reading them once per process turns thousands of redundant
// syscalls into a handful during commands like `experiment dump`. Failures are not
// cached, so they stay retryable, and the writers below drop the cache.
function readOnce<T>(read: () => Promise<T>): (() => Promise<T>) & { forget: () => void } {
  let cached: Promise<T> | undefined;
  const cachingRead = () => {
    cached ??= read().catch(e => {
      cached = undefined;
      throw e;
    });
    return cached;
  };
  cachingRead.forget = () => {
    cached = undefined;
  };
  return cachingRead;
}

// Keyed by directory rather than memoised outright: the path follows HOME, and caching
// a single "already created" flag would leave a later directory unmade while the writes
// into it fail. Failures are not cached, so an unwritable home stays retryable.
const createdDirectories = new Map<string, Promise<void>>();

function ensureConfigDirectoryExists(): Promise<void> {
  const directory = configDir();
  let created = createdDirectories.get(directory);
  if (!created) {
    created = fs.mkdir(directory, { recursive: true }).then(
      () => undefined,
      e => {
        createdDirectories.delete(directory);
        throw e;
      }
    );
    createdDirectories.set(directory, created);
  }
  return created;
}

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

const readProfiles = readOnce(async (): Promise<Profile[]> => {
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

const readActiveProfileName = readOnce(async (): Promise<string | undefined> => {
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

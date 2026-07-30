// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import fs from 'node:fs/promises';
import os from 'node:os';
import path from 'node:path';
import { beforeAll, describe, expect, it, vi } from 'vitest';

import { addProfile, getActiveProfile, getProfiles, setActiveProfile } from './service.ts';

// These tests write profile files, so HOME is redirected to a scratch directory. The
// guard refuses to run the suite if that ever stops working. A plain import is enough
// because service.ts resolves the config directory per call rather than at module load.
const fakeHome = await fs.mkdtemp(path.join(os.tmpdir(), 'steadybit-service-test-'));
process.env.HOME = fakeHome;
if (!os.homedir().startsWith(fakeHome)) {
  throw new Error(`refusing to run: home directory is ${os.homedir()}, not the scratch directory`);
}

const profilesFile = path.join(fakeHome, '.steadybit', 'profiles.json');

describe('profile service', () => {
  beforeAll(async () => {
    await addProfile({ name: 'first', baseUrl: 'https://one.example.com', apiAccessToken: 'a' });
    await addProfile({ name: 'second', baseUrl: 'https://two.example.com', apiAccessToken: 'b' });
  });

  it('should read the profiles file only once across repeated calls', async () => {
    const readSpy = vi.spyOn(fs, 'readFile');

    await getProfiles();
    await getProfiles();
    await getProfiles();

    const profileReads = readSpy.mock.calls.filter(([file]) => String(file) === profilesFile);
    expect(profileReads).toHaveLength(1);
    readSpy.mockRestore();
  });

  it('should resolve the active profile', async () => {
    await setActiveProfile('second');

    expect((await getActiveProfile())?.name).toBe('second');
  });

  it('should not serve a stale active profile after it changes', async () => {
    await setActiveProfile('second');
    expect((await getActiveProfile())?.name).toBe('second');

    await setActiveProfile('first');

    expect((await getActiveProfile())?.name).toBe('first');
  });

  it('should not serve stale profiles after one is added', async () => {
    expect((await getProfiles()).map(p => p.name)).not.toContain('third');

    await addProfile({ name: 'third', baseUrl: 'https://three.example.com', apiAccessToken: 'c' });

    expect((await getProfiles()).map(p => p.name)).toContain('third');
  });

  // The config directory used to be computed at module load, which meant nothing could
  // point the CLI at a different home once this module had been imported.
  it('should follow a home directory that changes after import', async () => {
    const otherHome = await fs.mkdtemp(path.join(os.tmpdir(), 'steadybit-other-home-'));
    const previous = process.env.HOME;
    process.env.HOME = otherHome;

    try {
      await setActiveProfile('written-elsewhere');
      expect(await fs.readFile(path.join(otherHome, '.steadybit', 'activeProfile'), 'utf8')).toEqual(
        'written-elsewhere'
      );
    } finally {
      process.env.HOME = previous;
    }
  });

  // The directory memo followed HOME while the read memo did not, so a changed home
  // produced correct directories with the previous home's contents in them.
  it('should read from a home directory that changes after import', async () => {
    const otherHome = await fs.mkdtemp(path.join(os.tmpdir(), 'steadybit-other-home-'));
    await fs.mkdir(path.join(otherHome, '.steadybit'), { recursive: true });
    await fs.writeFile(
      path.join(otherHome, '.steadybit', 'profiles.json'),
      JSON.stringify([{ name: 'only-over-here', apiAccessToken: 'z' }])
    );

    const previous = process.env.HOME;
    process.env.HOME = otherHome;
    try {
      expect((await getProfiles()).map(p => p.name)).toEqual(['only-over-here']);
    } finally {
      process.env.HOME = previous;
    }

    // ...and the original home is still answered correctly afterwards.
    expect((await getProfiles()).map(p => p.name)).toContain('first');
  });
});

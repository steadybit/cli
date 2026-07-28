// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import fs from 'node:fs/promises';
import os from 'node:os';
import path from 'node:path';
import { beforeAll, describe, expect, it, vi } from 'vitest';

// service.ts resolves the config directory from the home directory at module load, so
// HOME is redirected to a scratch directory before importing it. The guard below refuses
// to run the suite if that ever stops working, since these tests write profile files.
const fakeHome = await fs.mkdtemp(path.join(os.tmpdir(), 'steadybit-service-test-'));
process.env.HOME = fakeHome;
if (!os.homedir().startsWith(fakeHome)) {
  throw new Error(`refusing to run: home directory is ${os.homedir()}, not the scratch directory`);
}

const { addProfile, getActiveProfile, getProfiles, setActiveProfile } = await import('./service.ts');

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
});

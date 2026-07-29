// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import fs from 'node:fs/promises';
import os from 'node:os';
import path from 'node:path';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { wrapPrompt } from '@inquirer/testing/vitest';
import { answerPrompt, waitForPrompt } from '../../mocks/prompts.ts';
import { add } from './add.ts';
import type { Profile } from './types.ts';

vi.mock('@inquirer/input', async importOriginal => {
  const actual = await importOriginal<typeof import('@inquirer/input')>();
  return { ...actual, default: wrapPrompt(actual.default) };
});
vi.mock('@inquirer/password', async importOriginal => {
  const actual = await importOriginal<typeof import('@inquirer/password')>();
  return { ...actual, default: wrapPrompt(actual.default) };
});

// The flow writes a real profile store, so HOME points at a scratch directory.
const fakeHome = await fs.mkdtemp(path.join(os.tmpdir(), 'steadybit-add-test-'));
process.env.HOME = fakeHome;
if (!os.homedir().startsWith(fakeHome)) {
  throw new Error(`refusing to run: home directory is ${os.homedir()}, not the scratch directory`);
}

async function storedProfiles(): Promise<Profile[]> {
  return JSON.parse(await fs.readFile(path.join(fakeHome, '.steadybit', 'profiles.json'), 'utf8'));
}

describe('config profile add', () => {
  beforeEach(() => {
    // Otherwise the flow wipes the test runner's output.
    vi.spyOn(console, 'clear').mockImplementation(() => undefined);
    vi.spyOn(console, 'log').mockImplementation(() => undefined);
  });

  it('should ask for a name, a base url and a token, and store them', async () => {
    const done = add({} as never);

    await answerPrompt('Profile name:', 'from-the-prompt');
    await answerPrompt('Base URL of the Steadybit server:', 'https://platform.example.com');
    await answerPrompt('API access token:', 's3cr3t');
    await done;

    expect(await storedProfiles()).toContainEqual({
      name: 'from-the-prompt',
      baseUrl: 'https://platform.example.com',
      apiAccessToken: 's3cr3t',
    });
  });

  it('should fall back to the public platform when the base url is left empty', async () => {
    const done = add({} as never);

    await answerPrompt('Profile name:', 'defaulted');
    await answerPrompt('Base URL of the Steadybit server:', ''); // accept the offered default
    await answerPrompt('API access token:', 'tok');
    await done;

    expect((await storedProfiles()).find(p => p.name === 'defaulted')?.baseUrl).toEqual(
      'https://platform.steadybit.com'
    );
  });

  it('should refuse a blank name and ask again', async () => {
    const done = add({} as never);

    await answerPrompt('Profile name:', '   ');
    await waitForPrompt('You must provide a valid value');

    await answerPrompt('Profile name:', 'eventually-valid', { replace: true });
    await answerPrompt('Base URL of the Steadybit server:', '');
    await answerPrompt('API access token:', 'tok');
    await done;

    expect((await storedProfiles()).map(p => p.name)).toContain('eventually-valid');
  });

  it('should refuse a base url that is not http', async () => {
    const done = add({} as never);

    await answerPrompt('Profile name:', 'bad-url');
    await answerPrompt('Base URL of the Steadybit server:', 'ftp://files.example.com');
    await waitForPrompt('Unsupported protocol ftp:');

    await answerPrompt('Base URL of the Steadybit server:', 'https://platform.example.com', { replace: true });
    await answerPrompt('API access token:', 'tok');
    await done;

    expect((await storedProfiles()).map(p => p.name)).toContain('bad-url');
  });

  it('should skip the questions entirely when name and token are given', async () => {
    await add({ name: 'non-interactive', token: 'tok', baseUrl: 'https://given.example.com' });

    expect(await storedProfiles()).toContainEqual({
      name: 'non-interactive',
      baseUrl: 'https://given.example.com',
      apiAccessToken: 'tok',
    });
  });
});

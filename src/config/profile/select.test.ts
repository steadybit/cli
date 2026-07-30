// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import fs from 'node:fs/promises';
import os from 'node:os';
import path from 'node:path';
import { beforeAll, describe, expect, it, vi } from 'vitest';
import { wrapPrompt } from '@inquirer/testing/vitest';
import { answerPrompt, waitForPrompt } from '../../mocks/prompts.ts';
import { addProfile } from './service.ts';
import { select } from './select.ts';

vi.mock('@inquirer/select', async importOriginal => {
  const actual = await importOriginal<typeof import('@inquirer/select')>();
  return { ...actual, default: wrapPrompt(actual.default) };
});

const fakeHome = await fs.mkdtemp(path.join(os.tmpdir(), 'steadybit-select-test-'));
process.env.HOME = fakeHome;
if (!os.homedir().startsWith(fakeHome)) {
  throw new Error(`refusing to run: home directory is ${os.homedir()}, not the scratch directory`);
}

const activeProfileFile = path.join(fakeHome, '.steadybit', 'activeProfile');

describe('config profile select', () => {
  beforeAll(async () => {
    await addProfile({ name: 'alpha', apiAccessToken: 'a' });
    await addProfile({ name: 'beta', apiAccessToken: 'b' });
  });

  it('should offer every configured profile', async () => {
    const done = select();

    await waitForPrompt('Choose the new active profile:');
    await waitForPrompt('alpha');
    await waitForPrompt('beta');

    await answerPrompt('Choose the new active profile:', '');
    await done;
  });

  it('should make the chosen profile the active one', async () => {
    const done = select();

    await waitForPrompt('Choose the new active profile:');
    // Down to the second entry, then accept.
    await answerPrompt('Choose the new active profile:', '\x1b[B');
    await done;

    expect(await fs.readFile(activeProfileFile, 'utf8')).toEqual('beta');
  });
});

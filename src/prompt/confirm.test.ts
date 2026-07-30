// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { wrapPrompt } from '@inquirer/testing/vitest';
import { answerPrompt, waitForPrompt } from '../mocks/prompts.ts';
import { confirm } from './confirm.ts';

// confirm() reaches the prompt through a dynamic import, so the mock has to survive
// `await import(...)` rather than only a static one.
vi.mock('@inquirer/confirm', async importOriginal => {
  const actual = await importOriginal<typeof import('@inquirer/confirm')>();
  return { ...actual, default: wrapPrompt(actual.default) };
});

describe('confirm', () => {
  describe('without a terminal', () => {
    it('should answer with the non-interactive default instead of prompting', async () => {
      await expect(confirm('Run it?', { defaultWhenNonInteractive: false })).resolves.toBe(false);
      await expect(confirm('Run it?', { defaultWhenNonInteractive: true })).resolves.toBe(true);
    });
  });

  describe('with a terminal', () => {
    beforeEach(() => {
      // The guard in confirm() decides whether to prompt at all. Faking it is the only
      // way in-process; whether the guard reads the real terminal correctly is left to
      // the container tests, which get a genuine tty from `docker run -t`.
      Object.defineProperty(process.stdout, 'isTTY', { value: true, configurable: true });
    });

    afterEach(() => {
      Object.defineProperty(process.stdout, 'isTTY', { value: undefined, configurable: true });
    });

    it('should take the answer the user gives', async () => {
      const answer = confirm('Run it?', { defaultYes: false });

      await answerPrompt('Run it?', 'y');

      await expect(answer).resolves.toBe(true);
    });

    it('should offer the configured default', async () => {
      const answer = confirm('Run it?', { defaultYes: false });

      await waitForPrompt('(y/N)'); // the configured default is the one offered
      await answerPrompt('Run it?', '');

      await expect(answer).resolves.toBe(false);
    });
  });
});

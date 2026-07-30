// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { wrapPrompt } from '@inquirer/testing/vitest';
import { pressCtrlC, waitForPrompt } from '../mocks/prompts.ts';
import { confirm } from './confirm.ts';
import { cancelable } from './cancellation.ts';

vi.mock('@inquirer/confirm', async importOriginal => {
  const actual = await importOriginal<typeof import('@inquirer/confirm')>();
  return { ...actual, default: wrapPrompt(actual.default) };
});

describe('cancelable', () => {
  beforeEach(() => {
    Object.defineProperty(process.stdout, 'isTTY', { value: true, configurable: true });
  });

  afterEach(() => {
    Object.defineProperty(process.stdout, 'isTTY', { value: undefined, configurable: true });
    vi.restoreAllMocks();
  });

  // inquirer used to re-raise SIGINT; the @inquirer prompts reject instead, which would
  // otherwise surface as an unhandled rejection with a stack trace. Previously this was
  // only covered by throwing a hand-made ExitPromptError, never by a real cancellation.
  it('should exit quietly with the conventional SIGINT status when the user cancels', async () => {
    const exit = vi.spyOn(process, 'exit').mockImplementation((() => undefined) as never);

    // process.exit is stubbed, so cancelable() falls through to its rethrow instead of
    // ending the process. In production it never gets that far.
    void confirm('Run it?').catch(() => undefined);
    await waitForPrompt('Run it?');
    pressCtrlC();

    await vi.waitFor(() => expect(exit).toHaveBeenCalledWith(130));
  });

  it('should let every other failure through untouched', async () => {
    const boom = new Error('something else');

    await expect(cancelable(Promise.reject(boom))).rejects.toBe(boom);
  });

  it('should pass a normal answer straight through', async () => {
    await expect(cancelable(Promise.resolve('answered'))).resolves.toEqual('answered');
  });
});

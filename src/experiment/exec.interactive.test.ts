// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { wrapPrompt } from '@inquirer/testing/vitest';
import { answerPrompt } from '../mocks/prompts.ts';
import { givenAnotherExperimentIsRunning } from '../mocks/handlers.ts';
import { executeExperiments } from './exec.ts';

vi.mock('@inquirer/confirm', async importOriginal => {
  const actual = await importOriginal<typeof import('@inquirer/confirm')>();
  return { ...actual, default: wrapPrompt(actual.default) };
});

// These go through the prompts rather than past them. exec.test.ts covers the same
// commands with no terminal attached, which is what CI does and what the
// defaultWhenNonInteractive values are about; both branches matter.
describe('experiment run, with a terminal', () => {
  let logged: string[];

  beforeEach(() => {
    Object.defineProperty(process.stdout, 'isTTY', { value: true, configurable: true });
    logged = [];
    vi.spyOn(console, 'log').mockImplementation((...args) => {
      logged.push(args.join(' '));
    });
  });

  afterEach(() => {
    Object.defineProperty(process.stdout, 'isTTY', { value: undefined, configurable: true });
    vi.restoreAllMocks();
  });

  it('should not run anything when the confirmation is declined', async () => {
    // Throwing rather than returning, so that the flow stops here the way a real exit
    // would. A no-op stub lets execution carry on and run the experiment anyway.
    const exit = vi.spyOn(process, 'exit').mockImplementation(code => {
      throw new Error(`exited with ${code}`);
    });

    const done = executeExperiments({ key: 'TST-1', recursive: false });
    await answerPrompt('Are you sure you want to run the experiment?', 'n');

    await expect(done).rejects.toThrow('exited with 0');
    expect(exit).toHaveBeenCalledWith(0);
    expect(logged.join('\n')).not.toContain('Executing experiment');
  });

  // This recovery was unreachable until the response body it inspects stopped being
  // consumed before it got there, so it has never been exercised end to end.
  it('should offer to run in parallel when another experiment is already running', async () => {
    givenAnotherExperimentIsRunning();

    const done = executeExperiments({ key: 'TST-1', recursive: false });
    await answerPrompt('Are you sure you want to run the experiment?', 'y');
    await answerPrompt('There is already an experiment running', 'y');
    await done;

    expect(logged.join('\n')).toContain('Executing experiment: TST-1');
  });

  it('should give up when running in parallel is declined', async () => {
    givenAnotherExperimentIsRunning();

    const done = executeExperiments({ key: 'TST-1', recursive: false });
    await answerPrompt('Are you sure you want to run the experiment?', 'y');
    await answerPrompt('There is already an experiment running', 'n');

    await expect(done).rejects.toThrow('Failed to run experiment (TST-1)');
    expect(logged.join('\n')).not.toContain('Executing experiment');
  });
});

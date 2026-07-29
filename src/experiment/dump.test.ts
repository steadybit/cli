// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { afterEach, describe, expect, it, vi } from 'vitest';
import fs from 'node:fs/promises';
import path from 'node:path';
import { http, HttpResponse } from 'msw';
import { server } from '../mocks/server.ts';
import { givenExecutions } from '../mocks/handlers.ts';
import { getTempDir } from '../mocks/tempFiles.ts';
import type { Team } from '../team/types.ts';
import { dump, selectTeams } from './dump.ts';

const teams = [
  { key: 'TST', name: 'Test' },
  { key: 'ADM', name: 'Administrators' },
  { key: 'WEBHOOK', name: 'Webhook' },
] as Team[];

let directoryCount = 0;
function freshDirectory(): string {
  return path.join(getTempDir(), `dump-${directoryCount++}`);
}

describe('selectTeams', () => {
  it('should keep every team when none is named', () => {
    expect(selectTeams(teams, undefined)).toEqual(teams);
    expect(selectTeams(teams, [])).toEqual(teams);
  });

  it('should keep only the named teams', () => {
    expect(selectTeams(teams, ['WEBHOOK']).map(t => t.key)).toEqual(['WEBHOOK']);
    expect(selectTeams(teams, ['WEBHOOK', 'TST']).map(t => t.key)).toEqual(['TST', 'WEBHOOK']);
  });

  it('should match a key regardless of case', () => {
    expect(selectTeams(teams, ['webhook']).map(t => t.key)).toEqual(['WEBHOOK']);
  });

  // Skipping an unknown key would make a dump that covered less than was asked for look
  // exactly like one that covered everything.
  it('should refuse an unknown key rather than dumping less than asked', () => {
    expect(() => selectTeams(teams, ['NOPE'])).toThrow('No accessible team with key NOPE');
  });

  it('should name the available keys when one is unknown', () => {
    expect(() => selectTeams(teams, ['TST', 'NOPE'])).toThrow('Available: ADM, TST, WEBHOOK');
  });
});

describe('experiment dump', () => {
  afterEach(() => {
    process.exitCode = undefined;
  });

  it('should write the experiment and its executions', async () => {
    givenExecutions('TST-1', [1, 2, 3]);
    const directory = freshDirectory();

    await dump({ directory });

    expect((await fs.readdir(path.join(directory, 'TST-1'))).sort()).toEqual([
      'execution-1.yaml',
      'execution-2.yaml',
      'execution-3.yaml',
      'experiment.yaml',
    ]);
    expect(process.exitCode).toBeUndefined();
  });

  // The whole point of the change: one bad experiment must not discard the rest.
  it('should keep going when an experiment cannot be fetched, and report it', async () => {
    givenExecutions('TST-1', [1]);
    server.use(
      http.get('http://example.com/api/experiments/TST-1', () =>
        HttpResponse.json({ title: 'Server Error' }, { status: 500 })
      )
    );
    const problems = vi.spyOn(console, 'error').mockImplementation(() => undefined);
    const directory = freshDirectory();

    await dump({ directory });

    expect(problems.mock.calls.flat().join('\n')).toContain('TST-1');
    expect(problems.mock.calls.flat().join('\n')).toContain('1 experiments and 0 executions could not be dumped');
    problems.mockRestore();
  });

  it('should exit non-zero when an experiment could not be dumped', async () => {
    server.use(
      http.get('http://example.com/api/experiments/TST-1', () =>
        HttpResponse.json({ title: 'Server Error' }, { status: 500 })
      )
    );
    vi.spyOn(console, 'error').mockImplementation(() => undefined);

    await dump({ directory: freshDirectory() });

    expect(process.exitCode).toEqual(1);
    vi.restoreAllMocks();
  });

  // An execution that could not be fetched is as much a hole in the dump as an
  // experiment is, and used to be swallowed without a word or a non-zero status.
  it('should count and report executions it could not fetch', async () => {
    givenExecutions('TST-1', [1, 2, 3], [2, 3]);
    const problems = vi.spyOn(console, 'error').mockImplementation(() => undefined);
    const directory = freshDirectory();

    await dump({ directory });

    const written = await fs.readdir(path.join(directory, 'TST-1'));
    expect(written.sort()).toEqual(['execution-1.yaml', 'experiment.yaml']);

    const reported = problems.mock.calls.flat().join('\n');
    expect(reported).toContain('TST-1: 2 of 3 executions could not be fetched');
    expect(reported).toContain('0 experiments and 2 executions could not be dumped');
    expect(process.exitCode).toEqual(1);
    problems.mockRestore();
  });

  it('should dump only the requested team', async () => {
    const directory = freshDirectory();

    await dump({ directory, team: ['TST'] });

    expect(await fs.readdir(directory)).toEqual(['TST-1']);
  });
});

// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { afterEach, describe, expect, it, vi } from 'vitest';
import fs from 'node:fs/promises';
import path from 'node:path';
import { http, HttpResponse } from 'msw';
import { server } from '../mocks/server.ts';
import { givenExecutions } from '../mocks/handlers.ts';
import { getTempDir } from '../mocks/tempFiles.ts';
import { dump } from './dump.ts';

let directoryCount = 0;
function freshDirectory(): string {
  return path.join(getTempDir(), `dump-${directoryCount++}`);
}

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
});

// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import fs from 'node:fs/promises';
import path from 'node:path';
import { describe, expect, it, vi } from 'vitest';
import { HttpResponse } from 'msw';
import { respondTo } from '../mocks/recorder.ts';
import { getTempDir } from '../mocks/tempFiles.ts';
import { cancel } from './cancel.ts';
import { addProperty, setProperty } from './property.ts';
import { collectArtifacts, downloadArtifacts, listArtifacts, pathSegment } from './artifacts.ts';
import { getExecution } from './get.ts';
import type { Execution } from './api.ts';

const EXECUTION = {
  id: 42,
  key: 'ADM-1',
  state: 'COMPLETED',
  steps: [
    { stepType: 'wait', id: 'w' },
    {
      stepType: 'action',
      actionId: 'com.steadybit.extension_jmeter.run',
      targetExecutions: [
        { id: 'te-1', name: 'host-a', artifacts: ['report.zip', 'log.txt'] },
        { id: 'te-2', name: 'host-b', artifacts: ['report.zip'] },
      ],
    },
    {
      stepType: 'service-validation',
      customLabel: 'shop is healthy',
      validations: [
        { stepType: 'action', targetExecutions: [{ id: 'te-3', name: 'check', artifacts: ['result.json'] }] },
      ],
    },
  ],
} as unknown as Execution;

describe('execution', () => {
  describe('get', () => {
    it('prints the run as YAML by default', async () => {
      respondTo('get', '/api/experiments/executions/42', () => ({ json: { id: 42, state: 'RUNNING' } }));
      const logSpy = vi.spyOn(console, 'log');

      await getExecution({ id: 42 });

      expect(logSpy).toHaveBeenCalledWith('id: 42\nstate: RUNNING\n');
    });

    it('reports a run that does not exist', async () => {
      respondTo('get', '/api/experiments/executions/43', () => ({ status: 404 }));

      await expect(getExecution({ id: 43 })).rejects.toThrow('Experiment run 43 not found.');
    });
  });

  describe('cancel', () => {
    it('reports a cancel that was accepted', async () => {
      const requests = respondTo('post', '/api/experiments/executions/42/cancel', () => ({ status: 202 }));
      const logSpy = vi.spyOn(console, 'log');

      await cancel({ id: 42 });

      expect(requests).toHaveLength(1);
      expect(logSpy).toHaveBeenCalledWith('Experiment run %s is being canceled.', 42);
    });

    it('reports a run that had already ended', async () => {
      respondTo('post', '/api/experiments/executions/42/cancel', () => ({ status: 200 }));
      const logSpy = vi.spyOn(console, 'log');

      await cancel({ id: 42 });

      expect(logSpy).toHaveBeenCalledWith('Experiment run %s has already ended.', 42);
    });

    it('reports a run that does not exist', async () => {
      respondTo('post', '/api/experiments/executions/43/cancel', () => ({ status: 404 }));

      await expect(cancel({ id: 43 })).rejects.toThrow('Experiment run 43 not found.');
    });
  });

  describe('property', () => {
    it('sends a single value as a string, even when it looks like a number', async () => {
      const requests = respondTo('post', '/api/experiments/executions/42/properties/ticket/set', () => ({}));

      await setProperty({ id: 42, key: 'ticket', value: ['0042'] });

      expect(requests[0].body).toBe('0042');
    });

    it('sends several values as a list', async () => {
      const requests = respondTo('post', '/api/experiments/executions/42/properties/tickets/set', () => ({}));

      await setProperty({ id: 42, key: 'tickets', value: ['SHOP-1', 'SHOP-2'] });

      expect(requests[0].body).toEqual(['SHOP-1', 'SHOP-2']);
    });

    it('parses values as JSON on request', async () => {
      const requests = respondTo('post', '/api/experiments/executions/42/properties/score/set', () => ({}));

      await setProperty({ id: 42, key: 'score', value: ['7'], json: true });

      expect(requests[0].body).toBe(7);
    });

    it('rejects a value that is not JSON with --json', async () => {
      await expect(setProperty({ id: 42, key: 'score', value: ['seven'], json: true })).rejects.toThrow(
        "'seven' is not valid JSON"
      );
    });

    it('adds a value to a list property', async () => {
      const requests = respondTo('post', '/api/experiments/executions/42/properties/tickets/add', () => ({}));
      const logSpy = vi.spyOn(console, 'log');

      await addProperty({ id: 42, key: 'tickets', value: ['SHOP-3'] });

      expect(requests[0].body).toBe('SHOP-3');
      expect(logSpy).toHaveBeenCalledWith('Property %s of experiment run %s updated.', 'tickets', 42);
    });

    it('surfaces the platform validation error', async () => {
      respondTo('post', '/api/experiments/executions/42/properties/locked/set', () => ({
        status: 422,
        json: { title: 'Property locked is not editable in an execution' },
      }));

      await expect(setProperty({ id: 42, key: 'locked', value: ['x'] })).rejects.toThrow(
        'Property locked is not editable in an execution'
      );
    });
  });

  describe('artifacts', () => {
    it('collects artifacts of action steps and service validations', () => {
      expect(collectArtifacts(EXECUTION)).toEqual([
        {
          step: 'com.steadybit.extension_jmeter.run',
          target: 'host-a',
          targetExecutionId: 'te-1',
          artifactId: 'report.zip',
        },
        {
          step: 'com.steadybit.extension_jmeter.run',
          target: 'host-a',
          targetExecutionId: 'te-1',
          artifactId: 'log.txt',
        },
        {
          step: 'com.steadybit.extension_jmeter.run',
          target: 'host-b',
          targetExecutionId: 'te-2',
          artifactId: 'report.zip',
        },
        { step: 'shop is healthy', target: 'check', targetExecutionId: 'te-3', artifactId: 'result.json' },
      ]);
    });

    it('says so when a run has no artifacts', async () => {
      respondTo('get', '/api/experiments/executions/42', () => ({ json: { id: 42, steps: [] } }));
      const logSpy = vi.spyOn(console, 'log');

      await listArtifacts({ id: 42 });

      expect(logSpy).toHaveBeenCalledWith('Experiment run %s has no artifacts.', 42);
    });

    it('downloads every artifact into a directory per target execution', async () => {
      respondTo('get', '/api/experiments/executions/42', () => ({ json: EXECUTION as any }));
      respondTo(
        'get',
        '/api/experiments/executions/42/artifacts/:target/:artifact',
        ({ url }) => new HttpResponse(`content of ${url.pathname.split('/').slice(-2).join('/')}`)
      );
      const directory = path.join(getTempDir(), 'artifacts');

      await downloadArtifacts({ id: 42, directory });

      await expect(fs.readFile(path.join(directory, 'te-1', 'report.zip'), 'utf8')).resolves.toBe(
        'content of te-1/report.zip'
      );
      await expect(fs.readFile(path.join(directory, 'te-2', 'report.zip'), 'utf8')).resolves.toBe(
        'content of te-2/report.zip'
      );
      await expect(fs.readFile(path.join(directory, 'te-3', 'result.json'), 'utf8')).resolves.toBe(
        'content of te-3/result.json'
      );
    });

    it('writes a single artifact to --output', async () => {
      respondTo('get', '/api/experiments/executions/42', () => ({ json: EXECUTION as any }));
      respondTo('get', '/api/experiments/executions/42/artifacts/te-1/log.txt', () => new HttpResponse('the log'));
      const output = path.join(getTempDir(), 'my.log');

      await downloadArtifacts({ id: 42, directory: '.', artifact: 'log.txt', output });

      await expect(fs.readFile(output, 'utf8')).resolves.toBe('the log');
    });

    it('refuses --output when several artifacts match', async () => {
      respondTo('get', '/api/experiments/executions/42', () => ({ json: EXECUTION as any }));

      await expect(
        downloadArtifacts({ id: 42, directory: '.', artifact: 'report.zip', output: 'report.zip' })
      ).rejects.toThrow('2 artifacts match, but --output takes exactly one.');
    });

    it('keeps downloads inside the directory whatever the artifact is called', async () => {
      respondTo('get', '/api/experiments/executions/42', () => ({
        json: {
          id: 42,
          steps: [{ stepType: 'action', targetExecutions: [{ id: 'te-9', artifacts: ['../../evil'] }] }],
        },
      }));
      respondTo('get', '/api/experiments/executions/42/artifacts/:target/:artifact', () => new HttpResponse('x'));
      const directory = path.join(getTempDir(), 'contained');

      await downloadArtifacts({ id: 42, directory });

      await expect(fs.readFile(path.join(directory, 'te-9', 'evil'), 'utf8')).resolves.toBe('x');
    });

    it('never lets an id step out of a directory', () => {
      expect(['..', '.', '', '../..', 'a/../..'].map(pathSegment)).toEqual(['_', '_', '_', '_', '_']);
      expect(pathSegment('report.zip')).toBe('report.zip');
    });

    it('reports when nothing matches', async () => {
      respondTo('get', '/api/experiments/executions/42', () => ({ json: EXECUTION as any }));

      await expect(downloadArtifacts({ id: 42, directory: '.', artifact: 'nope' })).rejects.toThrow(
        'No matching artifacts found in experiment run 42.'
      );
    });
  });
});

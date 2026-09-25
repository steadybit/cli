// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { describe, expect, it, vi } from 'vitest';
import { applyExperiments } from './apply.ts';
import { executeExperiments } from './exec.ts';
import { respondTo } from '../mocks/recorder.ts';
import { writeFile } from '../mocks/tempFiles.ts';

const TEMPLATE = 'd7e65100-1d20-4980-be87-c351704910b8';

describe('experiment from template', () => {
  describe('apply', () => {
    it('creates an experiment with the placeholders from file and flags', async () => {
      const requests = respondTo('post', `/api/experiments/templates/${TEMPLATE}/experiment-create`, () => ({
        status: 201,
        headers: { location: 'http://example.com/api/experiments/ADM-12' },
      }));
      const placeholders = await writeFile('values.yml', { CLUSTER: 'dev', REPLICAS: 3 });
      const logSpy = vi.spyOn(console, 'log');

      await applyExperiments({
        recursive: false,
        template: TEMPLATE,
        team: 'ADM',
        environment: 'Global',
        externalId: 'shop-latency',
        placeholders,
        placeholder: { CLUSTER: 'prod' },
        variable: { endpoint: 'http://shop?x=1' },
        resetProperties: false,
      });

      expect(requests).toHaveLength(1);
      expect(requests[0].url.searchParams.get('resetProperties')).toBe('false');
      expect(requests[0].body).toEqual({
        team: 'ADM',
        environment: 'Global',
        externalId: 'shop-latency',
        placeholders: [
          { key: 'CLUSTER', value: 'prod' },
          { key: 'REPLICAS', value: 3 },
        ],
        experimentVariables: { endpoint: 'http://shop?x=1' },
      });
      expect(logSpy).toHaveBeenCalledWith('Experiment %s %s from template %s.', 'ADM-12', 'created', TEMPLATE);
    });

    it('accepts the platform list form of placeholders in a file', async () => {
      const requests = respondTo('post', `/api/experiments/templates/${TEMPLATE}/experiment-create`, () => ({
        status: 200,
        headers: { location: 'http://example.com/api/experiments/ADM-12' },
      }));
      const placeholders = await writeFile('values.json', [{ key: 'LIST', value: ['a', 'b'] }], 'json');

      await applyExperiments({
        recursive: false,
        template: TEMPLATE,
        team: 'ADM',
        placeholders,
        resetProperties: true,
      });

      expect((requests[0].body as any).placeholders).toEqual([{ key: 'LIST', value: ['a', 'b'] }]);
    });

    it('rejects a placeholders file that is neither a map nor a list of entries', async () => {
      const placeholders = await writeFile('values.yml', ['just', 'strings']);

      await expect(
        applyExperiments({ recursive: false, template: TEMPLATE, team: 'ADM', placeholders, resetProperties: true })
      ).rejects.toThrow('must be a map of key to value or a list of {key, value} entries');
    });

    it('requires a team to create an experiment', async () => {
      await expect(applyExperiments({ recursive: false, template: TEMPLATE, resetProperties: true })).rejects.toThrow(
        '--team is required to create an experiment from a template.'
      );
    });

    it('updates an existing experiment when a key is given', async () => {
      const requests = respondTo('post', `/api/experiments/templates/${TEMPLATE}/experiment-update/ADM-12`, () => ({
        status: 200,
      }));
      const logSpy = vi.spyOn(console, 'log');

      await applyExperiments({
        recursive: false,
        key: 'ADM-12',
        template: TEMPLATE,
        placeholder: { CLUSTER: 'prod' },
        resetProperties: true,
      });

      expect(requests[0].body).toEqual({ placeholders: [{ key: 'CLUSTER', value: 'prod' }] });
      expect(requests[0].url.searchParams.get('resetProperties')).toBe('true');
      expect(logSpy).toHaveBeenCalledWith('Experiment %s updated from template %s.', 'ADM-12', TEMPLATE);
    });

    it('refuses options an update would silently ignore', async () => {
      await expect(
        applyExperiments({ recursive: false, key: 'ADM-12', template: TEMPLATE, team: 'ADM', resetProperties: true })
      ).rejects.toThrow('only takes placeholders; remove --team');
    });

    it('reports a template that does not exist', async () => {
      respondTo('post', `/api/experiments/templates/${TEMPLATE}/experiment-create`, () => ({ status: 404 }));

      await expect(
        applyExperiments({ recursive: false, template: TEMPLATE, team: 'ADM', resetProperties: true })
      ).rejects.toThrow(`Experiment template ${TEMPLATE} not found.`);
    });

    it('still requires a file or a template', async () => {
      await expect(applyExperiments({ recursive: false })).rejects.toThrow(
        'Either --file or --template must be specified.'
      );
    });
  });

  describe('run', () => {
    const executed = (key: string, run: number) => ({
      status: 200,
      json: {
        key,
        executionId: run,
        apiLocation: `http://example.com/api/experiments/executions/${run}`,
        uiLocation: `http://example.com/experiments/edit/${key}/executions/${run}`,
      },
    });

    it('creates and runs an experiment from a template', async () => {
      const requests = respondTo('post', `/api/experiments/templates/${TEMPLATE}/experiment-execute`, () =>
        executed('ADM-12', 7)
      );
      const logSpy = vi.spyOn(console, 'log');

      await executeExperiments({
        recursive: false,
        yes: true,
        wait: false,
        template: TEMPLATE,
        team: 'ADM',
        placeholder: { CLUSTER: 'prod' },
        executionVariable: { region: 'eu' },
        resetProperties: true,
      });

      expect(requests[0].body).toEqual({
        team: 'ADM',
        placeholders: [{ key: 'CLUSTER', value: 'prod' }],
        executionVariables: { region: 'eu' },
      });
      expect(Object.fromEntries(requests[0].url.searchParams)).toEqual({
        resetProperties: 'true',
        allowParallel: 'false',
        forcePersist: 'true',
      });
      expect(logSpy).toHaveBeenCalledWith('Executing experiment:', 'ADM-12');
      expect(logSpy).toHaveBeenCalledWith('Experiment run API:', 'http://example.com/api/experiments/executions/7');
      expect(logSpy).toHaveBeenCalledWith(
        'Experiment run UI:',
        'http://example.com/experiments/edit/ADM-12/executions/7'
      );
    });

    it('retries validation errors without persisting the failed runs', async () => {
      let failures = 1;
      const requests = respondTo('post', `/api/experiments/templates/${TEMPLATE}/experiment-execute`, () =>
        failures-- > 0 ? { status: 422, json: { title: 'no targets' } } : executed('ADM-12', 8)
      );

      await executeExperiments({
        recursive: false,
        yes: true,
        wait: false,
        template: TEMPLATE,
        team: 'ADM',
        retries: 2,
        retryInterval: 0,
        resetProperties: true,
      });

      expect(requests).toHaveLength(2);
      expect(requests.every(r => r.url.searchParams.get('forcePersist') === 'false')).toBe(true);
    });

    it('runs in parallel with --yes when another experiment is running', async () => {
      const requests = respondTo('post', `/api/experiments/templates/${TEMPLATE}/experiment-execute`, ({ url }) =>
        url.searchParams.get('allowParallel') === 'true'
          ? executed('ADM-12', 9)
          : {
              status: 409,
              json: { type: 'https://steadybit.com/problems/another-experiment-running-exception' },
            }
      );

      await executeExperiments({
        recursive: false,
        yes: true,
        wait: false,
        template: TEMPLATE,
        team: 'ADM',
        resetProperties: true,
      });

      expect(requests.map(r => r.url.searchParams.get('allowParallel'))).toEqual(['false', 'true']);
    });

    it('refuses a key together with a template', async () => {
      await expect(
        executeExperiments({ recursive: false, yes: true, key: 'ADM-1', template: TEMPLATE, resetProperties: true })
      ).rejects.toThrow('--key cannot be combined with --template.');
    });
  });
});

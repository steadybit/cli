// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import fs from 'node:fs/promises';
import { describe, expect, it, vi } from 'vitest';
import { respondTo } from '../mocks/recorder.ts';
import { writeFile } from '../mocks/tempFiles.ts';
import { load } from '../yaml.ts';
import {
  applyServices,
  deleteService,
  getService,
  getServiceVariables,
  linkExperiment,
  listServiceExperiments,
  listServices,
  provideExperiment,
  setServiceVariables,
  showServiceRisk,
  unlinkExperiment,
} from './commands.ts';

const ID = '019cd80d-a4c9-775b-bdf8-2672a280ce7c';
const SERVICE = {
  id: ID,
  name: 'Checkout',
  environment: 'Global',
  team: 'ADM',
  logoId: 'service',
  logoColor: 'blue',
  query: 'k8s.namespace="shop"',
  validations: [],
  serviceProfile: 'Steadybit Starter',
  properties: {},
  variables: { region: 'eu' },
  version: 3,
  created: '2026-07-06T07:22:09Z',
  createdBy: { username: 'jane' },
  edited: '2026-07-06T07:22:09Z',
  editedBy: { username: 'jane' },
};

const printed = (spy: ReturnType<typeof vi.spyOn>) => spy.mock.calls.flat().join('\n');

describe('service', () => {
  it('lists services across pages with filters', async () => {
    const requests = respondTo('get', '/api/services', ({ url }) =>
      url.searchParams.get('page') === '0'
        ? { json: { items: [{ id: 'a', name: 'Checkout', team: 'ADM', environment: 'Global' }], nextPage: 1 } }
        : { json: { items: [{ id: 'b', name: 'Catalog', team: 'ADM', environment: 'Global' }] } }
    );
    const logSpy = vi.spyOn(console, 'log');

    await listServices({ team: ['ADM'], environment: ['Global'] });

    expect(requests[0].url.searchParams.getAll('teamKey')).toEqual(['ADM']);
    expect(requests[0].url.searchParams.getAll('environmentName')).toEqual(['Global']);
    expect(printed(logSpy)).toContain('Checkout');
    expect(printed(logSpy)).toContain('Catalog');
  });

  it('writes a service without the fields apply cannot send', async () => {
    respondTo('get', `/api/services/${ID}`, () => ({ json: SERVICE }));
    const file = await writeFile('service.yml', {});

    await getService({ id: ID, file });

    const written = load(await fs.readFile(file, 'utf8')) as Record<string, unknown>;
    expect(written).toMatchObject({ id: ID, name: 'Checkout', variables: { region: 'eu' } });
    for (const field of ['version', 'created', 'createdBy', 'edited', 'editedBy']) {
      expect(written).not.toHaveProperty(field);
    }
  });

  it('reports a service that does not exist', async () => {
    respondTo('get', '/api/services/nope', () => ({ status: 404 }));

    await expect(getService({ id: 'nope' })).rejects.toThrow('Service nope not found.');
  });

  it('creates a service from a file and writes the id back', async () => {
    const requests = respondTo('post', '/api/services', () => ({ status: 201, json: SERVICE }));
    const withoutId: Partial<typeof SERVICE> = { ...SERVICE };
    delete withoutId.id;
    delete withoutId.version;
    const file = await writeFile('new-service.yml', withoutId);
    const logSpy = vi.spyOn(console, 'log');

    await applyServices({ file: [file], recursive: false, deleteExperiments: true });

    expect(requests[0].url.searchParams.get('deleteExperiments')).toBe('true');
    expect(requests[0].body).not.toHaveProperty('createdBy');
    expect(requests[0].body).not.toHaveProperty('version');
    expect((load(await fs.readFile(file, 'utf8')) as Record<string, unknown>).id).toBe(ID);
    expect(logSpy).toHaveBeenCalledWith('Service %s (%s) %s.', 'Checkout', ID, 'created');
  });

  it('rejects a service file without a name', async () => {
    const file = await writeFile('nameless.yml', { team: 'ADM' });

    await expect(applyServices({ file: [file], recursive: false, deleteExperiments: false })).rejects.toThrow(
      'does not name the service'
    );
  });

  it('deletes a service', async () => {
    const requests = respondTo('delete', `/api/services/${ID}`, () => ({}));

    await deleteService({ id: ID });

    expect(requests).toHaveLength(1);
  });

  describe('risk', () => {
    const RISK = {
      risk: 42,
      categoryRisks: { Redundancy: { total: 40, experiment: 50, advice: 30 } },
      experimentRisks: [{ experimentKey: 'ADM-1', risk: 60 }],
      lastCalculated: '2026-09-25T09:03:46Z',
    };

    it('prints the overall, category and experiment risks', async () => {
      respondTo('get', `/api/services/${ID}/risk`, () => ({ json: RISK }));
      const logSpy = vi.spyOn(console, 'log');

      await showServiceRisk({ id: ID });

      expect(logSpy).toHaveBeenCalledWith('Risk of service %s: %s (calculated %s)', ID, 42, RISK.lastCalculated);
      expect(printed(logSpy)).toContain('Redundancy');
      expect(printed(logSpy)).toContain('ADM-1');
    });

    it('passes at or below the accepted risk and fails above it', async () => {
      respondTo('get', `/api/services/${ID}/risk`, () => ({ json: RISK }));

      await expect(showServiceRisk({ id: ID, failAbove: 42, type: 'json' })).resolves.toBeUndefined();
      await expect(showServiceRisk({ id: ID, failAbove: 41, type: 'json' })).rejects.toThrow(
        `Risk of service ${ID} is 42, above the accepted 41.`
      );
    });

    it('fails a gate when the risk has not been calculated', async () => {
      respondTo('get', `/api/services/${ID}/risk`, () => ({ json: {} }));

      await expect(showServiceRisk({ id: ID, failAbove: 100, type: 'json' })).rejects.toThrow('is unknown');
    });
  });

  describe('experiments', () => {
    it('lists provided and custom experiments, including ones not created yet', async () => {
      const requests = respondTo('get', `/api/services/${ID}/experiments`, () => ({
        json: {
          items: [
            { experimentKey: 'ADM-1', category: 'Redundancy', associationType: 'CUSTOM' },
            { templateId: 't-1', category: 'Scalability', associationType: 'PROVIDED' },
          ],
        },
      }));
      const logSpy = vi.spyOn(console, 'log');

      await listServiceExperiments({ id: ID, type: ['custom', 'provided'], category: ['Redundancy'] });

      expect(requests[0].url.searchParams.getAll('type')).toEqual(['CUSTOM', 'PROVIDED']);
      expect(requests[0].url.searchParams.getAll('category')).toEqual(['Redundancy']);
      expect(printed(logSpy)).toContain('(not created)');
    });

    it('says when nothing matches the filters rather than that there are none', async () => {
      respondTo('get', `/api/services/${ID}/experiments`, () => ({ json: { items: [] } }));
      const logSpy = vi.spyOn(console, 'log');

      await listServiceExperiments({ id: ID, type: ['custom'] });
      await listServiceExperiments({ id: ID });

      expect(logSpy).toHaveBeenCalledWith('Service %s has no matching experiments.', ID);
      expect(logSpy).toHaveBeenCalledWith('Service %s has no experiments.', ID);
    });

    it('links and unlinks a custom experiment', async () => {
      const links = respondTo('post', `/api/services/${ID}/experiments/custom`, () => ({ status: 201 }));
      const unlinks = respondTo('delete', `/api/services/${ID}/experiments/custom`, () => ({}));

      await linkExperiment({ id: ID, experiment: 'ADM-1', category: 'Redundancy' });
      await unlinkExperiment({ id: ID, experiment: 'ADM-1' });

      expect(links[0].body).toEqual({ experimentKey: 'ADM-1', category: 'Redundancy' });
      expect(unlinks[0].url.searchParams.get('experimentKey')).toBe('ADM-1');
    });

    it('creates a provided experiment from a template with placeholders', async () => {
      const requests = respondTo('post', `/api/services/${ID}/experiments/provided`, () => ({
        status: 201,
        headers: { location: 'http://example.com/api/experiments/ADM-9' },
      }));
      const logSpy = vi.spyOn(console, 'log');

      await provideExperiment({ id: ID, template: 't-1', placeholder: { REPLICAS: '3' }, resetProperties: false });

      expect(requests[0].body).toEqual({ templateId: 't-1', placeholders: [{ key: 'REPLICAS', value: '3' }] });
      expect(requests[0].url.searchParams.get('resetProperties')).toBe('false');
      expect(logSpy).toHaveBeenCalledWith(
        'Provided experiment %s of service %s %s from template %s.',
        'ADM-9',
        ID,
        'created',
        't-1'
      );
    });
  });

  describe('variables', () => {
    it('prints the variables', async () => {
      respondTo('get', `/api/services/${ID}/variables`, () => ({ json: { region: 'eu' } }));
      const logSpy = vi.spyOn(console, 'log');

      await getServiceVariables({ id: ID, type: 'json' });

      expect(logSpy).toHaveBeenCalledWith(JSON.stringify({ region: 'eu' }, undefined, 2));
    });

    it('merges variables from a file and arguments', async () => {
      const requests = respondTo('patch', `/api/services/${ID}/variables`, () => ({}));
      const file = await writeFile('variables.yml', { hosts: ['a', 'b'], region: 'us' });

      await setServiceVariables(['region=eu', 'url=http://x?a=b'], { id: ID, file, replace: false });

      expect(requests[0].body).toEqual({ hosts: ['a', 'b'], region: 'eu', url: 'http://x?a=b' });
    });

    it('replaces all variables with --replace, even with none given', async () => {
      const requests = respondTo('put', `/api/services/${ID}/variables`, () => ({}));

      await setServiceVariables([], { id: ID, replace: true });

      expect(requests[0].body).toEqual({});
    });

    it('rejects malformed arguments and empty merges', async () => {
      await expect(setServiceVariables(['novalue'], { id: ID, replace: false })).rejects.toThrow(
        "'novalue' is not in the form KEY=VALUE."
      );
      await expect(setServiceVariables([], { id: ID, replace: false })).rejects.toThrow('No variables given.');
    });
  });
});

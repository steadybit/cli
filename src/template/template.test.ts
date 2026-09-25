// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { describe, expect, it, vi } from 'vitest';
import { respondTo } from '../mocks/recorder.ts';
import { getTemplate, listTemplates } from './commands.ts';

const ID = 'd7e65100-1d20-4980-be87-c351704910b8';

describe('template', () => {
  it('lists templates matching the filters', async () => {
    const requests = respondTo('get', '/api/experiments/templates', () => ({
      json: { templates: [{ id: ID, templateTitle: 'Shop survives a database outage' }] },
    }));
    const logSpy = vi.spyOn(console, 'log');

    await listTemplates({ search: ['shop'], tag: ['k8s', 'db'] });

    expect(requests[0].url.searchParams.getAll('freeTextPhrases')).toEqual(['shop']);
    expect(requests[0].url.searchParams.getAll('tag')).toEqual(['k8s', 'db']);
    expect(logSpy.mock.calls.flat().join('\n')).toContain('Shop survives a database outage');
  });

  it('prints the placeholders of a template as a file to fill in', async () => {
    respondTo('get', `/api/experiments/templates/${ID}`, () => ({
      json: {
        id: ID,
        templateTitle: 't',
        placeholders: [
          { key: 'CLUSTER', name: 'Cluster', description: '' },
          { key: 'NAMESPACE', name: 'Namespace', description: '' },
        ],
      },
    }));
    const logSpy = vi.spyOn(console, 'log');

    await getTemplate({ id: ID, placeholders: true });

    expect(logSpy).toHaveBeenCalledWith("CLUSTER: ''\nNAMESPACE: ''\n");
  });

  it('reports a template that does not exist', async () => {
    respondTo('get', `/api/experiments/templates/${ID}`, () => ({ status: 404 }));

    await expect(getTemplate({ id: ID })).rejects.toThrow(`Experiment template ${ID} not found.`);
  });
});

// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import fs from 'node:fs/promises';
import { describe, expect, it, vi } from 'vitest';
import { respondTo } from '../mocks/recorder.ts';
import { writeFile } from '../mocks/tempFiles.ts';
import { load } from '../yaml.ts';
import { applyServiceProfiles, deleteServiceProfile, getServiceProfile, listServiceProfiles } from './commands.ts';

const ID = '019eacd7-fb2c-733a-bed5-99a935323db5';
const PROFILE = {
  id: ID,
  name: 'High Redundancy',
  origin: 'CUSTOM',
  templates: [{ category: 'Instance', templateIds: ['t-1', 't-2'] }],
  defaultProfile: false,
  version: 1,
  created: '2026-06-09T14:44:56Z',
  createdBy: 'jane',
  edited: '2026-06-09T14:49:13Z',
  editedBy: 'jane',
};

describe('service profile', () => {
  it('lists profiles with filters', async () => {
    const requests = respondTo('get', '/api/services/profiles', () => ({ json: { items: [PROFILE] } }));
    const logSpy = vi.spyOn(console, 'log');

    await listServiceProfiles({ name: 'Redund', origin: ['custom'], default: true });

    expect(Object.fromEntries(requests[0].url.searchParams)).toMatchObject({
      name: 'Redund',
      origin: 'CUSTOM',
      defaultProfile: 'true',
    });
    expect(logSpy.mock.calls.flat().join('\n')).toContain('High Redundancy');
  });

  it('writes a profile that can be applied again', async () => {
    respondTo('get', `/api/services/profiles/${ID}`, () => ({ json: PROFILE }));
    const file = await writeFile('profile.yml', {});

    await getServiceProfile({ id: ID, file });

    expect(load(await fs.readFile(file, 'utf8'))).toEqual({
      id: ID,
      name: 'High Redundancy',
      origin: 'CUSTOM',
      templates: PROFILE.templates,
    });
  });

  it('creates a custom profile when the file names no origin', async () => {
    const requests = respondTo('post', '/api/services/profiles', () => ({ status: 201, json: PROFILE }));
    const file = await writeFile('new-profile.yml', { name: 'High Redundancy', templates: PROFILE.templates });

    await applyServiceProfiles({ file: [file], recursive: false, deleteExperiments: false });

    expect(requests[0].body).toEqual({ name: 'High Redundancy', templates: PROFILE.templates, origin: 'CUSTOM' });
    expect(requests[0].url.searchParams.get('deleteExperiments')).toBe('false');
    expect((load(await fs.readFile(file, 'utf8')) as { id: string }).id).toBe(ID);
  });

  it('points at --delete-experiments when a change would remove provided experiments', async () => {
    respondTo('post', '/api/services/profiles', () => ({
      status: 422,
      json: {
        type: 'https://steadybit.com/problems/validation-exception',
        violations: [
          { field: 'templates', message: 'Cannot remove templates without setting `deleteExperiments` to true.' },
        ],
      },
    }));
    const file = await writeFile('shrunk-profile.yml', { ...PROFILE, templates: [] });

    await expect(applyServiceProfiles({ file: [file], recursive: false, deleteExperiments: false })).rejects.toThrow(
      'Service profile High Redundancy was not saved: the change would remove provided experiments. Pass --delete-experiments to delete them.'
    );
  });

  it('explains why a provided profile cannot be deleted', async () => {
    respondTo('delete', `/api/services/profiles/${ID}`, () => ({ status: 422 }));

    await expect(deleteServiceProfile({ id: ID })).rejects.toThrow('is provided by Steadybit and cannot be deleted.');
  });

  it('reports a profile that does not exist', async () => {
    respondTo('delete', '/api/services/profiles/nope', () => ({ status: 404 }));

    await expect(deleteServiceProfile({ id: 'nope' })).rejects.toThrow('Service profile nope not found.');
  });
});

// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import fs from 'node:fs/promises';
import { describe, expect, it, vi } from 'vitest';
import { respondTo } from '../mocks/recorder.ts';
import { writeFile } from '../mocks/tempFiles.ts';
import { load } from '../yaml.ts';
import {
  applySchedules,
  createSchedule,
  deleteSchedule,
  disableSchedule,
  enableSchedule,
  getSchedule,
  listSchedules,
  updateSchedule,
} from './commands.ts';

const ID = '01951394-727f-76a0-8675-c7519ebd0ff5';
const SCHEDULE = {
  id: ID,
  experimentKey: 'ADM-1',
  cron: '0 0 9 ? * MON-FRI',
  timezone: 'Europe/Berlin',
  enabled: true,
  allowParallel: false,
  editedBy: { username: 'jane' },
  lastUpdated: '2026-09-01T10:00:00Z',
};

describe('schedule', () => {
  describe('list', () => {
    it('filters by teams and experiments with repeated parameters', async () => {
      const requests = respondTo('get', '/api/experiments/schedules/v2', () => ({ json: [SCHEDULE] }));
      const logSpy = vi.spyOn(console, 'log');

      await listSchedules({ team: ['ADM', 'OPS'], experiment: ['ADM-1'] });

      expect(requests[0].url.searchParams.getAll('team')).toEqual(['ADM', 'OPS']);
      expect(requests[0].url.searchParams.getAll('experiment')).toEqual(['ADM-1']);
      const printed = logSpy.mock.calls.flat().join('\n');
      expect(printed).toContain(ID);
      expect(printed).toContain('0 0 9 ? * MON-FRI (Europe/Berlin)');
    });

    it('sends no filter when none is given', async () => {
      const requests = respondTo('get', '/api/experiments/schedules/v2', () => ({ json: [] }));
      const logSpy = vi.spyOn(console, 'log');

      await listSchedules({});

      expect(requests[0].url.search).toBe('');
      expect(logSpy).toHaveBeenCalledWith('No experiment schedules found.');
    });
  });

  describe('get', () => {
    it('leaves out the read-only fields so the file can be applied again', async () => {
      respondTo('get', `/api/experiments/schedules/${ID}`, () => ({ json: SCHEDULE }));
      const file = await writeFile('schedule.yml', {});

      await getSchedule({ id: ID, file });

      const expected: Partial<typeof SCHEDULE> = { ...SCHEDULE };
      delete expected.editedBy;
      delete expected.lastUpdated;
      expect(load(await fs.readFile(file, 'utf8'))).toEqual(expected);
    });

    it('reports a schedule that does not exist', async () => {
      respondTo('get', '/api/experiments/schedules/nope', () => ({ status: 404 }));

      await expect(getSchedule({ id: 'nope' })).rejects.toThrow('Experiment schedule nope not found.');
    });
  });

  describe('apply', () => {
    it('creates a schedule and writes its id back into the file', async () => {
      const requests = respondTo('post', '/api/experiments/schedules', () => ({ status: 201, json: SCHEDULE }));
      const file = await writeFile('new-schedule.yml', { experimentKey: 'ADM-1', cron: '0 0 9 ? * MON-FRI' });
      const logSpy = vi.spyOn(console, 'log');

      await applySchedules({ file: [file], recursive: false });

      expect(requests[0].body).toEqual({ experimentKey: 'ADM-1', cron: '0 0 9 ? * MON-FRI' });
      expect(load(await fs.readFile(file, 'utf8'))).toEqual({
        id: ID,
        experimentKey: 'ADM-1',
        cron: '0 0 9 ? * MON-FRI',
      });
      expect(logSpy).toHaveBeenCalledWith('Experiment schedule %s for %s %s.', ID, 'ADM-1', 'created');
    });

    it('updates a schedule that has an id without touching the file', async () => {
      const requests = respondTo('post', '/api/experiments/schedules', () => ({ status: 200, json: SCHEDULE }));
      const file = await writeFile('schedule.json', SCHEDULE, 'json');
      const before = await fs.readFile(file, 'utf8');

      await applySchedules({ file: [file], recursive: false });

      expect(requests[0].body).not.toHaveProperty('editedBy');
      expect(requests[0].body).toHaveProperty('id', ID);
      expect(await fs.readFile(file, 'utf8')).toBe(before);
    });

    it('rejects a file without an experiment key', async () => {
      const file = await writeFile('broken.yml', { cron: '* * * ? * *' });

      await expect(applySchedules({ file: [file], recursive: false })).rejects.toThrow(
        'does not name an experimentKey'
      );
    });
  });

  describe('create', () => {
    it('creates a schedule from flags', async () => {
      const requests = respondTo('post', '/api/experiments/schedules', () => ({ status: 201, json: SCHEDULE }));

      await createSchedule({
        experiment: 'ADM-1',
        cron: '0 0 9 ? * MON-FRI',
        timezone: 'Europe/Berlin',
        allowParallel: false,
        variable: { region: 'eu' },
      });

      expect(requests[0].body).toEqual({
        experimentKey: 'ADM-1',
        cron: '0 0 9 ? * MON-FRI',
        timezone: 'Europe/Berlin',
        allowParallel: false,
        enabled: true,
        variables: { region: 'eu' },
      });
    });

    it('creates a disabled one-off schedule', async () => {
      const requests = respondTo('post', '/api/experiments/schedules', () => ({ status: 201, json: SCHEDULE }));

      await createSchedule({ experiment: 'ADM-1', startAt: '2026-10-01T09:00:00Z', disabled: true });

      expect(requests[0].body).toEqual({ experimentKey: 'ADM-1', startAt: '2026-10-01T09:00:00Z', enabled: false });
    });

    it('needs either --cron or --start-at, not both', async () => {
      await expect(createSchedule({ experiment: 'ADM-1' })).rejects.toThrow(
        'Either --cron or --start-at must be specified.'
      );
      await expect(createSchedule({ experiment: 'ADM-1', cron: 'x', startAt: 'y' })).rejects.toThrow(
        '--cron and --start-at cannot be combined.'
      );
    });
  });

  describe('update', () => {
    it('patches only the given fields', async () => {
      const requests = respondTo('patch', `/api/experiments/schedules/${ID}`, () => ({ json: SCHEDULE }));

      await updateSchedule({ id: ID, cron: '0 30 8 ? * *' });

      expect(requests[0].body).toEqual({ cron: '0 30 8 ? * *' });
    });

    it('refuses an update without changes', async () => {
      await expect(updateSchedule({ id: ID })).rejects.toThrow('Nothing to update.');
    });

    it('enables and disables', async () => {
      const requests = respondTo('patch', `/api/experiments/schedules/${ID}`, () => ({ json: SCHEDULE }));
      const logSpy = vi.spyOn(console, 'log');

      await disableSchedule({ id: ID });
      await enableSchedule({ id: ID });

      expect(requests.map(r => r.body)).toEqual([{ enabled: false }, { enabled: true }]);
      expect(logSpy).toHaveBeenCalledWith('Experiment schedule %s for %s %s.', ID, 'ADM-1', 'disabled');
    });
  });

  describe('delete', () => {
    it('deletes a schedule', async () => {
      const requests = respondTo('delete', `/api/experiments/schedules/${ID}`, () => ({}));
      const logSpy = vi.spyOn(console, 'log');

      await deleteSchedule({ id: ID });

      expect(requests).toHaveLength(1);
      expect(logSpy).toHaveBeenCalledWith('Experiment schedule %s deleted.', ID);
    });

    it('reports a schedule that does not exist', async () => {
      respondTo('delete', '/api/experiments/schedules/nope', () => ({ status: 404 }));

      await expect(deleteSchedule({ id: 'nope' })).rejects.toThrow('Experiment schedule nope not found.');
    });
  });
});

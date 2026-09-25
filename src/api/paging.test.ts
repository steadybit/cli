// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { describe, expect, it } from 'vitest';
import { respondTo } from '../mocks/recorder.ts';
import { fetchAllPages } from './paging.ts';

describe('paging', () => {
  it('follows nextPage until the last page', async () => {
    const requests = respondTo('get', '/api/things', ({ url }) => {
      const page = Number(url.searchParams.get('page'));
      return { json: { items: [`item-${page}`], nextPage: page < 2 ? page + 1 : null } };
    });

    const items = await fetchAllPages<string>('/api/things', { team: ['A', 'B'] });

    expect(items).toEqual(['item-0', 'item-1', 'item-2']);
    expect(requests.map(r => r.url.searchParams.get('page'))).toEqual(['0', '1', '2']);
    expect(requests.every(r => r.url.searchParams.get('size') === '100')).toBe(true);
    expect(requests[0].url.searchParams.getAll('team')).toEqual(['A', 'B']);
  });

  it('stops when the platform omits nextPage', async () => {
    const requests = respondTo('get', '/api/things', () => ({ json: { items: ['only'] } }));

    await expect(fetchAllPages('/api/things')).resolves.toEqual(['only']);
    expect(requests).toHaveLength(1);
  });
});

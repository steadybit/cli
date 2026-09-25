// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import type { QueryParameters } from './common.ts';
import { executeApiCall } from './http.ts';

interface Page<T> {
  items?: T[];
  nextPage?: number | null;
}

// Paged endpoints return at most 100 items and point to the next page, so a listing
// that stopped at the first response would silently leave out the rest.
export async function fetchAllPages<T>(path: string, queryParameters: QueryParameters = {}): Promise<T[]> {
  const items: T[] = [];
  let page: number | null | undefined = 0;
  while (page !== undefined && page !== null) {
    const response = await executeApiCall({
      method: 'GET',
      path,
      queryParameters: { ...queryParameters, page: String(page), size: '100' },
    });
    const body = (await response.json()) as Page<T>;
    items.push(...(body.items ?? []));
    page = body.nextPage;
  }
  return items;
}

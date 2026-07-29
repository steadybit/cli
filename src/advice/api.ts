// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2024 Steadybit GmbH

import type { AdviceStatus, FetchAdviceRequest, FetchAdviceResponse } from './types.ts';

import { abortExecutionWithError } from '../errors.ts';
import { executeApiCall } from '../api/http.ts';
import { format } from 'node:util';

async function fetchAdvice(offset: number, environment: string, query?: string): Promise<FetchAdviceResponse> {
  const body: FetchAdviceRequest = {
    offset: offset,
    environmentName: environment,
  };
  if (query) {
    body.query = query;
  }
  try {
    const response = await executeApiCall({
      method: 'POST',
      path: '/api/advice',
      body,
    });
    return (await response.json()) as FetchAdviceResponse;
  } catch (e) {
    throw abortExecutionWithError(e, 'Failed to fetch advice status. HTTP request failed.');
  }
}

export async function fetchAllAdvice(environment: string, query?: string): Promise<AdviceStatus[]> {
  let offset = 0;
  const allAdvice: AdviceStatus[] = [];
  do {
    const response = await fetchAdvice(offset, environment, query);
    if (response.nextOffset) {
      offset = response.nextOffset;
    } else {
      offset = -1;
    }
    if (response.items.length > 0) {
      allAdvice.push(...response.items);
      console.log(format('Fetched %d of %d matching advice.', allAdvice.length, response.totalItems));
    } else {
      console.log('No matching advice.');
    }
  } while (offset > 0);
  return allAdvice;
}

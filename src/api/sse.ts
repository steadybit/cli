// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import EventSource from 'eventsource';
import { Observable, from, combineLatestWith, switchMap } from 'rxjs';
import { abortExecution } from '../errors';
import { getHeaders, toUrl } from './common';

export interface SseCallArguments {
  path: string;
  queryParameters?: Record<string, string>;
}

export function executeApiCall<T>({ path, queryParameters }: SseCallArguments): Observable<T> {
  return from(getHeaders())
    .pipe(combineLatestWith(from(toUrl(path, queryParameters))))
    .pipe(switchMap(([headers, url]) => executeApiCallInternal<T>(url, headers)));
}

function executeApiCallInternal<T>(url: string, headers: Record<string, string>): Observable<T> {
  return new Observable<T>(subscriber => {
    const es = new EventSource(url, {
      headers
    });

    es.onmessage = (event: MessageEvent) => {
      if (event.type !== 'heart-beat') {
        const data = JSON.parse(event.data);
        subscriber.next(data);
      }
    };

    es.onerror = (event: MessageEvent) => {
      // @ts-expect-error This prop is not part of the TypeScript definition, but it
      // is part of the documented EventSource API
      const status: any = event.status;
      if (status != null) {
        subscriber.error(new Error(`Failed to establish SSE connection to ${url}. Got status code ${status}`));
      } else {
        abortExecution('Connection to SSE endpoint %s failed: %s', url, event);
      }
    };

    return () => es.close();
  });
}

// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { http, HttpResponse, type JsonBodyType } from 'msw';
import { server } from './server.ts';

export interface RecordedRequest {
  method: string;
  url: URL;
  body: unknown;
}

type Method = 'get' | 'post' | 'put' | 'patch' | 'delete';

// Answers one endpoint and keeps what was sent to it, for tests that assert on the
// request a command makes rather than only on what it prints.
export function respondTo(
  method: Method,
  path: string,
  reply: (
    request: RecordedRequest
  ) => Response | { status?: number; json?: JsonBodyType; headers?: Record<string, string> }
): RecordedRequest[] {
  const requests: RecordedRequest[] = [];
  server.use(
    http[method](`http://example.com${path}`, async ({ request }) => {
      const text = await request.text();
      const recorded = { method: request.method, url: new URL(request.url), body: text ? JSON.parse(text) : undefined };
      requests.push(recorded);
      const response = reply(recorded);
      if (response instanceof Response) {
        return response;
      }
      return response.json === undefined
        ? new HttpResponse(null, { status: response.status ?? 200, headers: response.headers })
        : HttpResponse.json(response.json, { status: response.status ?? 200, headers: response.headers });
    })
  );
  return requests;
}

// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH
import { http, HttpResponse } from 'msw';
import { Experiment } from '../experiment/types';
import { FetchAdviceRequest, FetchAdviceResponse } from '../advice/types';
import * as url from 'node:url';

let retryCount = 0;
let runSequence = 1;
let experimentSequence = 1;
let experimentStore: Record<string, Experiment> = {};

export const resetExperiments = () => {
  retryCount = 0;
  experimentSequence = 1;
  runSequence = 1;
  experimentStore = { 'TST-1': EXPERIMENTS['TST-1'] };
};

export const EXPERIMENTS: Record<string, Experiment> = {
  'TST-1': {
    key: 'TST-1',
    name: 'Verify TTR fashion bestseller',
    team: 'TST',
    environment: 'Global',
    lanes: [
      {
        steps: [
          {
            type: 'action',
            ignoreFailure: false,
            parameters: {
              graceful: 'true',
            },
            actionType: 'container-stop-attack',
            radius: {
              targetType: 'container',
              predicate: {
                operator: 'AND',
                predicates: [
                  {
                    key: 'k8s.namespace',
                    operator: 'EQUALS',
                    values: ['steadybit-demo'],
                  },
                  {
                    key: 'k8s.deployment',
                    operator: 'EQUALS',
                    values: ['fashion-bestseller'],
                  },
                ],
              },
              query: null,
              percentage: 50,
            },
          },
        ],
      },
    ],
  },
  NEW: {
    name: 'Verify TTR fashion bestseller',
    team: 'TST',
    environment: 'Global',
    lanes: [
      {
        steps: [
          {
            type: 'action',
            ignoreFailure: false,
            parameters: {
              graceful: 'true',
            },
            actionType: 'container-stop-attack',
            radius: {
              targetType: 'container',
              predicate: {
                operator: 'AND',
                predicates: [
                  {
                    key: 'k8s.namespace',
                    operator: 'EQUALS',
                    values: ['steadybit-demo'],
                  },
                  {
                    key: 'k8s.deployment',
                    operator: 'EQUALS',
                    values: ['fashion-bestseller'],
                  },
                ],
              },
              query: null,
              percentage: 50,
            },
          },
        ],
      },
    ],
  },
};

const getTooManyRequestsHandler = http.get('http://example.com/api/status', async ({ request }) => {
  const headers: Record<string, string> = {};
  const query = url.parse(request.url, true).query;
  const reset = query.reset && String(query.reset);
  const times = Number(query.times);
  let code = Number(query.code) || 200;
  if (reset) {
    headers['RateLimit-Reset'] = reset;
  }
  if (times) {
    if (retryCount < times) {
      retryCount++;
    } else {
      code = 200;
    }
  }
  return HttpResponse.json('', { status: code, headers: headers });
});

const getExperimentHandler = http.get('http://example.com/api/experiments/:key', async ({ params }) => {
  const experiment = experimentStore[String(params.key)];
  if (experiment) {
    return HttpResponse.json(experiment);
  } else {
    return HttpResponse.json('', { status: 404 });
  }
});

const deleteExperimentHandler = http.delete('http://example.com/api/experiments/:key', async ({ params }) => {
  const experiment = experimentStore[String(params.key)];
  delete experimentStore[String(params.key)];
  return HttpResponse.json('', { status: experiment ? 200 : 404 });
});

const updateExperimentHandler = http.post('http://example.com/api/experiments/:key', async ({ request, params }) => {
  const experiment = experimentStore[String(params.key)];
  if (experiment) {
    experimentStore[String(params.key)] = request.json();
  }
  return HttpResponse.json('', { status: experiment ? 200 : 404 });
});

const upsertExperimentHandler = http.post('http://example.com/api/experiments', async ({ request }) => {
  const key = `NEW-${experimentSequence++}`;
  experimentStore[key] = request.json();
  return HttpResponse.json('', { status: 201, headers: { location: `http://example.com/api/experiments/${key}` } });
});

const executeExperimentHandler = http.post('http://example.com/api/experiments/:key/execute', ({ params }) => {
  const experiment = experimentStore[String(params.key)];
  const run = runSequence++;
  if (experiment) {
    return HttpResponse.json(
      {
        key: params.key,
        executionId: run,
        apiLocation: `http://example.com/api/experiments/executions/${run}`,
        uiLocation: `http://example.com/experiments/edit/${params.key}/executions/${run}?tenant=example&team=EXAMPLE`,
      },
      {
        status: 201,
        headers: { location: `http://example.com/api/experiments/executions/${run}` },
      }
    );
  } else {
    return HttpResponse.json('', { status: 404 });
  }
});

const executeUpsertExperimentHandler = http.post('http://example.com/api/experiments/execute', ({ request }) => {
  const key = `NEW-${experimentSequence++}`;
  const run = runSequence++;
  experimentStore[key] = request.json();
  return HttpResponse.json(
    {
      key: key,
      executionId: run,
      apiLocation: `http://example.com/api/experiments/executions/${run}`,
      uiLocation: `http://example.com/experiments/edit/${key}/executions/${run}?tenant=example&team=EXAMPLE`,
    },
    { status: 201, headers: { location: `http://example.com/api/experiments/executions/${run}` } }
  );
});

const fetchAdviceHandler = http.post('http://example.com/api/advice', async ({ request }) => {
  const body = (await request.json()) as FetchAdviceRequest;
  if (body.query === 'mock.response=ok') {
    const response: FetchAdviceResponse = {
      totalItems: 1,
      items: [
        {
          target: {
            reference: 'target-1-ref',
            label: 'target-1',
            type: 'host',
          },
          advice: {
            type: 'advice-type-1',
            label: 'advice-1',
            status: 'Implemented',
          },
          url: 'http://example.com/api/advice/1111',
        },
      ],
    };
    return HttpResponse.json(response);
  }
  if (body.offset === 0) {
    const response: FetchAdviceResponse = {
      nextOffset: 2,
      totalItems: 3,
      items: [
        {
          target: {
            reference: 'target-1-ref',
            label: 'target-1',
            type: 'host',
          },
          advice: {
            type: 'advice-type-1',
            label: 'advice-1',
            status: 'Validation needed',
          },
          url: 'http://example.com/api/advice/1111',
        },
        {
          target: {
            reference: 'target-2-ref',
            label: 'target-2',
            type: 'host',
          },
          advice: {
            type: 'advice-type-2',
            label: 'advice-2',
            status: 'Implemented',
          },
          url: 'http://example.com/api/advice/2222',
        },
      ],
    };
    return HttpResponse.json(response);
  } else {
    const response: FetchAdviceResponse = {
      totalItems: 3,
      items: [
        {
          target: {
            reference: 'target-3-ref',
            label: 'target-3',
            type: 'host',
          },
          advice: {
            type: 'advice-type-3',
            label: 'advice-3',
            status: 'Action needed',
          },
          url: 'http://example.com/api/advice/3333',
        },
      ],
    };
    return HttpResponse.json(response);
  }
});

export const handlers = [
  executeUpsertExperimentHandler,
  executeExperimentHandler,
  upsertExperimentHandler,
  updateExperimentHandler,
  deleteExperimentHandler,
  getExperimentHandler,
  fetchAdviceHandler,
  getTooManyRequestsHandler,
];

// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH
import { http, HttpResponse } from 'msw';
import type { Experiment } from '../experiment/types.ts';
import type { FetchAdviceRequest, FetchAdviceResponse } from '../advice/types.ts';

let retryCount = 0;
let runSequence = 1;
let experimentSequence = 1;
let experimentStore: Record<string, Experiment> = {};
let validationFailuresRemaining = 0;
let executionsPerExperiment: Record<string, number[]> = {};
let unfetchableExecutions = new Set<number>();
let anotherExperimentRunning = false;

// Lets a test reach the "an experiment is already running, run it in parallel?" recovery,
// which the platform only offers when the caller has not already asked for parallel.
export const givenAnotherExperimentIsRunning = () => {
  anotherExperimentRunning = true;
};

function executionsFor(key: string): { id: number }[] {
  return (executionsPerExperiment[key] ?? []).map(id => ({ id }));
}

// Lets a dump test set up an experiment whose executions cannot all be fetched.
export const givenExecutions = (key: string, ids: number[], unfetchable: number[] = []) => {
  executionsPerExperiment[key] = ids;
  unfetchable.forEach(id => unfetchableExecutions.add(id));
};

export const resetExperiments = () => {
  retryCount = 0;
  experimentSequence = 1;
  runSequence = 1;
  experimentStore = { 'TST-1': EXPERIMENTS['TST-1'] };
  validationFailuresRemaining = 0;
  executionsPerExperiment = {};
  unfetchableExecutions = new Set();
  anotherExperimentRunning = false;
};

export const setValidationFailures = (count: number) => {
  validationFailuresRemaining = count;
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
  const query = new URL(request.url).searchParams;
  const reset = query.get('reset');
  const times = Number(query.get('times'));
  let code = Number(query.get('code')) || 200;
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
  return HttpResponse.text(String(query.get('body')), { status: code, headers: headers });
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

const executeExperimentHandler = http.post('http://example.com/api/experiments/:key/execute', ({ params, request }) => {
  const experiment = experimentStore[String(params.key)];
  const requestUrl = new URL(request.url);
  const forcePersist = requestUrl.searchParams.get('forcePersist');

  if (validationFailuresRemaining > 0 && forcePersist === 'false') {
    validationFailuresRemaining--;
    return HttpResponse.json(
      {
        type: 'https://steadybit.com/problems/experiment-invalid-exception',
        title:
          'Had validation errors (lanes[0].steps[0].blastRadius.predicate: Please specify a query to select targets).',
        status: 422,
        instance: `/api/experiments/${params.key}/execute`,
      },
      { status: 422 }
    );
  }

  if (anotherExperimentRunning && requestUrl.searchParams.get('allowParallel') !== 'true') {
    return HttpResponse.json(
      {
        type: 'https://steadybit.com/problems/another-experiment-running-exception',
        title: 'Another experiment is currently running.',
        status: 409,
        instance: `/api/experiments/${params.key}/execute`,
      },
      { status: 409 }
    );
  }

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
  const requestUrl = new URL(request.url);
  const forcePersist = requestUrl.searchParams.get('forcePersist');

  if (validationFailuresRemaining > 0 && forcePersist === 'false') {
    validationFailuresRemaining--;
    return HttpResponse.json(
      {
        type: 'https://steadybit.com/problems/experiment-invalid-exception',
        title:
          'Had validation errors (lanes[0].steps[0].blastRadius.predicate: Please specify a query to select targets).',
        status: 422,
        instance: '/api/experiments/execute',
      },
      { status: 422 }
    );
  }

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
            status: 'IMPLEMENTED',
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
            status: 'VALIDATION_NEEDED',
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
            status: 'IMPLEMENTED',
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
            status: 'ACTION_NEEDED',
          },
          url: 'http://example.com/api/advice/3333',
        },
      ],
    };
    return HttpResponse.json(response);
  }
});

// The dump walk: teams, then each team's experiments, then each experiment's executions.
const getTeamsHandler = http.get('http://example.com/api/teams', () =>
  HttpResponse.json({ teams: [{ key: 'TST', name: 'Test Team' }] })
);

const listExperimentsHandler = http.get('http://example.com/api/experiments', ({ request }) => {
  const team = new URL(request.url).searchParams.get('team');
  const experiments = Object.values(experimentStore)
    .filter(experiment => experiment.team === team)
    .map(experiment => ({ key: experiment.key, name: experiment.name }));
  return HttpResponse.json({ experiments });
});

const listExecutionsHandler = http.get('http://example.com/api/experiments/:key/executions', ({ params }) =>
  experimentStore[String(params.key)]
    ? HttpResponse.json({ executions: executionsFor(String(params.key)) })
    : HttpResponse.json({ title: 'Not Found' }, { status: 404 })
);

const getExecutionHandler = http.get('http://example.com/api/experiments/executions/:id', ({ params }) =>
  unfetchableExecutions.has(Number(params.id))
    ? HttpResponse.json({ title: 'Server Error' }, { status: 500 })
    : HttpResponse.json({ id: Number(params.id), key: 'TST-1', state: 'COMPLETED' })
);

const getProblemHandler = http.get('http://example.com/api/problem', () =>
  HttpResponse.json({ type: 'https://steadybit.com/problems/another-experiment-running-exception' }, { status: 409 })
);

export const handlers = [
  getTeamsHandler,
  listExecutionsHandler,
  getExecutionHandler,
  listExperimentsHandler,
  getProblemHandler,
  executeUpsertExperimentHandler,
  executeExperimentHandler,
  upsertExperimentHandler,
  updateExperimentHandler,
  deleteExperimentHandler,
  getExperimentHandler,
  fetchAdviceHandler,
  getTooManyRequestsHandler,
];

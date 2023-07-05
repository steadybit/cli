// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH
import { rest } from 'msw';
import { Experiment } from '../experiment/types';

let runSequence = 1;
let experimentSequence = 1;
let experimentStore: Record<string, Experiment> = {};

export const resetExperiments = () => {
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
              graceful: 'true'
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
                    values: ['steadybit-demo']
                  },
                  {
                    key: 'k8s.deployment',
                    operator: 'EQUALS',
                    values: ['fashion-bestseller']
                  }
                ]
              },
              query: null,
              percentage: 50
            }
          }
        ]
      }
    ]
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
              graceful: 'true'
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
                    values: ['steadybit-demo']
                  },
                  {
                    key: 'k8s.deployment',
                    operator: 'EQUALS',
                    values: ['fashion-bestseller']
                  }
                ]
              },
              query: null,
              percentage: 50
            }
          }
        ]
      }
    ]
  }
};

const getExperimentHandler = rest.get('http://test/api/experiments/:key', async (req, res, ctx) => {
  const experiment = experimentStore[String(req.params.key)];
  if (experiment) {
    return res(ctx.json(experiment));
  } else {
    return res(ctx.status(404), ctx.body(''));
  }
});

const deleteExperimentHandler = rest.delete('http://test/api/experiments/:key', async (req, res, ctx) => {
  const experiment = experimentStore[String(req.params.key)];
  delete experimentStore[String(req.params.key)];
  return res(
    ctx.status(experiment ? 200 : 404),
    ctx.body('')
  );
});

const updateExperimentHandler = rest.post('http://test/api/experiments/:key', async (req, res, ctx) => {
  const experiment = experimentStore[String(req.params.key)];
  if (experiment) {
    experimentStore[String(req.params.key)] = req.json();
  }
  return res(
    ctx.status(experiment ? 200 : 404),
    ctx.body('')
  );
});

const upsertExperimentHandler = rest.post('http://test/api/experiments', async (req, res, ctx) => {
  const key = `NEW-${experimentSequence++}`;
  experimentStore[key] = req.json();
  return res(
    ctx.status(201),
    ctx.body(''),
    ctx.set({ location: `http://test/api/experiments/${key}` })
  );
});

const executeExperimentHandler = rest.post('http://test/api/experiments/:key/execute', async (req, res, ctx) => {
  const experiment = experimentStore[String(req.params.key)];
  if (experiment) {
    return res(
      ctx.status(201),
      ctx.body(''),
      ctx.set({ location: `http://test/api/experiments/executions/${runSequence++}` })
    );
  } else {
    return res(
      ctx.status(404),
      ctx.body('')
    );
  }
});

const executeUpsertExperimentHandler = rest.post('http://test/api/experiments/execute', async (req, res, ctx) => {
  const key = `NEW-${experimentSequence++}`;
  experimentStore[key] = req.json();
  return res(
    ctx.status(201),
    ctx.json({ key }),
    ctx.set({ location: `http://test/api/experiments/executions/${runSequence++}` })
  );
});

export const handlers = [
  executeUpsertExperimentHandler,
  executeExperimentHandler,
  upsertExperimentHandler,
  updateExperimentHandler,
  deleteExperimentHandler,
  getExperimentHandler
];

// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { setTimeout as sleep } from 'node:timers/promises';
import { confirm } from '../prompt/confirm.ts';
import { loadExperiment, resolveExperimentFiles, writeFile } from './files.ts';
import * as api from './api.ts';
import { ApiError } from '../api/error.ts';
import { abortExecution, abortExecutionWithError } from '../errors.ts';
import type { ExecuteResult } from './types.ts';

interface Options {
  key?: string;
  file?: string[];
  yes?: boolean;
  wait?: boolean;
  recursive: boolean;
  allowParallel?: boolean;
  retries?: number;
  retryInterval?: number;
}

export async function executeExperiments(options: Options) {
  if (!options.yes) {
    const confirmation = await confirm('Are you sure you want to run the experiment?', {
      defaultYes: false,
      defaultWhenNonInteractive: true,
    });
    if (!confirmation) {
      process.exit(0);
    }
  }

  if (!options.file && !options.key) {
    throw abortExecution('Either --key or --file must be specified.');
  } else if (options.file) {
    const files = await resolveExperimentFiles(options.file, options.recursive);
    if (files.length > 1 && options.key) {
      throw abortExecution('If --key is specified, at most one --file can be specified.');
    }

    for (const file of files) {
      const { experiment, datatype } = await loadExperiment(file);
      let key = options.key || experiment.key;
      let result: ExecuteResult;

      const hasRetries = (options.retries ?? 0) > 0;
      if (key) {
        await api.updateExperiment(key, experiment);
        result = await executeWithRetry(
          () => api.executeExperiment(key!, !!options.yes, options.allowParallel, !hasRetries),
          options.retries,
          options.retryInterval
        );
      } else {
        const upsertResult = await executeWithRetry(
          () => api.upsertAndExecuteExperiment(experiment, options.allowParallel, !hasRetries),
          options.retries,
          options.retryInterval
        );
        key = upsertResult.key;
        result = upsertResult;
        if (!experiment.key) {
          await writeFile(file, { key, ...experiment }, datatype);
        }
      }

      console.log('Executing experiment:', key);
      console.log('Experiment run API:', result.location);
      console.log('Experiment run UI:', result.uiLocation);
      /* eslint-disable @typescript-eslint/no-unused-expressions */
      options.wait && result.location && (await waitFor(result.location));
    }
  } else if (options.key) {
    const hasRetries = (options.retries ?? 0) > 0;
    const result = await executeWithRetry(
      () => api.executeExperiment(options.key!, !!options.yes, options.allowParallel, !hasRetries),
      options.retries,
      options.retryInterval
    );
    console.log('Experiment run API:', result.location);
    console.log('Experiment run UI:', result.uiLocation);
    console.log('Executing experiment:', options.key);
    /* eslint-disable @typescript-eslint/no-unused-expressions */
    options.wait && result.location && (await waitFor(result.location));
  }
}

async function executeWithRetry<T>(fn: () => Promise<T>, retries = 0, retryInterval = 10): Promise<T> {
  for (let attempt = 0; attempt <= retries; attempt++) {
    try {
      return await fn();
    } catch (e) {
      if (e instanceof ApiError && e.status === 422 && attempt < retries) {
        console.log(
          `Experiment has validation errors (attempt ${attempt + 1}/${retries + 1}). Retrying in ${retryInterval}s...`
        );
        await sleep(retryInterval * 1000);
        continue;
      }
      throw abortExecutionWithError(e, 'Failed to execute experiment');
    }
  }
  throw new Error('Unexpected end of retry loop');
}

async function waitFor(location: string): Promise<void> {
  // Loaded lazily because polling a run is the only thing in the CLI that needs rxjs,
  // and importing it costs far more than the rest of this subcommand's module graph.
  const { filter, firstValueFrom, from, interval, switchMap, tap } = await import('rxjs');

  const executionResult = await firstValueFrom(
    interval(5000)
      .pipe(switchMap(() => from(api.getExperimentExecutionUsingUrl(location ?? ''))))
      .pipe(tap(e => console.log('Current run state:', e.state.toLowerCase())))
      .pipe(
        filter(e => e.state === 'FAILED' || e.state === 'ERRORED' || e.state === 'CANCELED' || e.state === 'COMPLETED')
      )
  );

  if (executionResult && executionResult.state !== 'COMPLETED') {
    console.error(
      `Experiment ${executionResult.key} (#${executionResult.id}) ${executionResult.state.toLowerCase()}${executionResult.reason ? `, reason: ${executionResult.reason}` : ''}`
    );
    process.exit(1);
  }
}

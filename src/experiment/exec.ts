// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { confirm } from '../prompt/confirm';
import { loadExperiment, resolveExperimentFiles, writeYamlFile } from './files';
import * as api from './api';
import { filter, firstValueFrom, from, interval, switchMap, tap } from 'rxjs';
import { abortExecution } from '../errors';
import { ExecuteResult } from './types';

interface Options {
  key?: string;
  file?: string[];
  yes?: boolean;
  wait?: boolean;
  recursive: boolean;
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
      const experiment = await loadExperiment(file);
      let key = options.key || experiment.key;
      let result: ExecuteResult;

      if (key) {
        await api.updateExperiment(key, experiment);
        result = await api.executeExperiment(key, !!options.yes);
      } else {
        const upsertResult = await api.upsertAndExecuteExperiment(experiment);
        key = upsertResult.key;
        result = upsertResult;
        if (!experiment.key) {
          await writeYamlFile(file, { key, ...experiment });
        }
      }

      console.log('Executing experiment:', key);
      console.log('Experiment run API:', result.location);
      console.log('Experiment run UI:', result.uiLocation);
      /* eslint-disable @typescript-eslint/no-unused-expressions */
      options.wait && result.location && (await waitFor(result.location));
    }
  } else if (options.key) {
    const result = await api.executeExperiment(options.key, !!options.yes);
    console.log('Experiment run API:', result.location);
    console.log('Experiment run UI:', result.uiLocation);
    console.log('Executing experiment:', options.key);
    /* eslint-disable @typescript-eslint/no-unused-expressions */
    options.wait && result.location && (await waitFor(result.location));
  }
}

async function waitFor(location: string): Promise<void> {
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

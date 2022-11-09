// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { confirm } from '../prompt/confirm';
import { loadExperiment, writeExperiment } from './files';
import * as api from './api';
import { filter, firstValueFrom, from, interval, map, switchMap, tap } from 'rxjs';
import { abortExecution } from '../errors';

interface Options {
  key?: string;
  file?: string;
  yes?: boolean;
  wait?: boolean;
}

export async function executeExperiment(options: Options) {
  if (!options.yes) {
    const confirmation = await confirm('Are you sure you want to run the experiment?', {
      defaultYes: false,
      defaultWhenNonInteractive: true,
    });
    if (!confirmation) {
      process.exit(0);
    }
  }

  let location: string | undefined;

  if (!options.file) {
    if (!options.key) {
      throw abortExecution('Either --key or --file must be specified.');
    }
    const result = await api.executeExperiment(options.key);
    console.log('Executing experiment:', options.key);
    location = result.location;
  } else {
    const experiment = await loadExperiment(options.file);
    const key = options.key || experiment.key;

    if (key) {
      await api.updateExperiment(key, experiment);
      const result = await api.executeExperiment(key);
      console.log('Executing experiment:', key);
      location = result.location;
    } else {
      const result = await api.upsertAndExecuteExperiment(experiment);
      if (!experiment.key) {
        await writeExperiment(options.file, { key: result.key, ...experiment });
      }
      console.log('Executing experiment:', result.key);
      location = result.location;
    }
  }

  console.log('Experiment run:', location);

  if (!options.wait || !location) {
    return;
  }

  const result = await firstValueFrom(
    interval(5000)
      .pipe(switchMap(() => from(api.getExperimentExecution(location ?? ''))))
      .pipe(tap(e => console.log('Current run state:', e.state.toLowerCase())))
      .pipe(
        filter(e => e.state === 'FAILED' || e.state === 'ERRORED' || e.state === 'CANCELED' || e.state === 'COMPLETED')
      )
      .pipe(map(e => e.state === 'COMPLETED'))
  );

  process.exit(!result ? 1 : 0);
}

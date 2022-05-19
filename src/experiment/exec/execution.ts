// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { filter, firstValueFrom, from, interval, map, switchMap, tap } from 'rxjs';
import { getExperimentExecution, executeExperiment as sendExecuteExperimentHttpRequest } from '../api';

import { Options } from './types';

export async function executeExperiment(options: Options): Promise<boolean | undefined> {
  const result = await executeSingle(options);
  return result;
}

async function executeSingle(options: Options): Promise<boolean | undefined> {
  const experimentKey = options.key;
  console.log('Executing experiment:', experimentKey);

  const result = await sendExecuteExperimentHttpRequest(options.key);
  console.log('Experiment execution:', result.location);

  if (!options.wait || !result.location) {
    return undefined;
  }

  return await firstValueFrom(
    interval(5000)
      .pipe(switchMap(() => from(getExperimentExecution(result.location ?? ''))))
      .pipe(tap(e => console.log('Current execution state:', e.state.toLowerCase())))
      .pipe(filter(e => e.state === 'FAILED' || e.state === 'CANCELED' || e.state === 'COMPLETED'))
      .pipe(map(() => true))
  );
}

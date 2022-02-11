import { red, bold } from 'colors/safe';

import { ResilienceScoreServiceState } from './types';
import { abortExecutionWithError } from '../errors';
import { loadServiceDefinition } from './loading';
import { executeApiCall } from '../api';

export async function verify(serviceDefinitionPath: string) {
  const serviceDefinition = await loadServiceDefinition(serviceDefinitionPath);

  try {
    const response = await executeApiCall({
      method: 'GET',
      path: `/api/resilience-score/service/${serviceDefinition.id}`,
    });
    const state = (await response.json()) as ResilienceScoreServiceState;
    printState(state);

    if (state.desiredResilienceLevel !== state.actualResilienceLevel) {
      process.exit(1);
    }
  } catch (e) {
    throw abortExecutionWithError(e, 'Failed to read state for service %s', serviceDefinition.id);
  }
}

function printState(state: ResilienceScoreServiceState) {
  const desiredLevelAchieved = state.desiredResilienceLevel === state.actualResilienceLevel;

  let color: typeof red = s => s;
  if (!desiredLevelAchieved) {
    color = red;
  }

  console.log(color(bold(`Resilience Level`)));
  console.log(color(bold('================')));
  console.log('Desired: %s', state.desiredResilienceLevel);
  console.log(color('Actual:  %s'), state.actualResilienceLevel);
}

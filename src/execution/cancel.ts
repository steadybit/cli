// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { cancelExecution } from './api.ts';

export interface Options {
  id: number;
}

export async function cancel(options: Options) {
  const outcome = await cancelExecution(options.id);
  if (outcome === 'accepted') {
    console.log('Experiment run %s is being canceled.', options.id);
  } else {
    console.log('Experiment run %s has already ended.', options.id);
  }
}

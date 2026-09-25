// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { output } from '../structuredFiles.ts';
import { fetchExecution } from './api.ts';

export interface Options {
  id: number;
  file?: string;
  type?: string;
}

export async function getExecution(options: Options) {
  const execution = await fetchExecution(options.id);
  await output(execution, options);
  if (options.file) {
    console.log('Experiment run %s written to %s.', options.id, options.file);
  }
}

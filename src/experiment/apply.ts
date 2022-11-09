// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { loadExperiment, writeExperiment } from './files';
import { updateExperiment, upsertExperiment } from './api';

export interface Options {
  key?: string;
  file: string;
}

export async function applyExperiment(options: Options) {
  const experiment = await loadExperiment(options.file);

  const key = options.key || experiment.key;
  if (key) {
    await updateExperiment(key, experiment);
    console.log('Experiment %s updated.', key);
  } else {
    const result = await upsertExperiment(experiment);
    if (result.created) {
      await writeExperiment(options.file, {key: result.key, ...experiment});
      console.log('Experiment %s created.', result.key);
    } else {
      console.log('Experiment %s updated.', result.key);
    }
  }
}

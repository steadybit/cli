// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { loadExperiment, resolveExperimentFiles, writeYamlFile } from './files';
import { updateExperiment, upsertExperiment } from './api';
import { abortExecution } from '../errors';

export interface Options {
  key?: string;
  file: string[];
  recursive: boolean;
}

export async function applyExperiments(options: Options) {
  const files = await resolveExperimentFiles(options.file, options.recursive);

  if (options.key && files.length > 1) {
    throw abortExecution('If --key is specified, at most one --file can be specified.');
  }

  for (const file of files) {
    const experiment = await loadExperiment(file);
    const key = options.key || experiment.key;

    console.log('key: ', key, 'file: ', file);


    if (key) {
      await updateExperiment(key, experiment);
      console.log('Experiment %s updated.', key);
    } else {
      const result = await upsertExperiment(experiment);
      if (result.created) {
        await writeYamlFile(file, { key: result.key, ...experiment });
        console.log('Experiment %s created.', result.key);
      } else {
        console.log('Experiment %s updated.', result.key);
      }
    }
  }
}

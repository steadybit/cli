// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import yaml from 'js-yaml';
import { writeYamlFile } from './files';
import { fetchExperiment } from './api';

export interface Options {
  key: string;
  file?: string;
}

export async function getExperiment(options: Options) {
  const experiment = await fetchExperiment(options.key);

  if (!options.file) {
    console.log(yaml.dump(experiment));
  } else {
    await writeYamlFile(options.file, experiment);
    console.log('Experiment %s written to %s.', options.key, options.file);
  }
}

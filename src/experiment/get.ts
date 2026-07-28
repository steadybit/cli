// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { dump } from '../yaml.ts';
import { type Datatype, writeFile } from './files.ts';
import { fetchExperiment } from './api.ts';

export interface Options {
  key: string;
  file?: string;
  type?: Datatype;
}

export async function getExperiment(options: Options) {
  const experiment = await fetchExperiment(options.key);
  const datatype: Datatype = options.type ? options.type : options.file?.endsWith('.json') ? 'json' : 'yaml';

  if (!options.file) {
    if (datatype === 'json') {
      console.log(experiment);
    } else {
      console.log(dump(experiment));
    }
  } else {
    await writeFile(options.file, experiment, datatype);
    console.log('Experiment %s written to %s.', options.key, options.file);
  }
}

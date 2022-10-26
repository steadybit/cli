// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Experiment } from './types';
import fs from 'fs/promises';
import yaml from 'js-yaml';
import { PolicyBinding } from '../policyBinding/types';
import { abortExecution } from '../errors';

export async function writeExperiment(file: string, experiment: Experiment): Promise<void> {
  await fs.writeFile(file, yaml.dump(experiment), { encoding: 'utf8' });
}

export async function loadExperiment(file: string): Promise<Experiment> {
  let fileContent: string;
  try {
    fileContent = await fs.readFile(file, { encoding: 'utf8' });
  } catch (e) {
    throw abortExecution(
      'Failed to read experiment file at path \'%s\': %s',
      file,
      (e as Error)?.message ?? 'Unknown Cause'
    );
  }

  try {
    return yaml.load(fileContent) as PolicyBinding;
  } catch (e) {
    throw abortExecution(
      'Failed to parse experiment file at path \'%s\' as YAML/JSON: %s',
      file,
      (e as Error)?.message ?? 'Unknown Cause'
    );
  }
}

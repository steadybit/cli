// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Experiment } from './types';
import fs from 'fs/promises';
import yaml from 'js-yaml';
import { abortExecution } from '../errors';
import path from 'path';

export type Datatype = 'json' | 'yaml';

export interface ExperimentFromFile {
  experiment: Experiment;
  datatype: Datatype;
}

export async function resolveExperimentFiles(files: string[], recursive: boolean): Promise<string[]> {
  const results = [];

  for (const file of files) {
    try {
      const stat = await fs.stat(file);

      if (stat.isDirectory()) {
        const subDirectories = [];
        for (const entry of await fs.readdir(file, { withFileTypes: true })) {
          if (entry.isDirectory()) {
            subDirectories.push(path.join(file, entry.name));
          } else if (
            entry.isFile() &&
            (entry.name.toLowerCase().endsWith('.yaml') || entry.name.toLowerCase().endsWith('.yml'))
          ) {
            results.push(path.join(file, entry.name));
          }
        }
        if (recursive && subDirectories.length > 0) {
          results.push(...(await resolveExperimentFiles(subDirectories, recursive)));
        }
      } else {
        results.push(file);
      }
    } catch (e: any) {
      if (e.code === 'ENOENT') {
        throw abortExecution(`File or directory '${file}' not found.`);
      } else {
        throw e;
      }
    }
  }

  return results;
}

export async function writeFile(file: string, content: unknown, datatype: Datatype): Promise<void> {
  await fs.writeFile(file, datatype === 'json' ? JSON.stringify(content) : yaml.dump(content), { encoding: 'utf8' });
}

export async function loadExperiment(file: string): Promise<ExperimentFromFile> {
  let fileContent: string;
  try {
    fileContent = await fs.readFile(file, { encoding: 'utf8' });
  } catch (e) {
    throw abortExecution(
      "Failed to read experiment file at path '%s': %s",
      file,
      (e as Error)?.message ?? 'Unknown Cause'
    );
  }

  try {
    const experiment = JSON.parse(fileContent) as Experiment;
    return { experiment, datatype: 'json' };
  } catch {
    try {
      const experiment = yaml.load(fileContent) as Experiment;
      return { experiment, datatype: 'yaml' };
    } catch (e) {
      throw abortExecution(
        "Failed to parse experiment file at path '%s' as YAML/JSON: %s",
        file,
        (e as Error)?.message ?? 'Unknown Cause'
      );
    }
  }
}

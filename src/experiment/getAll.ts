// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import fs from 'fs/promises';
import {getAllTeams} from '../team/get';
import {Team} from '../team/types';
import {fetchExecutionsForExperiment, fetchExperiment, fetchExperiments, getExperimentExecution} from './api';
import {writeYamlFile} from './files';


export interface Options {
  directory: string;
}

export async function getAllExperiments(options: Options) {
  await ensureDirectoryExists(options.directory)
  let totalExperiments = 0;
  let totalExecutions = 0;

  for (const team of await getAllTeams(false)) {
    process.stdout.write(`Fetching experiments for team ${team.name} (${team.key})... `);
    const {countExperiments, countExecutions} = await getAllExperimentsForTeam(team, options.directory);
    totalExperiments += countExperiments;
    totalExecutions += countExecutions;
    process.stdout.write(`experiments: ${countExperiments}, executions: ${countExecutions}\n`);
  }
  console.log(`Written ${totalExperiments} experiments with ${totalExecutions} executions`);
}

async function getAllExperimentsForTeam(team: Team, dir: string) {
  const response = await fetchExperiments(team.key);

  let countExecutions = 0;
  await Promise.all(response.experiments.map(async (item) => {
    const subdir = `${dir}/${item.key}`
    await ensureDirectoryExists(subdir)

    const experiment = await fetchExperiment(item.key)
    await writeYamlFile(`${subdir}/experiment.yaml`, experiment);
    countExecutions += await getAllExecutionsForExperiment(item.key, subdir);
  }))

  return {
    countExperiments: response.experiments.length,
    countExecutions
  }
}

async function getAllExecutionsForExperiment(key: string, dir: string) {
  const response = await fetchExecutionsForExperiment(key);
  const res: number[] = await Promise.all(response.executions.map(async (item) => {
    try {
      const execution = await getExperimentExecution(item.id, false);
      await writeYamlFile(`${dir}/execution-${item.id}.yaml`, execution);
      return 1;
    } catch (e) {
      return 0;
    }
  }))
  return res.reduce((a, b) => a + b, 0);
}

async function ensureDirectoryExists(dir: string) {
  await fs.mkdir(dir, {recursive: true});
}



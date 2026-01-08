// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import fs from 'fs/promises';
import { getAllTeams } from '../team/get';
import { Team } from '../team/types';
import { fetchExecutionsForExperiment, fetchExperiment, fetchExperiments, getExperimentExecution } from './api';
import { Datatype, writeFile } from './files';

export interface Options {
  directory: string;
  type?: Datatype;
}

export async function dump(options: Options) {
  await ensureDirectoryExists(options.directory);
  let totalExperiments = 0;
  let totalExecutions = 0;

  for (const team of await getAllTeams(false)) {
    process.stdout.write(`Fetching experiments for team ${team.name} (${team.key})... `);
    const { countExperiments, countExecutions } = await getAllExperimentsForTeam(
      team,
      options.directory,
      options.type ?? 'yaml'
    );
    totalExperiments += countExperiments;
    totalExecutions += countExecutions;
    process.stdout.write(`experiments: ${countExperiments}, executions: ${countExecutions}\n`);
  }
  console.log(`Written ${totalExperiments} experiments with ${totalExecutions} executions`);
}

function removeDeprecatedFields(experiment: Record<string, any>) {
  if (Array.isArray(experiment.lanes)) {
    experiment.lanes.forEach(lane => {
      if (Array.isArray(lane.steps)) {
        lane.steps.forEach((step: Record<string, any>) => {
          if (step && typeof step === 'object' && 'radius' in step) {
            delete step.radius.query;
            delete step.radius.list;
          }
        });
      }
    });
  }

  return experiment;
}

async function getAllExperimentsForTeam(team: Team, dir: string, datatype: Datatype) {
  const response = await fetchExperiments(team.key);

  const c = await Promise.all(
    response.experiments.map(async item => {
      const subdir = `${dir}/${item.key}`;
      await ensureDirectoryExists(subdir);

      const experiment = await fetchExperiment(item.key);
      await writeFile(
        `${subdir}/experiment.${datatype === 'json' ? 'json' : 'yaml'}`,
        removeDeprecatedFields(experiment),
        datatype
      );
      return await getAllExecutionsForExperiment(item.key, subdir, datatype);
    })
  );

  return {
    countExperiments: response.experiments.length,
    countExecutions: c.reduce((a, b) => a + b, 0),
  };
}

async function getAllExecutionsForExperiment(key: string, dir: string, datatype: Datatype) {
  const response = await fetchExecutionsForExperiment(key);
  await new Promise(resolve => setTimeout(resolve, 1000));
  const c: number[] = await Promise.all(
    response.executions.map(async item => {
      let count = 0;
      try {
        if (count++ % 100 == 0) {
          //poor man's throttling
          await new Promise(resolve => setTimeout(resolve, 1000));
        }
        const execution = await getExperimentExecution(item.id, false);
        await writeFile(`${dir}/execution-${item.id}.${datatype === 'json' ? 'json' : 'yaml'}`, execution, datatype);
        return 1;
      } catch {
        return 0;
      }
    })
  );
  return c.reduce((a, b) => a + b, 0);
}

async function ensureDirectoryExists(dir: string) {
  await fs.mkdir(dir, { recursive: true });
}

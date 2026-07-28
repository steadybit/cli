// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import fs from 'node:fs/promises';
import { mapWithConcurrency } from '../concurrency.ts';
import { errorMessage } from '../errors.ts';
import { rateLimiter } from '../api/rateLimit.ts';
import { getAllTeams } from '../team/get.ts';
import { fetchExecutionsForExperiment, fetchExperiment, fetchExperiments, getExperimentExecution } from './api.ts';
import { type Datatype, writeFile } from './files.ts';
import type { ExecutionList, ExperimentList } from './types.ts';

export interface Options {
  directory: string;
  type?: Datatype;
}

// A dump walks every experiment of every team and every execution of every experiment.
// Both levels are bounded so the request volume stays predictable instead of scaling
// with tenant size, rather than relying on a connection cap: the global fetch, unlike
// the node-fetch agents it replaced, imposes no per-origin limit of its own. The
// resulting ceiling is in the same range as the 64 sockets that pool used to allow.
const EXPERIMENT_CONCURRENCY = 4;
const EXECUTION_CONCURRENCY = 16;

const LARGE_DUMP_EXPERIMENTS = 100;

export async function dump(options: Options) {
  await ensureDirectoryExists(options.directory);
  let totalExperiments = 0;
  let totalExecutions = 0;
  let totalFailedExperiments = 0;
  let totalFailedExecutions = 0;

  // The experiment lists are fetched up front, which costs nothing extra because each
  // team needs one anyway, so that the size of the walk is known before it starts.
  const teams = await getAllTeams(false);
  const listPerTeam = new Map<string, ExperimentList>();
  for (const team of teams) {
    listPerTeam.set(team.key, await fetchExperiments(team.key));
  }
  warnAboutLargeDump([...listPerTeam.values()].reduce((total, list) => total + list.experiments.length, 0));

  for (const team of teams) {
    process.stdout.write(`Fetching experiments for team ${team.name} (${team.key})... `);
    const teamDump = await getAllExperimentsForTeam(
      listPerTeam.get(team.key)!,
      options.directory,
      options.type ?? 'yaml'
    );
    totalExperiments += teamDump.countExperiments;
    totalExecutions += teamDump.countExecutions;
    totalFailedExperiments += teamDump.failedExperiments;
    totalFailedExecutions += teamDump.failedExecutions;

    const failed = teamDump.failedExperiments + teamDump.failedExecutions;
    process.stdout.write(
      `experiments: ${teamDump.countExperiments}, executions: ${teamDump.countExecutions}` +
        `${failed > 0 ? `, failed: ${failed}` : ''}\n`
    );
    // Only once the progress line above is terminated, so that the two streams stay
    // readable when they are redirected to different places.
    for (const problem of teamDump.problems) {
      console.error(`  ${problem}`);
    }
  }
  console.log(`Written ${totalExperiments} experiments with ${totalExecutions} executions`);

  if (totalFailedExperiments > 0 || totalFailedExecutions > 0) {
    // Everything that could be fetched has been written; the non-zero status is what
    // stops a pipeline treating an incomplete dump as a complete one. Executions count
    // for this too — a dump missing half its runs is not a complete dump either.
    console.error(
      `Incomplete: ${totalFailedExperiments} experiments and ${totalFailedExecutions} executions could not be dumped`
    );
    process.exitCode = 1;
  }
}

// Every experiment costs at least two requests, its design and its execution list, and
// each execution one more. Past the burst the platform meters those out slowly, so a
// large dump is a long job and saying so up front beats discovering it an hour later.
function warnAboutLargeDump(countExperiments: number): void {
  if (countExperiments <= LARGE_DUMP_EXPERIMENTS) {
    return;
  }
  const minimumMinutes = Math.ceil(rateLimiter.millisFor(countExperiments * 2) / 60000);
  console.error(
    `Dumping ${countExperiments} experiments. Requests are paced to the platform's rate limit, ` +
      `so this takes at least ${minimumMinutes} minutes, longer with executions.\n`
  );
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

export interface TeamDump {
  countExperiments: number;
  countExecutions: number;
  failedExperiments: number;
  failedExecutions: number;
  // Returned rather than printed, so the caller can finish its progress line first.
  // Emitting them as they happen tore the two streams across each other.
  problems: string[];
}

async function getAllExperimentsForTeam(response: ExperimentList, dir: string, datatype: Datatype): Promise<TeamDump> {
  // A single unlucky request must not discard the whole walk. Failures are counted and
  // reported at both levels instead of ending the command on the spot, and an execution
  // that could not be fetched is as much a hole in the dump as an experiment is.
  const results = await mapWithConcurrency(response.experiments, EXPERIMENT_CONCURRENCY, async item => {
    const subdir = `${dir}/${item.key}`;
    try {
      // The design, the execution list and the directory are independent of one another,
      // so waiting for them in turn would put two extra round trips on the critical path
      // of every experiment.
      const [, experiment, executions] = await Promise.all([
        ensureDirectoryExists(subdir),
        fetchExperiment(item.key, false),
        fetchExecutionsForExperiment(item.key, false),
      ]);

      await writeFile(`${subdir}/experiment.${datatype}`, removeDeprecatedFields(experiment), datatype);
      const { written, failed } = await writeExecutions(executions.executions, subdir, datatype);
      return {
        countExecutions: written,
        failedExperiments: 0,
        failedExecutions: failed,
        problems: failed > 0 ? [`${item.key}: ${failed} of ${written + failed} executions could not be fetched`] : [],
      };
    } catch (e) {
      return {
        countExecutions: 0,
        failedExperiments: 1,
        failedExecutions: 0,
        problems: [`${item.key}: ${errorMessage(e)}`],
      };
    }
  });

  const sum = (pick: (r: (typeof results)[number]) => number) => results.reduce((total, r) => total + pick(r), 0);
  const failedExperiments = sum(r => r.failedExperiments);
  return {
    countExperiments: response.experiments.length - failedExperiments,
    countExecutions: sum(r => r.countExecutions),
    failedExperiments,
    failedExecutions: sum(r => r.failedExecutions),
    problems: results.flatMap(r => r.problems),
  };
}

async function writeExecutions(executions: ExecutionList['executions'], dir: string, datatype: Datatype) {
  const written = await mapWithConcurrency(executions, EXECUTION_CONCURRENCY, async item => {
    try {
      const execution = await getExperimentExecution(item.id, false);
      await writeFile(`${dir}/execution-${item.id}.${datatype}`, execution, datatype);
      return true;
    } catch {
      return false;
    }
  });
  return { written: written.filter(Boolean).length, failed: written.filter(ok => !ok).length };
}

async function ensureDirectoryExists(dir: string) {
  await fs.mkdir(dir, { recursive: true });
}

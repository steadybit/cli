// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import fs from 'node:fs/promises';
import path from 'node:path';
import { Table } from 'console-table-printer';
import type { Schemas } from '../api/schemas.ts';
import { abortExecution } from '../errors.ts';
import { downloadArtifact, type Execution, fetchExecution } from './api.ts';

export interface Artifact {
  step: string;
  target: string;
  targetExecutionId: string;
  artifactId: string;
}

type TargetExecution = Schemas['TargetExecutionAO'];

// Artifacts hang off the target executions of action steps, and of the actions a
// service validation step runs on the step's behalf. The platform offers no listing of
// its own, so they are collected from the run.
export function collectArtifacts(execution: Execution): Artifact[] {
  const artifacts: Artifact[] = [];
  const addFrom = (step: string, targetExecutions: TargetExecution[] | undefined) => {
    for (const target of targetExecutions ?? []) {
      for (const artifactId of target.artifacts ?? []) {
        artifacts.push({ step, target: target.name ?? '', targetExecutionId: target.id ?? '', artifactId });
      }
    }
  };

  for (const step of execution.steps ?? []) {
    const label = step.customLabel || ('actionId' in step && step.actionId) || step.stepType;
    if ('targetExecutions' in step) {
      addFrom(label, step.targetExecutions);
    }
    if ('validations' in step) {
      for (const validation of step.validations ?? []) {
        addFrom(validation.customLabel || validation.actionId || label, validation.targetExecutions);
      }
    }
  }
  return artifacts;
}

export interface ListOptions {
  id: number;
}

export async function listArtifacts(options: ListOptions) {
  const artifacts = collectArtifacts(await fetchExecution(options.id));
  if (artifacts.length === 0) {
    console.log('Experiment run %s has no artifacts.', options.id);
    return;
  }
  const table = new Table({
    columns: [
      { name: 'artifactId', title: 'Artifact', alignment: 'left' },
      { name: 'target', title: 'Target', alignment: 'left' },
      { name: 'step', title: 'Step', alignment: 'left' },
      { name: 'targetExecutionId', title: 'Target execution', alignment: 'left' },
    ],
  });
  table.addRows(artifacts);
  table.printTable();
}

// The ids come from the platform, and are reduced to a single path segment so that one
// containing "../" cannot write outside the chosen directory. basename alone is not
// enough: it leaves ".." as it is.
export function pathSegment(id: string): string {
  const segment = path.basename(id);
  return segment === '' || segment === '.' || segment === '..' ? '_' : segment;
}

export interface DownloadOptions {
  id: number;
  artifact?: string;
  targetExecution?: string;
  directory: string;
  output?: string;
}

// Every artifact lands in <directory>/<target execution>/<artifact>. Two targets of the
// same step typically produce artifacts of the same name, so the name alone would let
// one download overwrite another.
export async function downloadArtifacts(options: DownloadOptions) {
  const selected = collectArtifacts(await fetchExecution(options.id)).filter(
    a =>
      (!options.artifact || a.artifactId === options.artifact) &&
      (!options.targetExecution || a.targetExecutionId === options.targetExecution)
  );

  if (selected.length === 0) {
    throw abortExecution('No matching artifacts found in experiment run %s.', options.id);
  }
  if (options.output && selected.length > 1) {
    throw abortExecution(
      '%d artifacts match, but --output takes exactly one. Narrow it down with --artifact and --target-execution.',
      selected.length
    );
  }

  for (const artifact of selected) {
    const file =
      options.output ??
      path.join(options.directory, pathSegment(artifact.targetExecutionId), pathSegment(artifact.artifactId));
    const content = await downloadArtifact(options.id, artifact.targetExecutionId, artifact.artifactId);
    await fs.mkdir(path.dirname(file), { recursive: true });
    await fs.writeFile(file, content);
    console.log('Artifact %s written to %s.', artifact.artifactId, file);
  }
}

// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { TaskDefinition, PolicyDefinition } from '../policyBinding/types';
import { forEachFile, IteratorArgs, ForEachAction } from './util';
import { getGitHubRepositoryName } from '../vcs/github';
import { abortExecution } from '../errors';

interface Options {
  version: string;
  directory: string;
}

export async function setVersion(options: Options) {
  const gitHubRepositoryName = await getGitHubRepositoryName(options.directory);
  if (!gitHubRepositoryName) {
    const error = await abortExecution(
      "Failed to identify GitHub reposity name from `git remote -v` output for directory '%s'.",
      options.directory
    );
    throw error;
  }

  await forEachFile<TaskDefinition>({
    pattern: '**/task.?(yml|yaml)',
    options: {
      nocase: true,
      cwd: options.directory,
    },
    iterator: args => taskIterator(options, args),
  });

  await forEachFile<PolicyDefinition>({
    pattern: '**/policy.?(yml|yaml)',
    options: {
      nocase: true,
      cwd: options.directory,
    },
    iterator: args => policyIterator(options, args, gitHubRepositoryName),
  });
}

async function taskIterator(options: Options, { fileContent }: IteratorArgs<TaskDefinition>): Promise<ForEachAction> {
  let action: ForEachAction;

  // Some task authors do not maintain a version number within the task definition file.
  if (fileContent.version) {
    fileContent.version = options.version;
    action = 'update';
  }

  return action;
}

async function policyIterator(
  options: Options,
  { fileContent }: IteratorArgs<PolicyDefinition>,
  gitHubRepositoryName: string
): Promise<ForEachAction> {
  let action: ForEachAction;

  // Some policy authors do not maintain a version number within the policy definition file.
  if (fileContent.version) {
    fileContent.version = options.version;
    action = 'update';
  }

  for (const task of fileContent.tasks) {
    if (task.name.startsWith(gitHubRepositoryName)) {
      task.version = options.version;
      action = 'update';
    }
  }

  return action;
}

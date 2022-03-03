/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import { TaskDefinition, PolicyDefinition } from '../serviceDefinition/types';
import { abortExecution, abortExecutionWithError } from '../errors';
import { forEachFile, IteratorArgs, ForEachAction } from './util';
import { executeShellCommand } from '../shell';

interface Options {
  version: string;
  cwd: string;
}

export async function setVersion(options: Options) {
  const gitHubRepositoryName = await getGitHubRepositoryName(options.cwd);
  if (!gitHubRepositoryName) {
    const error = await abortExecution(
      "Failed to identify GitHub reposity name from `git remote -v` output for CWD '%s'.",
      options.cwd
    );
    throw error;
  }

  await forEachFile<TaskDefinition>({
    pattern: '**/task.?(yml|yaml)',
    options: {
      nocase: true,
      cwd: options.cwd,
    },
    iterator: args => taskIterator(options, args),
  });

  await forEachFile<PolicyDefinition>({
    pattern: '**/policy.?(yml|yaml)',
    options: {
      nocase: true,
      cwd: options.cwd,
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

// exported for testing purposes
export async function getGitHubRepositoryName(cwd: string): Promise<string | undefined> {
  try {
    const { stdout } = await executeShellCommand('git remote -v', {
      cwd,
    });

    return getGitHubRepositoryNameFromGitRemotes(stdout);
  } catch (e) {
    const error = await abortExecutionWithError(e, 'Failed to identify GitHub repository name.');
    throw error;
  }
}

// exported for testing purposes
export function getGitHubRepositoryNameFromGitRemotes(gitRemotes: string): string | undefined {
  return gitRemotes.match(/github\.com[:/]([^/]+\/[^/ ]+).git/i)?.[1];
}

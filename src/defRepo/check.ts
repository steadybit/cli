// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import colors from 'colors/safe';
import path from 'path';

import { TaskDefinition, PolicyDefinition } from '../policyBinding/types';
import { getAbsolutePathToRepositoryRoot, getGitSha } from '../vcs/git';
import { abortExecution, abortExecutionWithError } from '../errors';
import { getGitHubRepositoryName } from '../vcs/github';
import { forEachFile, IteratorArgs } from './util';
import { executeApiCall } from '../api/http';
import { Response } from 'node-fetch';

interface Options {
  version: string;
  directory: string;
}

export async function check(options: Options) {
  const messages = await executeCheck(options);

  if (messages.length > 0) {
    console.log(messages);

    const error = await abortExecution(colors.red('Found issues while checking task and policy definitions. See above.'));
    throw error;
  } else {
    console.log(colors.green('No issues found!'));
  }
}

export async function executeCheck(options: Options): Promise<string> {
  if (!options.version) {
    try {
      options.version = await getGitSha(options.directory);
    } catch (e) {
      const error = await abortExecutionWithError(e, '-v, --version option not specified and git sha could not be identified automatically. Please specify a version.');
      throw error;
    }
  }

  const repositoryName = await getGitHubRepositoryName(options.directory);
  if (!repositoryName) {
    const error = await abortExecution(
      "Failed to identify GitHub reposity name from `git remote -v` output for directory '%s'.",
      options.directory
    );
    throw error;
  }

  const repositoryRoot = await getAbsolutePathToRepositoryRoot(options.directory);

  const messages: string[] = [];

  await forEachFile<TaskDefinition>({
    pattern: '**/task.?(yml|yaml)',
    options: {
      nocase: true,
      cwd: options.directory,
    },
    iterator: async args => {
      const message = await taskIterator(repositoryRoot, repositoryName, options, args);
      messages.push(message);
    },
  });

  await forEachFile<PolicyDefinition>({
    pattern: '**/policy.?(yml|yaml)',
    options: {
      nocase: true,
      cwd: options.directory,
    },
    iterator: async args => {
      const message = await policyIterator(repositoryRoot, repositoryName, options, args);
      messages.push(message);
    },
  });

  return messages.join('');
}

async function taskIterator(
  repositoryRoot: string,
  repositoryName: string,
  options: Options,
  { fileName }: IteratorArgs<TaskDefinition>
): Promise<string> {
  const taskReference = path.relative(repositoryRoot, fileName);
  const fullyQualifiedTaskName = `${repositoryName}/${taskReference}`.replace(/\/task\.ya?ml$/i, '');

  let message = '';

  try {
    const response: Response = await executeApiCall({
      method: 'get',
      path: '/api/task-definitions',
      queryParameters: {
        name: fullyQualifiedTaskName,
        version: options.version
      }
    });

    const taskDefinition = await response.json();

    if (taskDefinition.name !== fullyQualifiedTaskName) {
      message += ' - Name in task definition file does not match the fully qualified task name.\n';
      message += `   Got:      ${taskDefinition.name}\n`;
      message += `   Expected: ${fullyQualifiedTaskName}\n`;
    }
  } catch (error: any) {
    message += ' - Task retrieval/validation failed. Details:\n';

    if (typeof error?.response?.json === 'function') {
      try {
        const body = await (error.response as Response).json();
        message += JSON.stringify(body, undefined, 2);
      } catch (e: any) {
        message += ` - Unhandled error while trying to provide details: ${e?.message}\n`;
      }
    }
  }

  if (message.length > 0) {
    message = `${path.relative(repositoryRoot, fileName)}:\n${message}\n\n`;
  }

  return message;
}

async function policyIterator(
  repositoryRoot: string,
  repositoryName: string,
  options: Options,
  { fileName }: IteratorArgs<PolicyDefinition>
): Promise<string> {
  const policyReference = path.relative(repositoryRoot, fileName);
  const fullyQualifiedPolicyName = `${repositoryName}/${policyReference}`.replace(/\/policy\.ya?ml$/i, '');

  let message = '';

  try {
    const response: Response = await executeApiCall({
      method: 'get',
      path: '/api/policy-definitions',
      queryParameters: {
        name: fullyQualifiedPolicyName,
        version: options.version
      }
    });

    const policyDefinition = await response.json();

    if (policyDefinition.name !== fullyQualifiedPolicyName) {
      message += ' - Name in policy definition file does not match the fully qualified policy name.\n';
      message += `   Got:      ${policyDefinition.name}\n`;
      message += `   Expected: ${fullyQualifiedPolicyName}\n`;
    }
  } catch (error: any) {
    message += ' - Policy retrieval/validation failed. Details:\n';

    if (typeof error?.response?.json === 'function') {
      try {
        const body = await (error.response as Response).json();
        message += JSON.stringify(body, undefined, 2);
      } catch (e: any) {
        message += ` - Unhandled error while trying to provide details: ${e?.message}\n`;
      }
    }
  }

  if (message.length > 0) {
    message = `${path.relative(repositoryRoot, fileName)}:\n${message}\n\n`;
  }

  return message;
}

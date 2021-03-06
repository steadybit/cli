// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { abortExecutionWithError } from '../errors';
import { getGitOutput } from './git';

export async function   getGitHubRepositoryName(cwd: string): Promise<string | undefined> {
  try {
    const stdout = await getGitOutput({command: 'git remote -v', cwd});
    return getGitHubRepositoryNameFromGitRemotes(stdout);
  } catch (e) {
    const error = await abortExecutionWithError(e, 'Failed to identify GitHub repository name.');
    throw error;
  }
}

// exported for testing purposes
export function getGitHubRepositoryNameFromGitRemotes(gitRemotes: string): string | undefined {
  return gitRemotes.match(/github\.com[:/]([^/]+\/[^/ ]+).git/i)?.[1] ?? gitRemotes.match(/https:\/\/github\.com[:/]([^/]+\/[^/ ]+)/i)?.[1];
}

// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import path from 'path';

import { executeShellCommand } from '../shell';

export async function getGitRemote(filePath: string): Promise<string> {
  const output = await getGitOutput({
    command: 'git remote -v',
    cwd: path.dirname(filePath),
  });
  return output.split('\n')[0].split(/\s+/)[1];
}

export async function getPathToFileInGitRepository(filePath: string): Promise<string> {
  const output = await getGitOutput({
    command: 'git rev-parse --show-toplevel',
    cwd: path.dirname(filePath),
  });
  return path.relative(output, filePath);
}

export async function getAbsolutePathToRepositoryRoot(cwd: string): Promise<string> {
  return getGitOutput({
    command: 'git rev-parse --show-toplevel',
    cwd,
  });
}

export async function getGitRef(directory: string): Promise<string> {
  return getGitOutput({
    command: 'git symbolic-ref --short HEAD',
    cwd: directory,
  });
}

export async function getGitSha(directory: string): Promise<string> {
  return getGitOutput({
    command: 'git rev-parse --verify HEAD',
    cwd: directory,
  });
}

export interface GetGitOutputArgs {
  command: string;
  cwd: string;
}

export async function getGitOutput({ command, cwd }: GetGitOutputArgs): Promise<string> {
  try {
    const { stdout } = await executeShellCommand(command, {
      cwd,
    });

    return stdout.trim();
  } catch (e) {
    throw new Error(
      `Failed to execute command "${command}" in CWD "${cwd}".`,
      // TypeScript does not know about error causes yet.
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-expect-error
      { cause: e }
    );
  }
}

export async function getTags(policyBindingPath: string): Promise<Record<string, string>> {
  const records = await Promise.all([
    getAsTag(
      'git.sha',
      getGitSha(path.dirname(policyBindingPath)),
    ),
    getAsTag(
      'git.ref',
      getGitRef(path.dirname(policyBindingPath)),
    ),
    getAsTag(
      'git.repository_path',
      getPathToFileInGitRepository(policyBindingPath),
    ),
    getAsTag(
      'git.remote',
      getGitRemote(path.dirname(policyBindingPath)),
    )
  ]);

  return Object.fromEntries(records.flatMap(record => Object.entries(record)));
}

async function getAsTag(
  key: string,
  valuePromise: Promise<string>
): Promise<Record<string, string>> {
  try {
    const value = await valuePromise;

    if (value.length > 0) {
      return {
        [key]: value,
      };
    }
  } catch (e) {
    // We do not actually know whether we are operating inside a git repository. This makes it
    // a guess and for a guess we aren't going to bother the end-user.
  }

  return {};
}

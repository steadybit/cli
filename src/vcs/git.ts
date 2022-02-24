/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import path from 'path';

import { executeShellCommand } from '../shell';

export async function getTags(serviceDefinitionPath: string): Promise<Record<string, string>> {
  const records = await Promise.all([
    getGitOutputAsTag(serviceDefinitionPath, 'git rev-parse --verify HEAD', 'git.sha'),
    getGitOutputAsTag(serviceDefinitionPath, 'git symbolic-ref --short HEAD', 'git.ref'),
    getGitOutputAsTag(serviceDefinitionPath, 'git rev-parse --show-toplevel', 'git.repository_path', stdout =>
      path.relative(stdout, serviceDefinitionPath)
    ),
    getGitOutputAsTag(serviceDefinitionPath, 'git remote -v', 'git.remote', stdout =>
      stdout.split('\n')[0].split(/\s+/)[1]
    ),
  ]);

  return Object.fromEntries(records.flatMap(record => Object.entries(record)));
}

async function getGitOutputAsTag(
  serviceDefinitionPath: string,
  command: string,
  tag: string,
  processValue: (stdout: string) => string = s => s
): Promise<Record<string, string>> {
  try {
    let { stdout } = await executeShellCommand(command, {
      cwd: path.dirname(serviceDefinitionPath),
    });

    stdout = processValue(stdout);

    if (stdout.length > 0) {
      return {
        [tag]: stdout,
      };
    }
  } catch (e) {
    // We do not actually know whether we are operating inside a git repository. This makes it
    // a guess and for a guess we aren't going to bother the end-user.
  }

  return {};
}

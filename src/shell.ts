// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { exec as execWithoutPromises } from 'child_process';
import { promisify } from 'util';

const exec = promisify(execWithoutPromises);

export interface ShellCommandResult {
  stdout: string;
  stderr: string;
}

export interface ShellCommandOptions {
  trimOutput?: boolean;
  cwd?: string;
  timeout?: number;
}

export async function executeShellCommand(
  command: string,
  { trimOutput = true, cwd, timeout = 60_000 }: ShellCommandOptions = {}
): Promise<ShellCommandResult> {
  const result = await exec(command, {
    cwd,
    timeout,
  });

  if (trimOutput) {
    result.stdout = result.stdout.trim();
    result.stderr = result.stderr.trim();
  }

  return result;
}

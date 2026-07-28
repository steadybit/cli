// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// inquirer used to re-raise SIGINT so that Ctrl+C terminated the CLI quietly. The
// @inquirer prompts instead reject with an ExitPromptError, which would otherwise
// surface as an unhandled rejection with a stack trace. This restores the quiet exit
// with the conventional SIGINT status.
const SIGINT_EXIT_CODE = 130;

export async function cancelable<T>(prompt: Promise<T>): Promise<T> {
  try {
    return await prompt;
  } catch (e) {
    if (e instanceof Error && e.name === 'ExitPromptError') {
      process.exit(SIGINT_EXIT_CODE);
    }
    throw e;
  }
}

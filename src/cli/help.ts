// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import type { Command } from 'commander';

// Every command shows at least one example under --help. An option list says what can
// be passed, not which combination does the thing a pipeline author came for.
export function withExamples(command: Command, examples: string[]): Command {
  const lines = examples.map(example => `  $ ${example}`).join('\n');
  return command.addHelpText('after', `\nExamples:\n${lines}\n`);
}

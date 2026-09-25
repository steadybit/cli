// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { InvalidArgumentError } from 'commander';

// Not `parseInt` directly: commander invokes an argument parser as (value, previous),
// so passing it wholesale turns the previous value into the radix. Repeating an option
// then parses the second value in the base of the first, silently and wrongly.
export function parseDecimal(value: string): number {
  const parsed = Number.parseInt(value, 10);
  if (Number.isNaN(parsed)) {
    throw new InvalidArgumentError(`'${value}' is not a number.`);
  }
  return parsed;
}

// For repeatable `--flag KEY=VALUE` options. Only the first `=` separates, so a value
// may itself contain one, as a URL with a query string does.
export function collectKeyValue(value: string, previous: Record<string, string> = {}): Record<string, string> {
  const separator = value.indexOf('=');
  if (separator <= 0) {
    throw new InvalidArgumentError(`'${value}' is not in the form KEY=VALUE.`);
  }
  return { ...previous, [value.slice(0, separator)]: value.slice(separator + 1) };
}

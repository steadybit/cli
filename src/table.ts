// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { Table } from 'console-table-printer';
import { colorsSupported } from './colors.ts';

type TableOptions = NonNullable<ConstructorParameters<typeof Table>[0]>;

// console-table-printer colours its output on its own, whatever stdout is, which put
// escape codes into every listing piped into grep or a file. It follows the same terminal
// check as the rest of the CLI instead.
export function createTable(options: Exclude<TableOptions, string[]> = {}): Table {
  return new Table({ ...options, shouldDisableColors: !colorsSupported });
}

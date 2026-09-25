// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { describe, expect, it } from 'vitest';
import { createTable } from './table.ts';

describe('table', () => {
  // The test runner's stdout is not a terminal, just as a pipe is not.
  it('renders without escape codes when stdout is not a terminal', () => {
    const table = createTable({ columns: [{ name: 'id', title: 'Id' }] });
    table.addRow({ id: 'abc' }, { color: 'red' });

    const rendered = table.render();

    expect(rendered).toContain('abc');
    expect(rendered).not.toContain('\x1b[');
  });
});

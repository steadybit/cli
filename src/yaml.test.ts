// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { describe, expect, it } from 'vitest';
import { dump, load } from './yaml.ts';

describe('yaml', () => {
  it('should expand merge keys instead of keeping a literal << property', () => {
    const parsed = load(`
defaults: &defaults
  ignoreFailure: false
lanes:
  - steps:
      - <<: *defaults
        type: action
`) as any;

    expect(parsed.lanes[0].steps[0]).toEqual({ ignoreFailure: false, type: 'action' });
    expect(parsed.lanes[0].steps[0]).not.toHaveProperty('<<');
  });

  it('should parse timestamps into dates', () => {
    const parsed = load('created: 2024-01-15T10:30:00Z') as any;

    expect(parsed.created).toBeInstanceOf(Date);
  });

  it('should keep quoted scalars as strings', () => {
    const parsed = load("graceful: 'true'") as any;

    expect(parsed.graceful).toBe('true');
  });

  it('should round-trip dumped documents', () => {
    const original = { key: 'TST-1', lanes: [{ steps: [{ type: 'action', ignoreFailure: false }] }] };

    expect(load(dump(original))).toEqual(original);
  });
});

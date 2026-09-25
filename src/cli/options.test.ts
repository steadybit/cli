// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { describe, expect, it } from 'vitest';
import { collectKeyValue, parseDecimal } from './options.ts';

describe('options', () => {
  it('collects repeated KEY=VALUE pairs, splitting on the first =', () => {
    const first = collectKeyValue('url=http://shop?a=b');
    expect(collectKeyValue('CLUSTER=prod', first)).toEqual({ url: 'http://shop?a=b', CLUSTER: 'prod' });
  });

  it('accepts an empty value but not an empty key', () => {
    expect(collectKeyValue('EMPTY=')).toEqual({ EMPTY: '' });
    expect(() => collectKeyValue('=value')).toThrow("'=value' is not in the form KEY=VALUE.");
    expect(() => collectKeyValue('novalue')).toThrow("'novalue' is not in the form KEY=VALUE.");
  });

  it('parses decimals regardless of the previous value', () => {
    expect(parseDecimal('010')).toBe(10);
    expect(() => parseDecimal('ten')).toThrow("'ten' is not a number.");
  });
});

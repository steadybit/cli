// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { describe, expect, it } from 'vitest';
import { mapWithConcurrency } from './concurrency.ts';

function trackingMapper(limitObserver: { inFlight: number; maxInFlight: number }) {
  return async (item: number) => {
    limitObserver.inFlight++;
    limitObserver.maxInFlight = Math.max(limitObserver.maxInFlight, limitObserver.inFlight);
    await new Promise(resolve => setImmediate(resolve));
    limitObserver.inFlight--;
    return item * 2;
  };
}

describe('mapWithConcurrency', () => {
  it('should never exceed the given limit', async () => {
    const observer = { inFlight: 0, maxInFlight: 0 };

    await mapWithConcurrency(
      Array.from({ length: 100 }, (_, i) => i),
      4,
      trackingMapper(observer)
    );

    expect(observer.maxInFlight).toBe(4);
  });

  it('should return results in input order', async () => {
    const items = [5, 1, 4, 2, 3];

    const results = await mapWithConcurrency(items, 2, async item => {
      await new Promise(resolve => setTimeout(resolve, item));
      return item * 10;
    });

    expect(results).toEqual([50, 10, 40, 20, 30]);
  });

  it('should handle fewer items than the limit', async () => {
    const observer = { inFlight: 0, maxInFlight: 0 };

    const results = await mapWithConcurrency([1, 2], 16, trackingMapper(observer));

    expect(results).toEqual([2, 4]);
    expect(observer.maxInFlight).toBe(2);
  });

  it('should handle an empty input', async () => {
    expect(await mapWithConcurrency([], 4, async () => 1)).toEqual([]);
  });

  it('should reject when the mapper fails', async () => {
    await expect(
      mapWithConcurrency([1, 2, 3], 2, async item => {
        if (item === 2) {
          throw new Error('boom');
        }
        return item;
      })
    ).rejects.toThrow('boom');
  });
});

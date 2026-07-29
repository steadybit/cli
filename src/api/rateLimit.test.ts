// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { describe, expect, it } from 'vitest';
import { type Clock, RateLimiter, bucketFromEnvironment, defaultBucket } from './rateLimit.ts';

// Wall-clock assertions made these tests fail whenever the machine was busy. This clock
// advances only when the limiter asks to sleep, so the numbers below are exact.
function fakeClock(): Clock & { elapsed(): number } {
  let millis = 0;
  return {
    now: () => millis,
    sleep: async requested => {
      millis += requested;
    },
    elapsed: () => millis,
  };
}

describe('RateLimiter', () => {
  it('should let the whole burst through without pacing', async () => {
    const clock = fakeClock();
    const limiter = new RateLimiter({ burst: 20, refillTokens: 1, refillIntervalMillis: 10000 }, clock);

    for (let i = 0; i < 20; i++) {
      await limiter.acquire();
    }

    expect(clock.elapsed()).toEqual(0);
  });

  it('should pace at the refill rate once the burst is spent', async () => {
    const clock = fakeClock();
    const limiter = new RateLimiter({ burst: 4, refillTokens: 1, refillIntervalMillis: 100 }, clock);

    for (let i = 0; i < 7; i++) {
      await limiter.acquire();
    }

    // Four free, then one per 100ms interval for the remaining three.
    expect(clock.elapsed()).toEqual(300);
  });

  it('should bound concurrent callers, not just sequential ones', async () => {
    const clock = fakeClock();
    const limiter = new RateLimiter({ burst: 3, refillTokens: 1, refillIntervalMillis: 200 }, clock);

    const admittedAt: number[] = [];
    await Promise.all(Array.from({ length: 6 }, () => limiter.acquire().then(() => admittedAt.push(clock.now()))));

    expect(admittedAt).toEqual([0, 0, 0, 200, 400, 600]);
  });

  it('should refill over time rather than all at once', async () => {
    const clock = fakeClock();
    const limiter = new RateLimiter({ burst: 2, refillTokens: 2, refillIntervalMillis: 200 }, clock);
    await limiter.acquire();
    await limiter.acquire();

    // A full interval restores the burst, and no more than the burst.
    await clock.sleep(200);
    await limiter.acquire();
    await limiter.acquire();
    expect(clock.elapsed()).toEqual(200);

    await limiter.acquire();
    expect(clock.elapsed()).toBeGreaterThan(200);
  });

  it('should not accumulate more than the burst while idle', async () => {
    const clock = fakeClock();
    const limiter = new RateLimiter({ burst: 2, refillTokens: 2, refillIntervalMillis: 100 }, clock);

    await clock.sleep(10_000); // idle far longer than it takes to refill

    await limiter.acquire();
    await limiter.acquire();
    expect(clock.elapsed()).toEqual(10_000);

    await limiter.acquire();
    expect(clock.elapsed()).toBeGreaterThan(10_000);
  });

  describe('millisFor', () => {
    it('should charge nothing for a run inside the burst', () => {
      expect(new RateLimiter().millisFor(100)).toEqual(0);
    });

    it('should charge the refill rate beyond the burst', () => {
      // 200 requests is 100 beyond the burst, which is four refills of 15s.
      expect(new RateLimiter().millisFor(200)).toEqual(60000);
    });
  });

  describe('bucketFromEnvironment', () => {
    it('should use the documented allowance when nothing is set', () => {
      expect(bucketFromEnvironment({})).toEqual(defaultBucket);
    });

    it('should take overrides, with the interval given in seconds', () => {
      expect(
        bucketFromEnvironment({
          STEADYBIT_RATE_LIMIT_BURST: '10',
          STEADYBIT_RATE_LIMIT_REFILL: '5',
          STEADYBIT_RATE_LIMIT_INTERVAL: '30',
        })
      ).toEqual({ burst: 10, refillTokens: 5, refillIntervalMillis: 30000 });
    });

    it('should take a value with surrounding whitespace', () => {
      expect(bucketFromEnvironment({ STEADYBIT_RATE_LIMIT_BURST: ' 10 ' }).burst).toEqual(10);
    });

    // `Number` would have read these as 1000, 16 and 2.5. None is how a request count
    // gets written on purpose, and accepting them changes the pacing silently.
    it.each(['0', '-1', 'abc', 'NaN', '1e3', '0x10', '2.5', '10,5', '+5'])(
      'should fall back and complain about %s',
      value => {
        const complaints: string[] = [];
        const original = console.error;
        console.error = (m: string) => complaints.push(m);
        try {
          expect(bucketFromEnvironment({ STEADYBIT_RATE_LIMIT_BURST: value }).burst).toEqual(defaultBucket.burst);
        } finally {
          console.error = original;
        }
        expect(complaints).toHaveLength(1);
        expect(complaints[0]).toContain('STEADYBIT_RATE_LIMIT_BURST');
      }
    );
  });

  it('should default to the allowance the platform documents', () => {
    // A burst of 100, refilling by 25 every 15s.
    expect(defaultBucket).toEqual({ burst: 100, refillTokens: 25, refillIntervalMillis: 15000 });
  });
});

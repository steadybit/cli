// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { setTimeout as sleep } from 'node:timers/promises';

// The platform limits requests with a token bucket: a burst of 100 is allowed, and the
// allowance then refills by 25 every 15 seconds. A fan-out like `experiment dump` issues
// far more than that, and because every rejected request retries into the window it just
// exhausted, the retries make the pressure worse rather than better.
//
// Pacing to the documented allowance from the first request keeps 429s rare instead of
// routine. Nothing is learned from response headers: `RateLimit-Limit: 100;w=15` reports
// the burst but not the refill rate, and taking it at face value paces four times too
// fast. Where the limit is exceeded anyway, the Retry-After handling in http.ts adapts.
export const defaultBucket = {
  burst: 100,
  refillTokens: 25,
  refillIntervalMillis: 15000,
};

export type BucketOptions = typeof defaultBucket;

// Injected so that tests can drive the bucket deterministically. Asserting on wall-clock
// windows made them fail whenever the machine was busy.
export interface Clock {
  now(): number;
  sleep(millis: number): Promise<void>;
}

const systemClock: Clock = {
  now: () => Date.now(),
  sleep: millis => sleep(millis),
};

// A deployment may be configured with a different allowance, and the platform does not
// advertise one that could be read instead: the ratelimit-* headers appear only on the
// 429 itself, by which point remaining is zero.
export function bucketFromEnvironment(env: NodeJS.ProcessEnv = process.env): BucketOptions {
  return {
    burst: positiveInteger(env.STEADYBIT_RATE_LIMIT_BURST, 'STEADYBIT_RATE_LIMIT_BURST', defaultBucket.burst),
    refillTokens: positiveInteger(
      env.STEADYBIT_RATE_LIMIT_REFILL,
      'STEADYBIT_RATE_LIMIT_REFILL',
      defaultBucket.refillTokens
    ),
    refillIntervalMillis:
      positiveInteger(
        env.STEADYBIT_RATE_LIMIT_INTERVAL,
        'STEADYBIT_RATE_LIMIT_INTERVAL',
        defaultBucket.refillIntervalMillis / 1000
      ) * 1000,
  };
}

// Plain decimal digits only. `Number` would also have taken '1e3', '0x10' and '2.5',
// which are not how anyone means to write a request count, and reading '0x10' as 16
// would quietly change how hard the CLI hits the platform.
const POSITIVE_INTEGER = /^\d+$/;

function positiveInteger(value: string | undefined, name: string, fallback: number): number {
  if (value === undefined || value.trim() === '') {
    return fallback;
  }
  const trimmed = value.trim();
  if (!POSITIVE_INTEGER.test(trimmed) || Number(trimmed) <= 0) {
    // Warned about rather than ignored: a typo here silently changes how hard the CLI
    // hits the platform, which is the last thing that should fail quietly.
    console.error(`Ignoring ${name}: '${value}' is not a positive whole number. Using ${fallback}.`);
    return fallback;
  }
  return Number(trimmed);
}

export class RateLimiter {
  private tokens: number;
  private lastRefill: number;
  private gate: Promise<unknown> = Promise.resolve();

  constructor(
    private readonly bucket: BucketOptions = defaultBucket,
    private readonly clock: Clock = systemClock
  ) {
    this.tokens = bucket.burst;
    this.lastRefill = clock.now();
  }

  // How long `count` requests take once the burst is spent, which is what makes the
  // scale of a large dump visible before it starts rather than an hour into it.
  millisFor(count: number): number {
    const beyondBurst = Math.max(0, count - this.bucket.burst);
    return (beyondBurst / this.bucket.refillTokens) * this.bucket.refillIntervalMillis;
  }

  // Serialised, so that concurrent callers cannot all spend the same token.
  acquire(): Promise<void> {
    const admitted = this.gate.then(() => this.reserve());
    this.gate = admitted.catch(() => undefined);
    return admitted;
  }

  private async reserve(): Promise<void> {
    for (;;) {
      this.refill();
      if (this.tokens >= 1) {
        this.tokens--;
        return;
      }
      await this.clock.sleep(this.millisUntilNextToken());
    }
  }

  private refill(): void {
    const now = this.clock.now();
    this.tokens = Math.min(this.bucket.burst, this.tokens + (now - this.lastRefill) * this.tokensPerMilli());
    this.lastRefill = now;
  }

  private tokensPerMilli(): number {
    return this.bucket.refillTokens / this.bucket.refillIntervalMillis;
  }

  private millisUntilNextToken(): number {
    return Math.max(1, Math.ceil((1 - this.tokens) / this.tokensPerMilli()));
  }
}

let shared: RateLimiter | undefined;

// Built on first use, not at import. The parent `steadybit` process pulls this module in
// through http.ts just to dispatch to a subcommand and never makes a request, so eager
// construction had it read the environment — and complain about a bad value — twice.
function sharedLimiter(): RateLimiter {
  return (shared ??= new RateLimiter(bucketFromEnvironment()));
}

export const rateLimiter = {
  acquire: () => sharedLimiter().acquire(),
  millisFor: (count: number) => sharedLimiter().millisFor(count),
};

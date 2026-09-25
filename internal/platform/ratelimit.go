// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package platform

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// The platform limits requests with a token bucket: a burst of 100, refilled by 25 every
// 15 seconds. A fan-out like `experiment dump` issues far more than that, and every
// rejected request would retry into the window it just exhausted, so requests are paced
// to the documented allowance from the first one. The ratelimit-* headers cannot be used
// instead: they appear only on the 429 itself, and report the burst but not the refill.
type Bucket struct {
	Burst          int
	RefillTokens   int
	RefillInterval time.Duration
}

var DefaultBucket = Bucket{Burst: 100, RefillTokens: 25, RefillInterval: 15 * time.Second}

// Plain decimal digits only: '1e3' or '0x10' are not how anyone writes a request count.
var positiveInteger = regexp.MustCompile(`^\d+$`)

// BucketFromEnvironment reads the STEADYBIT_RATE_LIMIT_* overrides. An invalid value is
// warned about rather than ignored quietly, since it would change how hard the CLI hits
// the platform.
func BucketFromEnvironment() Bucket {
	read := func(name string, fallback int) int {
		value, ok := os.LookupEnv(name)
		trimmed := strings.TrimSpace(value)
		if !ok || trimmed == "" {
			return fallback
		}
		n, err := strconv.Atoi(trimmed)
		if !positiveInteger.MatchString(trimmed) || err != nil || n <= 0 {
			fmt.Fprintf(os.Stderr, "Ignoring %s: '%s' is not a positive whole number. Using %d.\n", name, value, fallback)
			return fallback
		}
		return n
	}
	return Bucket{
		Burst:          read("STEADYBIT_RATE_LIMIT_BURST", DefaultBucket.Burst),
		RefillTokens:   read("STEADYBIT_RATE_LIMIT_REFILL", DefaultBucket.RefillTokens),
		RefillInterval: time.Duration(read("STEADYBIT_RATE_LIMIT_INTERVAL", int(DefaultBucket.RefillInterval/time.Second))) * time.Second,
	}
}

// Clock is injected so tests can drive the bucket deterministically.
type Clock interface {
	Now() time.Time
	Sleep(time.Duration)
}

type systemClock struct{}

func (systemClock) Now() time.Time        { return time.Now() }
func (systemClock) Sleep(d time.Duration) { time.Sleep(d) }

type RateLimiter struct {
	mu         sync.Mutex
	bucket     Bucket
	clock      Clock
	tokens     float64
	lastRefill time.Time
}

func NewRateLimiter(bucket Bucket, clock Clock) *RateLimiter {
	if clock == nil {
		clock = systemClock{}
	}
	return &RateLimiter{bucket: bucket, clock: clock, tokens: float64(bucket.Burst), lastRefill: clock.Now()}
}

// Acquire blocks until a request may be sent. Callers are served one at a time, so
// concurrent ones cannot all spend the same token.
func (r *RateLimiter) Acquire() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for {
		r.refill()
		if r.tokens >= 1 {
			r.tokens--
			return
		}
		r.clock.Sleep(r.untilNextToken())
	}
}

// DurationFor is how long `count` requests take once the burst is spent, which makes
// the scale of a large dump visible before it starts rather than an hour into it.
func (r *RateLimiter) DurationFor(count int) time.Duration {
	beyond := max(0, count-r.bucket.Burst)
	return time.Duration(float64(beyond) / float64(r.bucket.RefillTokens) * float64(r.bucket.RefillInterval))
}

func (r *RateLimiter) perNanosecond() float64 {
	return float64(r.bucket.RefillTokens) / float64(r.bucket.RefillInterval)
}

func (r *RateLimiter) refill() {
	now := r.clock.Now()
	r.tokens = math.Min(float64(r.bucket.Burst), r.tokens+float64(now.Sub(r.lastRefill))*r.perNanosecond())
	r.lastRefill = now
}

func (r *RateLimiter) untilNextToken() time.Duration {
	return max(time.Millisecond, time.Duration(math.Ceil((1-r.tokens)/r.perNanosecond())))
}

var (
	sharedLimiter     *RateLimiter
	sharedLimiterOnce sync.Once
)

// SetLimiter replaces the shared limiter. Tests use it so that a suite's requests are
// not paced to the platform's allowance.
func SetLimiter(l *RateLimiter) {
	sharedLimiterOnce.Do(func() {})
	sharedLimiter = l
}

// Limiter is built on first use, so that a command that sends nothing never reads, or
// complains about, the environment.
func Limiter() *RateLimiter {
	sharedLimiterOnce.Do(func() { sharedLimiter = NewRateLimiter(BucketFromEnvironment(), nil) })
	return sharedLimiter
}

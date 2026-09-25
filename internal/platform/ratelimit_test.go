// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package platform

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type fakeClock struct{ now time.Time }

func (c *fakeClock) Now() time.Time        { return c.now }
func (c *fakeClock) Sleep(d time.Duration) { c.now = c.now.Add(d) }

func TestSpendsTheBurstWithoutWaiting(t *testing.T) {
	clock := &fakeClock{now: time.Unix(0, 0)}
	limiter := NewRateLimiter(Bucket{Burst: 3, RefillTokens: 1, RefillInterval: time.Second}, clock)

	for range 3 {
		limiter.Acquire()
	}

	assert.Equal(t, time.Unix(0, 0), clock.now)
}

func TestPacesToTheRefillOnceTheBurstIsSpent(t *testing.T) {
	clock := &fakeClock{now: time.Unix(0, 0)}
	limiter := NewRateLimiter(Bucket{Burst: 2, RefillTokens: 25, RefillInterval: 15 * time.Second}, clock)

	for range 2 + 25 {
		limiter.Acquire()
	}

	assert.InDelta(t, 15*time.Second, clock.now.Sub(time.Unix(0, 0)), float64(10*time.Millisecond))
}

func TestEstimatesTheDurationOfALargeWalk(t *testing.T) {
	limiter := NewRateLimiter(DefaultBucket, &fakeClock{})

	assert.Equal(t, time.Duration(0), limiter.DurationFor(100))
	assert.Equal(t, 60*time.Second, limiter.DurationFor(200))
}

func TestReadsOverridesAndIgnoresInvalidOnes(t *testing.T) {
	t.Setenv("STEADYBIT_RATE_LIMIT_BURST", " 50 ")
	t.Setenv("STEADYBIT_RATE_LIMIT_REFILL", "1e3")
	t.Setenv("STEADYBIT_RATE_LIMIT_INTERVAL", "0x10")

	assert.Equal(t, Bucket{Burst: 50, RefillTokens: 25, RefillInterval: 15 * time.Second}, BucketFromEnvironment())
}

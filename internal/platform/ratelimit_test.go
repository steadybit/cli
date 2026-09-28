// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

package platform

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeClock struct{ now time.Time }

func (c *fakeClock) Now() time.Time { return c.now }
func (c *fakeClock) Sleep(ctx context.Context, d time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	c.now = c.now.Add(d)
	return nil
}

func TestSpendsTheBurstWithoutWaiting(t *testing.T) {
	clock := &fakeClock{now: time.Unix(0, 0)}
	limiter := NewRateLimiter(Bucket{Burst: 3, RefillTokens: 1, RefillInterval: time.Second}, clock)

	for range 3 {
		require.NoError(t, limiter.Wait(context.Background()))
	}

	assert.Equal(t, time.Unix(0, 0), clock.now)
}

func TestPacesToTheRefillOnceTheBurstIsSpent(t *testing.T) {
	clock := &fakeClock{now: time.Unix(0, 0)}
	limiter := NewRateLimiter(Bucket{Burst: 2, RefillTokens: 25, RefillInterval: 15 * time.Second}, clock)

	for range 2 + 25 {
		require.NoError(t, limiter.Wait(context.Background()))
	}

	assert.InDelta(t, 15*time.Second, clock.now.Sub(time.Unix(0, 0)), float64(10*time.Millisecond))
}

// A caller that gives up, as completion does after its timeout, does not spend a token.
func TestACanceledWaitHandsItsTokenBack(t *testing.T) {
	clock := &fakeClock{now: time.Unix(0, 0)}
	limiter := NewRateLimiter(Bucket{Burst: 1, RefillTokens: 1, RefillInterval: time.Second}, clock)
	require.NoError(t, limiter.Wait(context.Background()))
	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	assert.ErrorIs(t, limiter.Wait(canceled), context.Canceled)
	require.NoError(t, limiter.Wait(context.Background()))
	assert.Equal(t, time.Second, clock.now.Sub(time.Unix(0, 0)))
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

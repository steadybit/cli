// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Package interrupt decides what Ctrl-C and SIGTERM do. By default they end the CLI with
// 130 or 143 and no stack trace. Code that must clean up first, such as `run --wait`
// cancelling the experiment it started, pushes a handler for as long as it runs.
package interrupt

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

var (
	mu       sync.Mutex
	handlers []*func(os.Signal)
	// interrupted is set before the handlers run: something started after RunHandlers
	// took its snapshot has no handler that will clean it up, so it asks Interrupted.
	interrupted atomic.Bool
)

func init() {
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		first := <-signals
		// A second signal while cleaning up ends the CLI at once: someone pressing
		// Ctrl-C twice wants out, whatever the cleanup was doing.
		go func() {
			<-signals
			os.Exit(code(first))
		}()
		RunHandlers(first)
		fmt.Fprintln(os.Stderr)
		os.Exit(code(first))
	}()
}

// RunHandlers runs what is registered for sig, most recent first, without exiting. The
// signal handler calls it before exiting; tests call it to stand in for Ctrl-C.
func RunHandlers(sig os.Signal) {
	interrupted.Store(true)
	mu.Lock()
	pending := make([]*func(os.Signal), len(handlers))
	copy(pending, handlers)
	mu.Unlock()
	for i := len(pending) - 1; i >= 0; i-- {
		(*pending[i])(sig)
	}
}

// Interrupted reports whether the CLI was interrupted. Code that starts something to
// clean up checks it after pushing its handler, which may have come too late to run.
func Interrupted() bool {
	return interrupted.Load()
}

// Reset forgets an interrupt, for tests that stand in for one with RunHandlers.
func Reset() {
	interrupted.Store(false)
}

func code(s os.Signal) int {
	if s == syscall.SIGTERM {
		return 143
	}
	return 130
}

// Push runs fn when the CLI is interrupted, before it exits, until the returned pop is
// called. Handlers run most recent first.
func Push(fn func(os.Signal)) (pop func()) {
	mu.Lock()
	defer mu.Unlock()
	h := &fn
	handlers = append(handlers, h)
	return func() {
		mu.Lock()
		defer mu.Unlock()
		for i, existing := range handlers {
			if existing == h {
				handlers = append(handlers[:i], handlers[i+1:]...)
				return
			}
		}
	}
}

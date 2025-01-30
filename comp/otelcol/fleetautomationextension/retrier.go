// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

package fleetautomationextension

import (
	"context"
	"math/rand"
	"time"
)

// BackoffStrategy defines how backoffRetrier will wait between retries, and how many times
// to retry. The first retry is random between initialWait[0] and initialWait[1].
// Then for each retry: newWait = previousWait * multiplier + jitter
// If MaxWait is not zero, then the wait will never be above MaxWait.
type BackoffStrategy struct {
	InitialWait []time.Duration
	Multiplier  float64
	Jitter      []time.Duration
	MaxWait     time.Duration
}

// A RetriableCode advises whether a retrier should retry.
type RetriableCode int

// Enumeration of available RetriableCodes
const (
	ShouldRetry = iota
	DoNotRetry
	ForceKeepRetrying
)

type backoffRetrier struct {
	Strategy    BackoffStrategy
	Times       int
	Retriable   func(error) RetriableCode // Optional: func(error) shouldRetryCode
	Notifer     func(error, time.Duration)
	numAttempts int

	lastBackoff time.Duration
}

// NextRetry gives you the next deadline when you should retry
func (r *backoffRetrier) NextRetry() time.Time {
	var waitTime time.Duration

	if r.lastBackoff == 0 { // No previous wait, randomize
		diff := int64(r.Strategy.InitialWait[1]) - int64(r.Strategy.InitialWait[0])
		waitTime = r.Strategy.InitialWait[0]
		if diff != 0 {
			waitTime += time.Duration(rand.Int63n(diff))
		}
	} else {
		jitterDiff := int64(r.Strategy.Jitter[1]) - int64(r.Strategy.Jitter[0])
		var jitter time.Duration
		if jitterDiff != 0 {
			jitter = r.Strategy.Jitter[0] + time.Duration(rand.Int63n(jitterDiff))
		}

		waitTime = time.Duration(float64(r.lastBackoff) * r.Strategy.Multiplier)
		waitTime += jitter
	}

	if waitTime > r.Strategy.MaxWait && r.Strategy.MaxWait != 0 {
		waitTime = r.Strategy.MaxWait
	}

	r.lastBackoff = waitTime
	return time.Now().Add(waitTime)
}

// Reset resets the backoff to initial values
func (r *backoffRetrier) Reset() {
	r.lastBackoff = 0
}

// Do runs of your function x times until it succeeds with backoff delays
func (r *backoffRetrier) Do(f func() error) (err error) {
	return r.DoContext(context.Background(), f)
}

// DoContext runs of your function x times until it succeeds with backoff delays or until context expires
func (r *backoffRetrier) DoContext(ctx context.Context, f func() error) (err error) {
	return r.doContext(ctx, func(context.Context) error {
		return f()
	})
}

// DoContextV2 runs of your function x times until it succeeds with backoff delays or until context expires
// This version allows you to pass a function that takes a context as an argument so APM spans can be connected
func (r *backoffRetrier) doContext(ctx context.Context, f func(context.Context) error) (err error) {

	// Don't do anything if context is already expired
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// default to 3 attempts if it's set to 0
	totalNumAttempts := r.Times
	if totalNumAttempts == 0 {
		totalNumAttempts = 3
	}
	for i := 0; i < totalNumAttempts; i++ {
		r.numAttempts++
		if err = f(ctx); err == nil {
			return nil
		}
		if r.Retriable != nil {
			switch r.Retriable(err) {
			case ShouldRetry:
			case DoNotRetry:
				return err
			case ForceKeepRetrying:
				i--
			default:
				return err
			}
		}
		if i+1 >= totalNumAttempts {
			// if we exhausted retries, break from the for loop
			break
		}

		//nolint:S1024
		sleepDuration := time.Until(r.NextRetry())
		if r.Notifer != nil {
			r.Notifer(err, sleepDuration)
		}
	}
	return err
}

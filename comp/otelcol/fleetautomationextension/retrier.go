// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

package fleetautomationextension

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// Retry calls fn at most n times, retrying on error
func Retry(n int, fn func() error) error {
	rt := Retrier{Times: n}
	return rt.Do(fn)
}

// RetryContext calls fn at most n times, retrying on error
func RetryContext(ctx context.Context, n int, fn func(fctx context.Context) error) error {
	rt := Retrier{Times: n}
	return rt.DoContext(ctx, fn)
}

// Retrier is a helper to retry a function x times
// if Retriable is provided it will only retry if
// Retriable returns true for the error
type Retrier struct {
	// Times is the maximum number of times the function will be called
	Times     int
	Retriable func(error) bool
}

// Do runs your function x times until it succeeds
func (r *Retrier) Do(f func() error) (err error) {
	for i := 0; i < r.Times; i++ {
		err = f()
		if err != nil && (r.Retriable == nil || r.Retriable(err)) {
			continue
		}
		break
	}
	return err
}

func (r *Retrier) DoContext(ctx context.Context, f func(fctx context.Context) error) (err error) {
	for i := 0; i < r.Times; i++ {
		err = f(ctx)
		if err != nil && (r.Retriable == nil || r.Retriable(err)) {
			continue
		}
		break
	}
	return err
}

// BackoffStrategy defines how BackoffRetrier will wait between retries, and how many times
// to retry. The first retry is random between initialWait[0] and initialWait[1].
// Then for each retry: newWait = previousWait * multiplier + jitter
// If MaxWait is not zero, then the wait will never be above MaxWait.
type BackoffStrategy struct {
	InitialWait []time.Duration
	Multiplier  float64
	Jitter      []time.Duration
	MaxWait     time.Duration
}

// IsZero checks if the strategy has not been initialized
func (s *BackoffStrategy) IsZero() bool {
	return s.Multiplier == 0
}

// GetDefaultBackoffStrategy gives you a good default backoff strategy
// (retry spaced by seconds, max 5 min wait)
func GetDefaultBackoffStrategy() BackoffStrategy {
	return BackoffStrategy{
		[]time.Duration{4 * time.Second, 8 * time.Second}, // First retry, wait between 4 and 8s
		1.5,
		[]time.Duration{0, 2 * time.Second},
		5 * time.Minute,
	}
}

// A RetriableCode advises whether a retrier should retry.
type RetriableCode int

// Enumeration of available RetriableCodes
const (
	ShouldRetry = iota
	DoNotRetry
	ForceKeepRetrying
)

func (r RetriableCode) String() string {
	switch r {
	case ShouldRetry:
		return "ShouldRetry"
	case DoNotRetry:
		return "DoNotRetry"
	case ForceKeepRetrying:
		return "ForceKeepRetrying"
	default:
		return "Unknown"
	}
}

// BackoffRetrier is a simple structure to help you implement a exponential backoff retrier
// You should call your function, and if it fails, call the backoff.NextRetry()
// to get the next time you should retry
// You should call Reset() when you function succeed, and you want to reset the timer
// Usage:
//
//	b := BackoffRetrier{
//	  Strategy: util.GetDefaultBackoffStrategy(),
//	}
//
//	if time.Now().Before(deadline) {
//	  continue
//	}
//
// resp, err := myFunction()
//
//	if err != nil {
//	  // Handle error
//	  deadline = b.NextRetry()
//	} else {
//
//	  b.Reset()
//	}
type BackoffRetrier struct {
	Strategy    BackoffStrategy
	Times       int
	Retriable   func(error) RetriableCode // Optional: func(error) shouldRetryCode
	Notifer     func(error, time.Duration)
	numAttempts int

	lastBackoff time.Duration
}

// NextRetry gives you the next deadline when you should retry
func (r *BackoffRetrier) NextRetry() time.Time {
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
func (r *BackoffRetrier) Reset() {
	r.lastBackoff = 0
}

// Do runs of your function x times until it succeeds with backoff delays
func (r *BackoffRetrier) Do(f func() error) (err error) {
	return r.DoContext(context.Background(), f)
}

// DoContext runs of your function x times until it succeeds with backoff delays or until context expires
func (r *BackoffRetrier) DoContext(ctx context.Context, f func() error) (err error) {
	return r.doContext(ctx, func(context.Context) error {
		return f()
	}, false)
}

func (r *BackoffRetrier) DoContextWithTracing(ctx context.Context, f func(context.Context) error) (err error) {
	return r.doContext(ctx, f, true)
}

// DoContextV2 runs of your function x times until it succeeds with backoff delays or until context expires
// This version allows you to pass a function that takes a context as an argument so APM spans can be connected
func (r *BackoffRetrier) doContext(ctx context.Context, f func(context.Context) error, withAPM bool) (err error) {

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
		sleepDuration := r.NextRetry().Sub(time.Now())
		if r.Notifer != nil {
			r.Notifer(err, sleepDuration)
		}
	}
	return err
}

func (r *BackoffRetrier) NumberRetries() int {
	return r.numAttempts - 1
}

// BackoffDuration calculates the backoff duration using full-jitter
func BackoffDuration(base, max time.Duration, attempt int, multiplier float64) time.Duration {
	scale := rand.Float64()
	next := base
	for i := 0; i < attempt; i++ {
		next = time.Duration(float64(next) * multiplier)
		if next > max {
			next = max
			break
		}
	}
	next = time.Duration(float64(next) * scale)
	return next
}

// ExponentialBackoffDuration calculates backoff duration without any jitter
// it returns min(base * (multiplier ^ attempt), max)
func ExponentialBackoffDuration(base, max time.Duration, attempt int, multiplier float64) time.Duration {
	// check if values are going to overflow
	// we need to know if base * (multiplier ^ attempt) > math.MaxInt64
	// by applying math.Log it's the same as checking if
	// log(base) + attempt*log(multiplier) > log(math.MaxInt64)
	if math.Log(float64(base))+float64(attempt)*math.Log(multiplier) > math.Log(math.MaxInt64) {
		return max
	}

	next := time.Duration(float64(base) * math.Pow(multiplier, float64(attempt)))
	if next > max {
		return max
	}
	return next
}

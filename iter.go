package backoff

import (
	"context"
	"time"
)

var after = time.After

type Backoff interface {
	Interval(n int) time.Duration
}

// Iter returns a iterator function that yields the i-th attempt and the backoff interval.
func Iter(n int, b Backoff) func(yield func(int, time.Duration) bool) {
	return func(yield func(int, time.Duration) bool) {
		for i := 0; i < n; i++ {
			if !yield(i, b.Interval(i)) {
				break
			}
		}
	}
}

// SleepIter returns a iterator function that yields the i-th attempt and sleeps for the backoff interval.
func SleepIter(n int, backoff Backoff) func(yield func(int) bool) {
	return SleepIterContext(context.Background(), n, backoff)
}

// SleepIterContext returns a iterator function that yields the i-th attempt and sleeps for the backoff interval.
// The iterator function stops the iteration when the context is canceled.
func SleepIterContext(ctx context.Context, n int, backoff Backoff) func(yield func(int) bool) {
	return func(yield func(int) bool) {
		Iter(n, backoff)(func(i int, d time.Duration) bool {
			if !yield(i) {
				return false
			}
			if i < n-1 { // skip sleep on last iteration
				select {
				case <-after(d):
				case <-ctx.Done():
					return false
				}
			}
			return true
		})
	}
}

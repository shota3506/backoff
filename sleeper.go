package backoff

import (
	"context"
	"time"
)

type Backoff interface {
	Iter(n int) func(yield func(int, time.Duration) bool)
}

type Sleeper struct {
	backoff Backoff
	after   func(d time.Duration) <-chan time.Time
}

func NewSleeper(backoff Backoff) *Sleeper {
	return &Sleeper{
		backoff: backoff,
		after:   time.After,
	}
}

func (s *Sleeper) Iter(n int) func(yield func(int) bool) {
	return s.IterContext(context.Background(), n)
}

func (s *Sleeper) IterContext(ctx context.Context, n int) func(yield func(int) bool) {
	return func(yield func(int) bool) {
		s.backoff.Iter(n)(func(i int, d time.Duration) bool {
			if !yield(i) {
				return false
			}
			if i < n-1 { // skip sleep on last iteration
				select {
				case <-s.after(d):
				case <-ctx.Done():
					return false
				}
			}
			return true
		})
	}
}

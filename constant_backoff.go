package backoff

import "time"

// ConstantBackoff is a backoff strategy that always returns the same interval.
type ConstantBackoff struct {
	interval time.Duration
}

// NewConstantBackoff returns a new ConstantBackoff with the given interval.
func NewConstantBackoff(interval time.Duration) *ConstantBackoff {
	return &ConstantBackoff{
		interval: interval,
	}
}

// Interval returns the backoff interval for the i-th attempt.
// The interval is always the same.
func (b *ConstantBackoff) Interval(n int) time.Duration {
	return b.interval
}

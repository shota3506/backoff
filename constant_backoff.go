package backoff

import "time"

type ConstantBackoff struct {
	interval time.Duration
}

func NewConstantBackoff(interval time.Duration) *ConstantBackoff {
	return &ConstantBackoff{
		interval: interval,
	}
}

func (b *ConstantBackoff) Iter(n int) func(yield func(int, time.Duration) bool) {
	return func(yield func(int, time.Duration) bool) {
		for i := 0; i < n; i++ {
			if !yield(i, b.interval) {
				break
			}
		}
	}
}

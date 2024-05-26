package backoff

import (
	"slices"
	"testing"
	"time"
)

func TestConstantBackoff(t *testing.T) {
	type result struct {
		i int
		d time.Duration
	}

	for _, tt := range []struct {
		name     string
		interval time.Duration
		expected []result
	}{
		{
			name:     "1 second",
			interval: time.Second,
			expected: []result{
				{i: 0, d: time.Second},
				{i: 1, d: time.Second},
				{i: 2, d: time.Second},
				{i: 3, d: time.Second},
				{i: 4, d: time.Second},
			},
		},
		{
			name:     "10 seconds",
			interval: 10 * time.Second,
			expected: []result{
				{i: 0, d: 10 * time.Second},
				{i: 1, d: 10 * time.Second},
				{i: 2, d: 10 * time.Second},
				{i: 3, d: 10 * time.Second},
				{i: 4, d: 10 * time.Second},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			b := NewConstantBackoff(tt.interval)

			var got []result
			for i, duration := range b.Iter(5) {
				got = append(got, result{i: i, d: duration})
			}

			if !slices.Equal(got, tt.expected) {
				t.Errorf("got: %v, expected: %v", got, tt.expected)
			}
		})
	}
}

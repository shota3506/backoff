package backoff

import (
	"math/rand/v2"
	"slices"
	"testing"
	"time"
)

func TestExponentialBackoff(t *testing.T) {
	type result struct {
		i int
		d time.Duration
	}

	for _, tt := range []struct {
		name     string
		config   ExponentialBackoffConfig
		expected []result
	}{
		{
			name: "multiplier 1 and randomization factor 0",
			config: ExponentialBackoffConfig{
				InitialInterval:     time.Second,
				Multiplier:          1,
				RandomizationFactor: 0,
			},
			expected: []result{
				{i: 0, d: time.Second},
				{i: 1, d: time.Second},
				{i: 2, d: time.Second},
				{i: 3, d: time.Second},
				{i: 4, d: time.Second},
			},
		},
		{
			name: "multiplier 2 and randomization factor 0.5",
			config: ExponentialBackoffConfig{
				InitialInterval:     time.Second,
				Multiplier:          2,
				RandomizationFactor: 0.5,
			},
			expected: []result{
				{i: 0, d: 819144156},
				{i: 1, d: 1447210901},
				{i: 2, d: 2156976007},
				{i: 3, d: 9211213750},
				{i: 4, d: 20432719687},
			},
		},
		{
			name: "exceeds max interval",
			config: ExponentialBackoffConfig{
				InitialInterval:     time.Second,
				MaxInterval:         10 * time.Second,
				Multiplier:          2,
				RandomizationFactor: 0.5,
			},
			expected: []result{
				{i: 0, d: 819144156},
				{i: 1, d: 1447210901},
				{i: 2, d: 2156976007},
				{i: 3, d: 9211213750},
				{i: 4, d: 10000000000},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewExponentialBackoff(tt.config)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			// set fixed seed for testing
			b.rand = rand.New(rand.NewChaCha8([32]byte{}))

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

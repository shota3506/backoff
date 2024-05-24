package backoff

import (
	"slices"
	"testing"
	"time"
)

func TestConstantBackoff(t *testing.T) {
	duration := time.Second
	b := NewConstantBackoff(duration)

	var got []struct {
		i int
		d time.Duration
	}

	for i, duration := range b.Iter(5) {
		got = append(got, struct {
			i int
			d time.Duration
		}{i: i, d: duration})
	}

	expected := []struct {
		i int
		d time.Duration
	}{
		{i: 0, d: time.Second},
		{i: 1, d: time.Second},
		{i: 2, d: time.Second},
		{i: 3, d: time.Second},
		{i: 4, d: time.Second},
	}
	if !slices.Equal(got, expected) {
		t.Errorf("got %v, expected %v", got, expected)
	}
}

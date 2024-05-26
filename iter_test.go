package backoff

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"testing"
	"time"
)

func TestSleeper_Iter(t *testing.T) {
	callStack := []string{}

	var mu sync.Mutex
	after = func(d time.Duration) <-chan time.Time {
		ch := make(chan time.Time, 1)
		go func() {
			now := <-time.After(time.Microsecond)
			mu.Lock()
			callStack = append(callStack, fmt.Sprintf("sleep %s", d))
			mu.Unlock()
			ch <- now
		}()
		return ch
	}
	t.Cleanup(func() {
		after = time.After
	})

	duration := time.Second
	backoff := NewConstantBackoff(duration)

	for i := range SleepIter(5, backoff) {
		callStack = append(callStack, fmt.Sprintf("yield %d", i))
	}

	expected := []string{
		"yield 0",
		"sleep 1s",
		"yield 1",
		"sleep 1s",
		"yield 2",
		"sleep 1s",
		"yield 3",
		"sleep 1s",
		"yield 4",
	}
	if !slices.Equal(callStack, expected) {
		t.Errorf("got: %v, expected: %v", callStack, expected)
	}
}

func TestSleeper_IterContext(t *testing.T) {
	callStack := []string{}

	var mu sync.Mutex
	after = func(d time.Duration) <-chan time.Time {
		ch := make(chan time.Time, 1)
		go func() {
			now := <-time.After(time.Microsecond)
			mu.Lock()
			callStack = append(callStack, fmt.Sprintf("sleep %s", d))
			mu.Unlock()
			ch <- now
		}()
		return ch
	}
	t.Cleanup(func() {
		after = time.After
	})

	duration := time.Second
	backoff := NewConstantBackoff(duration)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := range SleepIterContext(ctx, 5, backoff) {
		callStack = append(callStack, fmt.Sprintf("yield %d", i))
		if i == 2 {
			cancel()
		}
	}

	expected := []string{
		"yield 0",
		"sleep 1s",
		"yield 1",
		"sleep 1s",
		"yield 2",
	}
	if !slices.Equal(callStack, expected) {
		t.Errorf("got: %v, expected: %v", callStack, expected)
	}
}

package backoff_test

import (
	"errors"
	"log"
	"time"

	"github.com/shota3506/backoff"
)

func ExampleExponentialBackoff() {
	do := func(i int) error {
		log.Printf("do(%d): %s\n", i, time.Now())
		if i < 2 {
			return errors.New("error")
		}
		return nil
	}

	b, err := backoff.NewExponentialBackoff(backoff.ExponentialBackoffConfig{
		InitialInterval:     100 * time.Millisecond,
		RandomizationFactor: 0.5,
		Multiplier:          2,
	})
	if err != nil {
		log.Fatal(err)
	}

	for i := range backoff.SleepIter(5, b) {
		if err := do(i); err == nil {
			break
		}
		// retry after backoff
	}
}

package backoff

import (
	cryptorand "crypto/rand"
	"errors"
	"math"
	"math/rand/v2"
	"time"
)

// ExponentialBackoffConfig is the configuration for ExponentialBackoff.
type ExponentialBackoffConfig struct {
	// InitialInterval specifies the initial interval.
	InitialInterval time.Duration

	// RandomizationFactor specifies the degree of randomization.
	// The backoff interval is multiplied by a random value in the range [1-randomizationFactor, 1+randomizationFactor].
	RandomizationFactor float64

	// Multiplier specifies the multiplier for the next interval.
	// The next interval is calculated by multiplying the previous interval by the multiplier.
	Multiplier float64

	// MaxInterval specifies the maximum backoff interval.
	// The backoff interval is capped to the max interval.
	// If the value is 0, the backoff interval is not capped.
	MaxInterval time.Duration
}

// ExponentialBackoff is a backoff strategy that increases the interval exponentially.
type ExponentialBackoff struct {
	initialInterval     time.Duration
	randomizationFactor float64
	multiplier          float64
	maxInterval         time.Duration

	rand *rand.Rand
}

// NewExponentialBackoff returns a new ExponentialBackoff with the given configuration.
func NewExponentialBackoff(config ExponentialBackoffConfig) (*ExponentialBackoff, error) {
	if config.InitialInterval < 0 {
		return nil, errors.New("initial interval must be greater than or equal to 0")
	}
	if config.RandomizationFactor < 0 || config.RandomizationFactor > 1 {
		return nil, errors.New("randomization factor must be greater than or equal to 0 and less than or equal to 1")
	}
	if config.Multiplier < 0 {
		return nil, errors.New("multiplier must be greater than or equal to 0")
	}

	var seed [32]byte
	if _, err := cryptorand.Read(seed[:]); err != nil {
		return nil, err
	}
	r := rand.New(rand.NewChaCha8(seed))

	return &ExponentialBackoff{
		initialInterval:     config.InitialInterval,
		randomizationFactor: config.RandomizationFactor,
		multiplier:          config.Multiplier,
		maxInterval:         config.MaxInterval,

		rand: r,
	}, nil
}

// Interval returns the backoff interval for the i-th attempt.
// The interval is calculated by multiplying the initial interval by the multiplier to the power of i.
func (b *ExponentialBackoff) Interval(i int) time.Duration {
	interval := time.Duration(float64(b.initialInterval) * math.Pow(b.multiplier, float64(i)))
	if b.randomizationFactor > 0 {
		interval = b.randomize(interval)
	}
	if b.maxInterval > 0 {
		return min(interval, b.maxInterval)
	}
	return interval
}

func (b *ExponentialBackoff) randomize(interval time.Duration) time.Duration {
	f := float64(b.rand.Uint64()<<11>>11) / ((1 << 53) - 1) // pseudo-random number in closed interval [0, 1]
	return time.Duration(float64(interval) * (1 + (2*f-1)*b.randomizationFactor))
}

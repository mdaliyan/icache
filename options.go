package icache

import "time"

type customizablePot interface {
	setTTL(TTL time.Duration)
}

// Option is the modifier function for the pot. icache uses options pattern.
type Option func(pot customizablePot)

// WithTTL sets the global ttl invalidation strategy for the pot.
// after setting an entry, it will be removed after the given duration time.
func WithTTL(duration time.Duration) Option {
	return func(pot customizablePot) {
		pot.setTTL(duration)
	}
}

package clients

import (
	"github.com/sony/gobreaker"

	"github.com/cyber/test-project/config"
)

type CircuitBreaker interface {
	Execute(func() (any, error)) (any, error)
}

type noCircuitBreaker struct{}

func (b noCircuitBreaker) Execute(req func() (any, error)) (any, error) {
	return req()
}

func NewCircuitBreaker(cfg config.CircuitBreakerConfig) CircuitBreaker {
	cbSett := gobreaker.Settings{
		Name:        cfg.Name,
		MaxRequests: cfg.MaxRequests,
		Timeout:     cfg.Timeout,
	}

	if cfg.MaxFailures > 0 {
		cbSett.ReadyToTrip = func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= cfg.MaxFailures
		}
	}

	return gobreaker.NewCircuitBreaker(cbSett)
}

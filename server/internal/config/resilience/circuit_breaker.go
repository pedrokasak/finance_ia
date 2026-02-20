package resilience

import (
	"fmt"
	"time"

	"github.com/sony/gobreaker"
)

// CircuitBreakers holds all circuit breakers for external services
type CircuitBreakers struct {
	Stripe *gobreaker.CircuitBreaker
	AI     *gobreaker.CircuitBreaker
}

// NewCircuitBreakers initializes circuit breakers for Stripe and AI providers
func NewCircuitBreakers() *CircuitBreakers {
	return &CircuitBreakers{
		Stripe: newBreaker("stripe", 5, 30*time.Second, 10*time.Second),
		AI:     newBreaker("ai", 3, 60*time.Second, 5*time.Second),
	}
}

func newBreaker(name string, maxRequests uint32, interval, timeout time.Duration) *gobreaker.CircuitBreaker {
	st := gobreaker.Settings{
		Name:        name,
		MaxRequests: maxRequests, // half-open max requests
		Interval:    interval,    // closed state rolling window
		Timeout:     timeout,     // open -> half-open timeout
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Trip if more than 5 consecutive failures OR >50% failure rate over 5+ requests
			return counts.ConsecutiveFailures > 3 ||
				(counts.Requests >= 5 && counts.TotalFailures*2 > counts.Requests)
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			fmt.Printf("[circuit-breaker] %s: %s → %s\n", name, from, to)
		},
	}
	return gobreaker.NewCircuitBreaker(st)
}

// CallStripe executes fn within the Stripe circuit breaker
func (cb *CircuitBreakers) CallStripe(fn func() (interface{}, error)) (interface{}, error) {
	return cb.Stripe.Execute(fn)
}

// CallAI executes fn within the AI circuit breaker
func (cb *CircuitBreakers) CallAI(fn func() (interface{}, error)) (interface{}, error) {
	return cb.AI.Execute(fn)
}

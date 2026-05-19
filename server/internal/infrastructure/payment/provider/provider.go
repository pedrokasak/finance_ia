package provider

import "strings"

const (
	Stripe     = "stripe"
	AbacatePay = "abacatepay"
)

func Normalize(v string) string {
	s := strings.ToLower(strings.TrimSpace(v))
	switch s {
	case "", Stripe:
		return Stripe
	case "abacate", AbacatePay:
		return AbacatePay
	default:
		return Stripe
	}
}

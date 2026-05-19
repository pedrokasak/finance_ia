package main

import (
	subDomain "finance-ia/internal/domain/subscription"
	"fmt"
)

type stripeCatalogPlanRepo interface {
	FindBySlug(slug string) (*subDomain.Plan, error)
	Upsert(plan *subDomain.Plan) error
}

type stripeCatalogGateway interface {
	EnsurePlanCatalog(slug, name, description string, monthly, yearly float64, currency string) (productID, monthlyPriceID, yearlyPriceID string, err error)
}

func syncStripeCatalog(repo stripeCatalogPlanRepo, gateway stripeCatalogGateway) error {
	for _, slug := range []string{"pro", "premium"} {
		plan, err := repo.FindBySlug(slug)
		if err != nil {
			return err
		}
		if plan == nil {
			return fmt.Errorf("plan %s not found for stripe sync", slug)
		}

		// Skip free or malformed plans; only paid plans need Stripe catalog.
		if plan.PriceMonthly <= 0 || plan.PriceYearly <= 0 {
			continue
		}

		// Idempotent short-circuit: if IDs are already present, avoid remote calls.
		if plan.StripeProductID != "" &&
			plan.StripePriceIDMonthly != "" &&
			plan.StripePriceIDYearly != "" {
			continue
		}

		productID, monthlyID, yearlyID, err := gateway.EnsurePlanCatalog(
			plan.Slug,
			plan.Name,
			plan.Description,
			plan.PriceMonthly,
			plan.PriceYearly,
			"brl",
		)
		if err != nil {
			return err
		}

		plan.StripeProductID = productID
		plan.StripePriceIDMonthly = monthlyID
		plan.StripePriceIDYearly = yearlyID
		if err := repo.Upsert(plan); err != nil {
			return err
		}
	}
	return nil
}

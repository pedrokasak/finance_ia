package main

import (
	subDomain "finance-ia/internal/domain/subscription"
	"fmt"
)

type abacateCatalogPlanRepo interface {
	FindBySlug(slug string) (*subDomain.Plan, error)
	Upsert(plan *subDomain.Plan) error
}

type abacateCatalogGateway interface {
	EnsurePlanCatalog(slug, name, description string, monthly, yearly float64, currency string) (productID, monthlyProductID, yearlyProductID string, err error)
}

func syncAbacateCatalog(repo abacateCatalogPlanRepo, gateway abacateCatalogGateway) error {
	for _, slug := range []string{"pro", "premium"} {
		plan, err := repo.FindBySlug(slug)
		if err != nil {
			return err
		}
		if plan == nil {
			return fmt.Errorf("plan %s not found for abacatepay sync", slug)
		}
		if plan.PriceMonthly <= 0 || plan.PriceYearly <= 0 {
			continue
		}

		// Reuse monthly/yearly fields as generic checkout item IDs for non-Stripe providers.
		if plan.StripePriceIDMonthly != "" && plan.StripePriceIDYearly != "" {
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

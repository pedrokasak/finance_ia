package main

import (
	subDomain "finance-ia/internal/domain/subscription"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakePlanCatalogRepo struct {
	plans map[string]*subDomain.Plan
}

func (f *fakePlanCatalogRepo) FindBySlug(slug string) (*subDomain.Plan, error) {
	if p, ok := f.plans[slug]; ok {
		cp := *p
		return &cp, nil
	}
	return nil, nil
}

func (f *fakePlanCatalogRepo) Upsert(plan *subDomain.Plan) error {
	cp := *plan
	f.plans[plan.Slug] = &cp
	return nil
}

type fakeStripeCatalogGateway struct {
	calls map[string]int
}

func (f *fakeStripeCatalogGateway) EnsurePlanCatalog(slug, name, description string, monthly, yearly float64, currency string) (string, string, string, error) {
	if f.calls == nil {
		f.calls = map[string]int{}
	}
	f.calls[slug]++
	return "prod_" + slug, "price_" + slug + "_monthly", "price_" + slug + "_yearly", nil
}

func TestSyncStripeCatalog_FillsMissingIDsAndIsIdempotent(t *testing.T) {
	repo := &fakePlanCatalogRepo{
		plans: map[string]*subDomain.Plan{
			"pro": {
				Slug:         "pro",
				Name:         "Pro",
				Description:  "Pro",
				PriceMonthly: 29.90,
				PriceYearly:  299.00,
			},
			"premium": {
				Slug:         "premium",
				Name:         "Premium",
				Description:  "Premium",
				PriceMonthly: 49.90,
				PriceYearly:  499.00,
			},
		},
	}
	gw := &fakeStripeCatalogGateway{}

	err := syncStripeCatalog(repo, gw)
	require.NoError(t, err)

	err = syncStripeCatalog(repo, gw)
	require.NoError(t, err)

	require.Equal(t, 1, gw.calls["pro"])
	require.Equal(t, 1, gw.calls["premium"])

	pro := repo.plans["pro"]
	require.Equal(t, "prod_pro", pro.StripeProductID)
	require.Equal(t, "price_pro_monthly", pro.StripePriceIDMonthly)
	require.Equal(t, "price_pro_yearly", pro.StripePriceIDYearly)
}

func TestSyncStripeCatalog_SkipsPlanWhenAlreadyConfigured(t *testing.T) {
	repo := &fakePlanCatalogRepo{
		plans: map[string]*subDomain.Plan{
			"pro": {
				Slug:                 "pro",
				Name:                 "Pro",
				Description:          "Pro",
				PriceMonthly:         29.90,
				PriceYearly:          299.00,
				StripeProductID:      "prod_existing",
				StripePriceIDMonthly: "price_existing_monthly",
				StripePriceIDYearly:  "price_existing_yearly",
			},
			"premium": {
				Slug:         "premium",
				Name:         "Premium",
				Description:  "Premium",
				PriceMonthly: 49.90,
				PriceYearly:  499.00,
			},
		},
	}
	gw := &fakeStripeCatalogGateway{}

	err := syncStripeCatalog(repo, gw)
	require.NoError(t, err)

	require.Equal(t, 0, gw.calls["pro"])
	require.Equal(t, 1, gw.calls["premium"])
}

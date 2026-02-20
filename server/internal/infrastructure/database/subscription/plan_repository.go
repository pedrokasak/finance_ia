package subscription

import (
	"finance-ia/internal/domain/subscription"

	"gorm.io/gorm"
)

type PlanRepository struct {
	db *gorm.DB
}

func NewPlanRepository(db *gorm.DB) *PlanRepository {
	return &PlanRepository{db: db}
}

func (r *PlanRepository) FindAll() ([]*subscription.Plan, error) {
	var plans []*subscription.Plan
	if err := r.db.Where("is_active = ?", true).Order("price_monthly ASC").Find(&plans).Error; err != nil {
		return nil, err
	}
	// Populate features in-memory (avoids JSON column complexity)
	for _, p := range plans {
		p.Features = featuresForSlug(p.Slug)
	}
	return plans, nil
}

func (r *PlanRepository) FindBySlug(slug string) (*subscription.Plan, error) {
	var plan subscription.Plan
	if err := r.db.Where("slug = ?", slug).First(&plan).Error; err != nil {
		return nil, err
	}
	plan.Features = featuresForSlug(plan.Slug)
	return &plan, nil
}

func (r *PlanRepository) Upsert(plan *subscription.Plan) error {
	return r.db.Save(plan).Error
}

func featuresForSlug(slug string) []string {
	switch slug {
	case "pro":
		return []string{
			"Transações ilimitadas",
			"Categorização automática",
			"Gráficos avançados e interativos",
			"Metas financeiras personalizadas",
			"Relatórios detalhados",
			"Exportação em Excel/PDF",
			"Sincronização em nuvem",
			"Suporte prioritário",
		}
	case "premium":
		return []string{
			"Todos os recursos do Pro",
			"Análise de IA personalizada",
			"Projeções e previsões avançadas",
			"Alertas inteligentes",
			"Consultoria financeira automatizada",
			"Dashboard executivo",
			"API para integrações",
			"Suporte 24/7 via chat",
		}
	default: // free
		return []string{
			"Controle básico de receitas e despesas",
			"Categorização manual",
			"Gráficos simples",
			"Até 100 transações/mês",
			"Suporte por email",
		}
	}
}

var _ subscription.PlanRepository = (*PlanRepository)(nil)

package finance

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"finance-ia/internal/domain/finance"
)

type PostgresFinancialMethodRepository struct {
	db *gorm.DB
}

func NewPostgresFinancialMethodRepository(db *gorm.DB) finance.FinancialMethodRepository {
	return &PostgresFinancialMethodRepository{db: db}
}

func (r *PostgresFinancialMethodRepository) List() ([]*finance.FinancialMethod, error) {
	var methods []finance.FinancialMethod
	if err := r.db.Where("is_active = ?", true).Find(&methods).Error; err != nil {
		return nil, err
	}

	result := make([]*finance.FinancialMethod, len(methods))
	for i := range methods {
		result[i] = &methods[i]
	}
	return result, nil
}

func (r *PostgresFinancialMethodRepository) FindByID(id uuid.UUID) (*finance.FinancialMethod, error) {
	var method finance.FinancialMethod
	if err := r.db.First(&method, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil on not found based on architecture style, or could return err
		}
		return nil, err
	}
	return &method, nil
}

func (r *PostgresFinancialMethodRepository) FindByKey(key string) (*finance.FinancialMethod, error) {
	var method finance.FinancialMethod
	if err := r.db.First(&method, "key = ?", key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &method, nil
}

func (r *PostgresFinancialMethodRepository) Create(method *finance.FinancialMethod) error {
	return r.db.Create(method).Error
}

func (r *PostgresFinancialMethodRepository) SeedDefaults() error {
	methods := []finance.FinancialMethod{
		{
			Key:         "5-3-2",
			Name:        "Regra 5-3-2",
			Tagline:     "Foco em equilíbrio simples",
			Icon:        "CircleEllipsis",
			Color:       "text-lime-400",
			Bg:          "bg-lime-500/10 border-lime-500/30",
			Description: "50% para necessidades essenciais, 30% para estilo de vida e 20% para investimentos e reserva. Leitura rápida da 50-30-20 em formato compacto.",
			ForWho:      "Ideal para onboarding rápido com uma regra fácil de memorizar.",
			SplitRaw:    `[{"label":"Necessidades","percent":50,"color":"bg-lime-500"},{"label":"Estilo de Vida","percent":30,"color":"bg-teal-500"},{"label":"Investimentos","percent":20,"color":"bg-cyan-500"}]`,
		},
		{
			Key:         "60-10-10-10-10",
			Name:        "Regra 60-10-10-10-10",
			Tagline:     "Mais granular para metas simultâneas",
			Icon:        "ListTree",
			Color:       "text-indigo-400",
			Bg:          "bg-indigo-500/10 border-indigo-500/30",
			Description: "60% para custos fixos, 10% para investimentos de longo prazo, 10% para reserva, 10% para lazer e 10% para objetivos pessoais.",
			ForWho:      "Ideal para quem quer dividir melhor os 40% restantes em metas específicas.",
			SplitRaw:    `[{"label":"Custos Fixos","percent":60,"color":"bg-indigo-500"},{"label":"Invest. Longo Prazo","percent":10,"color":"bg-violet-500"},{"label":"Reserva","percent":10,"color":"bg-sky-500"},{"label":"Lazer","percent":10,"color":"bg-pink-500"},{"label":"Objetivos","percent":10,"color":"bg-amber-500"}]`,
		},
		{
			Key:         "50-30-20",
			Name:        "Regra 50-30-20",
			Tagline:     "O método mais popular do mundo",
			Icon:        "PieChart",
			Color:       "text-emerald-400",
			Bg:          "bg-emerald-500/10 border-emerald-500/30",
			Description: "50% da renda para necessidades básicas (moradia, alimentação, saúde), 30% para desejos (lazer, restaurantes) e 20% para poupança e investimentos.",
			ForWho:      "Ideal para quem está começando a organizar as finanças.",
			SplitRaw:    `[{"label":"Necessidades","percent":50,"color":"bg-emerald-500"},{"label":"Desejos","percent":30,"color":"bg-blue-500"},{"label":"Investimentos","percent":20,"color":"bg-purple-500"}]`,
		},
		{
			Key:         "60-20-20",
			Name:        "Regra 60-20-20",
			Tagline:     "Mais moderna e adaptável",
			Icon:        "Layers",
			Color:       "text-cyan-400",
			Bg:          "bg-cyan-500/10 border-cyan-500/30",
			Description: "60% para despesas fixas, 20% para investimentos e 20% para lazer e objetivos. Mais flexível que a 50-30-20 para famílias ou quem mora em cidade cara.",
			ForWho:      "Ideal para famílias, quem mora em cidade cara ou precisa de flexibilidade.",
			SplitRaw:    `[{"label":"Despesas Fixas","percent":60,"color":"bg-cyan-500"},{"label":"Investimentos","percent":20,"color":"bg-violet-500"},{"label":"Lazer e Metas","percent":20,"color":"bg-pink-500"}]`,
		},
		{
			Key:         "pay-yourself-first",
			Name:        "Pague-se Primeiro",
			Tagline:     "Invista antes de gastar",
			Icon:        "Coins",
			Color:       "text-yellow-400",
			Bg:          "bg-yellow-500/10 border-yellow-500/30",
			Description: "Assim que receber o salário, invista imediatamente antes de qualquer gasto. O percentual de investimento é definido por você — o foco está na ordem, não na proporção fixa.",
			ForWho:      "Excelente para construção de patrimônio e quem já tem alguma disciplina financeira.",
			SplitRaw:    `[{"label":"Investimento Imediato","percent":20,"color":"bg-yellow-500"},{"label":"Viva com o resto","percent":80,"color":"bg-orange-400"}]`,
		},
		{
			Key:         "emergency-reserve",
			Name:        "Reserva de Emergência",
			Tagline:     "Sua primeira meta financeira",
			Icon:        "Landmark",
			Color:       "text-sky-400",
			Bg:          "bg-sky-500/10 border-sky-500/30",
			Description: "Foco em construir uma reserva de 3 a 6 meses do custo de vida (ou 12 meses para autônomos). Guardar em Tesouro Selic ou CDB com liquidez diária.",
			ForWho:      "Essencial para quem não possui reserva de emergência ainda.",
			SplitRaw:    `[{"label":"Reserva (3-6 meses)","percent":30,"color":"bg-sky-500"},{"label":"Gastos Mensais","percent":70,"color":"bg-slate-400"}]`,
		},
		{
			Key:         "envelopes",
			Name:        "Método dos Envelopes",
			Tagline:     "Controle total por categoria",
			Icon:        "Wallet",
			Color:       "text-amber-400",
			Bg:          "bg-amber-500/10 border-amber-500/30",
			Description: "Cada categoria recebe um 'envelope' com valor fixo. Quando o envelope acaba, parou de gastar naquela categoria. Simples e visual.",
			ForWho:      "Ideal para quem gasta de forma impulsiva por categorias.",
			SplitRaw:    `[{"label":"Moradia","percent":30,"color":"bg-amber-500"},{"label":"Alimentação","percent":20,"color":"bg-orange-500"},{"label":"Transporte","percent":15,"color":"bg-yellow-500"},{"label":"Outros","percent":35,"color":"bg-red-500"}]`,
		},
		{
			Key:         "zero-based",
			Name:        "Orçamento Base Zero",
			Tagline:     "Cada real tem uma destinação",
			Icon:        "Target",
			Color:       "text-blue-400",
			Bg:          "bg-blue-500/10 border-blue-500/30",
			Description: "Renda - Despesas = 0. Todo real é alocado intencionalmente. Máximo controle, pois nenhum dinheiro 'some' sem destino definido.",
			ForWho:      "Ideal para quem quer controle total e obsessivo das finanças.",
			SplitRaw:    `[{"label":"Fixos","percent":40,"color":"bg-blue-500"},{"label":"Variáveis","percent":35,"color":"bg-indigo-500"},{"label":"Reservas","percent":25,"color":"bg-violet-500"}]`,
		},
		{
			Key:         "goal-based",
			Name:        "Planejamento por Objetivos",
			Tagline:     "Separe por metas: curto, médio e longo prazo",
			Icon:        "TrendingUp",
			Color:       "text-green-400",
			Bg:          "bg-green-500/10 border-green-500/30",
			Description: "Cada objetivo tem seu horizonte temporal e tipo de investimento. Viagem → renda fixa. Carro → renda fixa moderada. Aposentadoria → renda variável.",
			ForWho:      "Ideal para quem tem múltiplos objetivos financeiros simultâneos.",
			SplitRaw:    `[{"label":"Curto Prazo (até 2a)","percent":20,"color":"bg-green-500"},{"label":"Médio Prazo (2-5a)","percent":25,"color":"bg-teal-500"},{"label":"Longo Prazo (5a+)","percent":20,"color":"bg-emerald-700"},{"label":"Gastos","percent":35,"color":"bg-gray-400"}]`,
		},
		{
			Key:         "savings-rate",
			Name:        "Taxa de Poupança",
			Tagline:     "Maximize o quanto você poupa",
			Icon:        "Percent",
			Color:       "text-rose-400",
			Bg:          "bg-rose-500/10 border-rose-500/30",
			Description: "Taxa de poupança = valor investido ÷ renda. 10% é básico, 20% é bom, 30%+ acelera a riqueza. Referências: 10% (básico), 20% (bom), 50%+ (modo monge financeiro).",
			ForWho:      "Para quem quer acelerar a construção de patrimônio.",
			SplitRaw:    `[{"label":"Investimentos (30%+)","percent":30,"color":"bg-rose-500"},{"label":"Gastos (70% ou menos)","percent":70,"color":"bg-pink-300"}]`,
		},
		{
			Key:         "70-20-10",
			Name:        "Regra 70-20-10",
			Tagline:     "Gastos, poupança e dívidas",
			Icon:        "Flame",
			Color:       "text-red-400",
			Bg:          "bg-red-500/10 border-red-500/30",
			Description: "70% para gastos mensais (necessários e desejos), 20% para poupança e investimentos, 10% para quitação de dívidas ou doações.",
			ForWho:      "Ideal para quem tem dívidas e quer estruturar a saída delas.",
			SplitRaw:    `[{"label":"Gastos","percent":70,"color":"bg-red-500"},{"label":"Poupança","percent":20,"color":"bg-pink-500"},{"label":"Dívidas/Doação","percent":10,"color":"bg-fuchsia-500"}]`,
		},
	}

	for _, m := range methods {
		var existing finance.FinancialMethod
		err := r.db.Where("key = ?", m.Key).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			if err := r.db.Create(&m).Error; err != nil {
				return err
			}
		} else if err == nil {
			// Update the strings if the structure changed
			r.db.Model(&existing).Updates(map[string]interface{}{
				"name":        m.Name,
				"tagline":     m.Tagline,
				"description": m.Description,
				"for_who":     m.ForWho,
				"split_raw":   m.SplitRaw,
				"icon":        m.Icon,
				"color":       m.Color,
				"bg":          m.Bg,
			})
		}
	}
	return nil
}

package main

import (
	authapp "finance-ia/internal/application/usecase/auth"
	userapp "finance-ia/internal/application/usecase/user"
	"finance-ia/internal/config"
	"finance-ia/internal/config/database"
	"finance-ia/internal/domain/ai"
	"finance-ia/internal/domain/auth"
	"finance-ia/internal/domain/email"
	"finance-ia/internal/domain/finance"
	subDomain "finance-ia/internal/domain/subscription"
	"finance-ia/internal/domain/user"
	infraAI "finance-ia/internal/infrastructure/ai"
	infraAIRepo "finance-ia/internal/infrastructure/database/ai"
	infraauth "finance-ia/internal/infrastructure/database/auth"
	infraEmailRepo "finance-ia/internal/infrastructure/database/email"
	infraFinance "finance-ia/internal/infrastructure/database/finance"
	infraSubRepo "finance-ia/internal/infrastructure/database/subscription"
	infrauser "finance-ia/internal/infrastructure/database/user"
	stripeAdapter "finance-ia/internal/infrastructure/payment/stripe"
	aiHandlers "finance-ia/internal/interfaces/handlers/ai"
	authHandlers "finance-ia/internal/interfaces/handlers/auth"
	financeHandlers "finance-ia/internal/interfaces/handlers/finance"
	subHandlers "finance-ia/internal/interfaces/handlers/subscription"
	useHandlers "finance-ia/internal/interfaces/handlers/user"
	"finance-ia/internal/router"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado")
	}

	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Falha ao conectar com o banco:", err)
	}

	// Run migrations (creates all tables)
	if err := database.Migrate(db); err != nil {
		log.Fatal("Falha ao executar migrations:", err)
	}

	// ────────────────────────────────────
	// Auth domain
	// ────────────────────────────────────
	userRepo := infrauser.NewUserRepository(db)
	authRepo := infraauth.NewAuthRepository(db)
	tokenRepo := infraEmailRepo.NewTokenRepository(db)
	_ = db.AutoMigrate(&infraEmailRepo.PasswordResetToken{}, &finance.FinancialMethod{})

	emailService := email.NewSMTPService()
	userService := user.NewService(userRepo)
	authService := auth.NewService(authRepo, tokenRepo, emailService)

	userUsecase := userapp.NewUseCase(userService)
	authUsecase := authapp.NewUseCase(authService)

	userHandler := useHandlers.NewUserHandler(userUsecase)
	authHandler := authHandlers.NewAuthHandler(authUsecase)

	// ────────────────────────────────────
	// Finance domain
	// ────────────────────────────────────
	txRepo := infraFinance.NewTransactionRepository(db)
	categoryRepo := infraFinance.NewCategoryRepository(db)
	budgetRepo := infraFinance.NewBudgetRepository(db)
	methodRepo := infraFinance.NewPostgresFinancialMethodRepository(db)

	// Seed default categories on startup
	if err := categoryRepo.SeedDefaults(); err != nil {
		log.Printf("Warning: failed to seed default categories: %v", err)
	}

	// Seed financial methods
	if err := methodRepo.(*infraFinance.PostgresFinancialMethodRepository).SeedDefaults(); err != nil {
		log.Printf("Warning: failed to seed financial methods: %v", err)
	}

	financeService := finance.NewService(txRepo, categoryRepo, budgetRepo, methodRepo)
	financeHandler := financeHandlers.NewFinanceHandler(financeService)

	// ────────────────────────────────────
	// AI domain
	// ────────────────────────────────────
	insightRepo := infraAIRepo.NewInsightRepository(db)

	var aiProvider ai.AIProvider
	geminiProvider, err := infraAI.NewGeminiProvider()
	if err != nil {
		log.Printf("Warning: Gemini AI unavailable (%v) — AI features will return error", err)
		aiProvider = nil
	} else {
		aiProvider = geminiProvider
	}

	var aiService *ai.Service
	if aiProvider != nil {
		aiService = ai.NewService(insightRepo, aiProvider)
	}

	var aiHandler *aiHandlers.AIHandler
	if aiService != nil {
		aiHandler = aiHandlers.NewAIHandler(aiService, financeService)
	}

	// ────────────────────────────────────
	// Subscription / Payment domain
	// ────────────────────────────────────
	subRepo := infraSubRepo.NewSubscriptionRepository(db)
	planRepo := infraSubRepo.NewPlanRepository(db)

	// AutoMigrate Plan table
	if err := db.AutoMigrate(&subDomain.Plan{}); err != nil {
		log.Printf("Warning: failed to migrate Plan table: %v", err)
	}

	// Seed subscription plans (links to Stripe Price IDs from env vars)
	if err := seedPlans(planRepo); err != nil {
		log.Printf("Warning: failed to seed plans: %v", err)
	}

	gateway := stripeAdapter.NewAdapter()
	subscriptionHandler := subHandlers.NewSubscriptionHandler(gateway, subRepo, planRepo, userRepo, db)
	planHandler := subHandlers.NewPlanHandler(gateway, planRepo)

	// ────────────────────────────────────
	// Subscription service (for future use cases)
	// ────────────────────────────────────
	_ = subDomain.GetPlanFeatures("free") // ensure package is used

	// ────────────────────────────────────
	// Router
	// ────────────────────────────────────
	registrars := []router.RouteRegistrar{
		userHandler,
		authHandler,
		financeHandler,
		subscriptionHandler,
		planHandler,
	}

	if aiHandler != nil {
		registrars = append(registrars, aiHandler)
	}

	r := router.NewRouter(registrars...)

	log.Printf("Servidor rodando na porta %s", cfg.AppPort)
	log.Fatal(r.Run(":" + cfg.AppPort))
}

// seedPlans creates or updates subscription plans in the DB using Stripe Price IDs from env vars.
// Plans are seeded idempotently — safe to run on every startup.
func seedPlans(repo *infraSubRepo.PlanRepository) error {
	plans := []subDomain.Plan{
		{
			Slug:                 "free",
			Name:                 "Básico",
			Description:          "Perfeito para começar",
			PriceMonthly:         0,
			PriceYearly:          0,
			StripePriceIDMonthly: "",
			StripePriceIDYearly:  "",
			MaxTransactions:      100,
			AIInsights:           false,
			ExportData:           false,
			IsActive:             true,
			Features: []subDomain.PlanFeature{
				{Description: "Controle básico de receitas e despesas"},
				{Description: "Categorização manual"},
				{Description: "Gráficos simples"},
				{Description: "Até 100 transações/mês"},
				{Description: "Suporte por email"},
			},
		},
		{
			Slug:                 "pro",
			Name:                 "Pro",
			Description:          "Ideal para controle avançado",
			PriceMonthly:         29.90,
			PriceYearly:          299.00,
			StripePriceIDMonthly: getEnv("STRIPE_PRICE_PRO"),
			StripePriceIDYearly:  getEnv("STRIPE_PRICE_PRO_YEARLY"),
			MaxTransactions:      -1,
			AIInsights:           true,
			ExportData:           true,
			IsActive:             true,
			Features: []subDomain.PlanFeature{
				{Description: "Transações ilimitadas"},
				{Description: "Categorização automática"},
				{Description: "Gráficos avançados e interativos"},
				{Description: "Metas financeiras personalizadas"},
				{Description: "Relatórios detalhados"},
				{Description: "Exportação em Excel/PDF"},
				{Description: "Sincronização em nuvem"},
				{Description: "Suporte prioritário"},
			},
		},
		{
			Slug:                 "premium",
			Name:                 "Premium",
			Description:          "Máximo controle financeiro com IA",
			PriceMonthly:         49.90,
			PriceYearly:          499.00,
			StripePriceIDMonthly: getEnv("STRIPE_PRICE_PREMIUM"),
			StripePriceIDYearly:  getEnv("STRIPE_PRICE_PREMIUM_YEARLY"),
			MaxTransactions:      -1,
			AIInsights:           true,
			ExportData:           true,
			IsActive:             true,
			Features: []subDomain.PlanFeature{
				{Description: "Todos os recursos do Pro"},
				{Description: "Análise de IA personalizada"},
				{Description: "Projeções e previsões avançadas"},
				{Description: "Alertas inteligentes"},
				{Description: "Consultoria financeira automatizada"},
				{Description: "Dashboard executivo"},
				{Description: "API para integrações"},
				{Description: "Suporte 24/7 via chat"},
			},
		},
	}
	for i := range plans {
		// Only seed if plan doesn't exist or we want to overwrite.
		// For idempotency and dynamic features, we'll check if exists first.
		existing, err := repo.FindBySlug(plans[i].Slug)
		if err != nil || existing == nil {
			if err := repo.Upsert(&plans[i]); err != nil {
				return err
			}
		} else if len(existing.Features) == 0 {
			// If plan exists but has no features, add default ones
			existing.Features = plans[i].Features
			if err := repo.Upsert(existing); err != nil {
				return err
			}
		}
	}
	return nil
}

func getEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return ""
}

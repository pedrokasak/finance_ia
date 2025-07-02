package main

import (
	userapp "finance-ia/internal/application/user"
	"finance-ia/internal/config"
	"finance-ia/internal/config/database"
	"finance-ia/internal/config/middleware"
	"finance-ia/internal/domain/user"
	infrauser "finance-ia/internal/infrastructure/database/user"
	"finance-ia/internal/interfaces/handlers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado")
	}

	// Inicializar configuração
	cfg := config.Load()
	
	// Conectar ao banco de dados
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Falha ao conectar com o banco:", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal("Falha ao executar migrations:", err)
	}

		// Repositório do usuário
    userRepo := infrauser.NewUserRepository(db)
		
    // Service do domínio
    userService := user.NewService(userRepo);

		userUsecase := userapp.NewUseCase(userService)
		
    // Handler
		userHandler := handlers.NewUserHandler(userUsecase)

    // Configurar router
    router := setupRouter(userHandler)

	// Iniciar servidor
	log.Printf("Servidor rodando na porta %s", cfg.AppPort)
	log.Fatal(router.Run(":" + cfg.AppPort))

}

func setupRouter(uh *handlers.UserHandler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// Rotas públicas
	public := router.Group("/api/v1")
	{
		public.POST("/register", uh.Register)
		public.POST("/login", uh.Login)
	}

	// Rotas protegidas
	protected := router.Group("/api/v1")
	protected.Use(middleware.JWTAuth())
	{
		// Usuários
		// protected.GET("/profile", uh.GetProfile)
		// protected.PUT("/profile", uh.UpdateProfile)

		// Transações
		// protected.POST("/transactions", th.Create)
		// protected.GET("/transactions", th.List)
		// protected.GET("/transactions/:id", th.GetByID)
		// protected.PUT("/transactions/:id", th.Update)
		// protected.DELETE("/transactions/:id", th.Delete)

		// Categorias
		// protected.POST("/categories", ch.Create)
		// protected.GET("/categories", ch.List)
		// protected.PUT("/categories/:id", ch.Update)
		// protected.DELETE("/categories/:id", ch.Delete)

		// Relatórios
		// protected.GET("/reports/monthly", th.MonthlyReport)
		// protected.GET("/reports/dashboard", th.DashboardData)
	}

	return router
}
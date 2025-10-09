package main

import (
	userapp "finance-ia/internal/application/usecase/user"
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
		public.POST("/auth/signup", uh.Register)
		public.POST("/auth/login", uh.Login)
		// public.POST("/auth/refresh", uh.Refresh)
		public.POST("/auth/forgot-password", uh.ForgotPassword)
		public.POST("/auth/reset-password", uh.ResetPassword)
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	// Rotas protegidas
	protected := router.Group("/api/v1")
	protected.Use(middleware.JWTAuth())
	{
		// Usuários
		protected.GET("/users", uh.GetAllUsers)
		protected.GET("/user/:id", uh.GetUserByID)
		protected.PUT("/user/:id", uh.UpdateUser)
		protected.DELETE("/user/:id", uh.DeleteUser)


		// Relatórios
		// protected.GET("/reports/monthly", th.MonthlyReport)
		// protected.GET("/reports/dashboard", th.DashboardData)
	}

	return router
}
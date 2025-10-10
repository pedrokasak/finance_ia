package main

import (
	authapp "finance-ia/internal/application/usecase/auth"
	userapp "finance-ia/internal/application/usecase/user"
	"finance-ia/internal/config"
	"finance-ia/internal/config/database"
	"finance-ia/internal/config/middleware"
	"finance-ia/internal/domain/auth"
	"finance-ia/internal/domain/user"
	infraauth "finance-ia/internal/infrastructure/database/auth"
	infrauser "finance-ia/internal/infrastructure/database/user"
	authHandlers "finance-ia/internal/interfaces/handlers/auth"
	useHandlers "finance-ia/internal/interfaces/handlers/user"
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
		authRepo := infraauth.NewAuthRepository(db)
		
    // Service do domínio
    userService := user.NewService(userRepo)
		authService := auth.NewService(authRepo)

		// Use case da aplicação

		userUsecase := userapp.NewUseCase(userService)
		authUsecase := authapp.NewUseCase(authService)

    // Handler
		userHandler := useHandlers.NewUserHandler(userUsecase)
		authHandler := authHandlers.NewAuthHandler(authUsecase)

    // Configurar router
    router := setupRouter(userHandler, authHandler)

	// Iniciar servidor
	log.Printf("Servidor rodando na porta %s", cfg.AppPort)
	log.Fatal(router.Run(":" + cfg.AppPort))

}


func setupRouter(uh *useHandlers.UserHandler, ah *authHandlers.AuthHandler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// Rotas públicas
	public := router.Group("/api/v1")
	{
		public.POST("/auth/signup", uh.Register)
		public.POST("/auth/login", ah.Login)
		// public.POST("/auth/refresh", uh.Refresh)
		public.POST("/auth/forgot-password", ah.ForgotPassword)
		public.POST("/auth/reset-password", ah.ResetPassword)
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
package main

import (
	authapp "finance-ia/internal/application/usecase/auth"
	userapp "finance-ia/internal/application/usecase/user"
	"finance-ia/internal/config"
	"finance-ia/internal/config/database"
	"finance-ia/internal/domain/auth"
	"finance-ia/internal/domain/email"
	"finance-ia/internal/domain/user"
	infraauth "finance-ia/internal/infrastructure/database/auth"
	emailRepository "finance-ia/internal/infrastructure/database/email"
	infrauser "finance-ia/internal/infrastructure/database/user"
	authHandlers "finance-ia/internal/interfaces/handlers/auth"
	useHandlers "finance-ia/internal/interfaces/handlers/user"
	"finance-ia/internal/router"
	"log"

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
	
	// Migrations
	db.AutoMigrate(&emailRepository.PasswordResetToken{})

	if err != nil {
		log.Fatal("Falha ao conectar com o banco:", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal("Falha ao executar migrations:", err)
	}

		// Repositório do usuário
    userRepo := infrauser.NewUserRepository(db)
		authRepo := infraauth.NewAuthRepository(db)
		tokenRepo := emailRepository.NewTokenRepository(db)

		// Inicializa serviço de email
		emailService := email.NewSMTPService()
		
    // Service do domínio
    userService := user.NewService(userRepo)
		authService := auth.NewService(authRepo, tokenRepo, emailService)

		// Use case da aplicação

		userUsecase := userapp.NewUseCase(userService)
		authUsecase := authapp.NewUseCase(authService)

    // Handler
		userHandler := useHandlers.NewUserHandler(userUsecase)
		authHandler := authHandlers.NewAuthHandler(authUsecase)

    // Configurar router
    router := router.NewRouter(userHandler, authHandler)

	// Iniciar servidor
	log.Printf("Servidor rodando na porta %s", cfg.AppPort)
	log.Fatal(router.Run(":" + cfg.AppPort))

}

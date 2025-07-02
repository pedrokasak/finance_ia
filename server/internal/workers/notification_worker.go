package workers

import (
	"context"
	"finance-ia/internal/domain/user"

	"log"
	"time"

	"gorm.io/gorm"
)

type NotificationWorker struct {
	db *gorm.DB
}

func NewNotificationWorker(db *gorm.DB) *NotificationWorker {
	return &NotificationWorker{db: db}
}

// Goroutine para processar notificações de gastos mensais
func (w *NotificationWorker) StartMonthlySpendingNotifications(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour) // Executa diariamente
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Parando worker de notificações...")
			return
		case <-ticker.C:
			w.checkMonthlySpendingLimits()
		}
	}
}

func (w *NotificationWorker) checkMonthlySpendingLimits() {
	log.Println("Verificando limites de gastos mensais...")

	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)

	// Buscar todos os usuários
	var users []user.User
	if err := w.db.Find(&users).Error; err != nil {
		log.Printf("Erro ao buscar usuários: %v", err)
		return
	}

	for _, user := range users {
		var monthlyExpense float64
		// Exemplo: notificar se gastou mais de R$ 5000 no mês
		if monthlyExpense > 5000 {
			log.Printf("Usuário %s atingiu limite de gastos: R$ %.2f", user.Email, monthlyExpense)
			// Aqui você pode implementar envio de email, push notification, etc.
		}
	}
}
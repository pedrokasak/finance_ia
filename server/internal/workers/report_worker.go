package workers

import (
	"context"
	"finance-ia/internal/domain/user"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportWorker struct {
	db *gorm.DB
}

func NewReportWorker(db *gorm.DB) *ReportWorker {
	return &ReportWorker{db: db}
}

// Goroutine para gerar relatórios automáticos
func (w *ReportWorker) StartMonthlyReportGeneration(ctx context.Context) {
	// Executa no primeiro dia de cada mês
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Parando worker de relatórios...")
			return
		case <-ticker.C:
			if time.Now().Day() == 1 { // Primeiro dia do mês
				w.generateMonthlyReportsForAllUsers()
			}
		}
	}
}

func (w *ReportWorker) generateMonthlyReportsForAllUsers() {
	log.Println("Gerando relatórios mensais para todos os usuários...")

	lastMonth := time.Now().AddDate(0, -1, 0)
	
	var users []user.User
	if err := w.db.Find(&users).Error; err != nil {
		log.Printf("Erro ao buscar usuários: %v", err)
		return
	}

	for _, user := range users {
		go w.generateUserMonthlyReport(user.ID, lastMonth.Year(), lastMonth.Month())
	}
}

func (w *ReportWorker) generateUserMonthlyReport(userID uuid.UUID, year int, month time.Month) {
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	var totalIncome, totalExpense float64

	// // Calcular receitas
	// w.db.Model(&transaction.Transaction{}).
	// 	Where("user_id = ? AND type = ? AND date >= ? AND date <= ?", userID, "income", startDate, endDate).
	// 	Select("COALESCE(SUM(amount), 0)").
	// 	Scan(&totalIncome)

	// // Calcular despesas
	// w.db.Model(&transaction.Transaction{}).
	// 	Where("user_id = ? AND type = ? AND date >= ? AND date <= ?", userID, "expense", startDate, endDate).
	// 	Select("COALESCE(SUM(amount), 0)").
	// 	Scan(&totalExpense)

	log.Printf("Relatório gerado para usuário %s - Receitas: R$ %.2f, Despesas: R$ %.2f", 
		userID, totalIncome, totalExpense)
	
	// Aqui você pode salvar o relatório em uma tabela específica ou enviar por email
}
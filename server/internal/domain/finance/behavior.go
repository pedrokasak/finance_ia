package finance

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

// AnalyzeBehavior performs advanced behavioral analysis on the user's transactions
// from the last 3 months to detect spending patterns.
func (s *Service) AnalyzeBehavior(userID uuid.UUID) (*BehavioralAnalysis, error) {
	now := time.Now()
	threeMonthsAgo := now.AddDate(0, -3, 0)

	txs, err := s.txRepo.FindByPeriod(userID, threeMonthsAgo, now)
	if err != nil {
		return nil, fmt.Errorf("analyze behavior: fetch transactions: %w", err)
	}

	if len(txs) == 0 {
		return &BehavioralAnalysis{RiskLevel: "low", RiskMessage: "Adicione mais transações para análise completa."}, nil
	}

	analysis := &BehavioralAnalysis{}

	// 1. Weekend vs Weekday ratio & emotional spending
	weekdayTotals := make([]float64, 7) // 0=Sun, 6=Sat
	weekdayCount := make([]int, 7)
	var weekendTotal, weekdayTotal float64
	var weekendCount, weekdayCount2 int

	// 2. Salary effect (spike after income days)
	var incomeDates []time.Time

	// 3. Impulse purchases
	var impCount int
	var impTotal float64

	// 4. Category monthly totals (for money leaks)
	type monthKey struct{ year, month int }
	catMonthly := make(map[string]map[monthKey]float64) // catName -> month -> total
	catColors := make(map[string]string)

	for _, tx := range txs {
		if tx.Type == TransactionTypeIncome {
			incomeDates = append(incomeDates, tx.Date)
			continue
		}
		// expense
		wd := int(tx.Date.Weekday())
		weekdayTotals[wd] += tx.Amount
		weekdayCount[wd]++

		isWeekend := wd == 0 || wd == 6
		if isWeekend {
			weekendTotal += tx.Amount
			weekendCount++
		} else {
			weekdayTotal += tx.Amount
			weekdayCount2++
		}

		// Impulse: < R$50, not recurring
		if tx.Amount < 50 && !tx.IsRecurring {
			impCount++
			impTotal += tx.Amount
		}

		// Category monthly totals
		if tx.Category != nil {
			key := monthKey{tx.Date.Year(), int(tx.Date.Month())}
			if catMonthly[tx.Category.Name] == nil {
				catMonthly[tx.Category.Name] = make(map[monthKey]float64)
				catColors[tx.Category.Name] = tx.Category.Color
			}
			catMonthly[tx.Category.Name][key] += tx.Amount
		}
	}

	analysis.ImpulsivePurchaseCount = impCount
	analysis.ImpulsivePurchaseTotal = impTotal

	// Weekend vs weekday ratio
	avgWeekend := 0.0
	avgWeekday := 0.0
	if weekendCount > 0 {
		avgWeekend = weekendTotal / float64(weekendCount)
	}
	if weekdayCount2 > 0 {
		avgWeekday = weekdayTotal / float64(weekdayCount2)
	}
	if avgWeekday > 0 {
		analysis.WeekendVsWeekdayRatio = roundTo2(avgWeekend / avgWeekday)
	}

	// Emotional spending score: % of expenses on weekends
	totalExp := weekendTotal + weekdayTotal
	if totalExp > 0 {
		analysis.EmotionalSpendingScore = roundTo2((weekendTotal / totalExp) * 100)
	}

	// Top weekday by expense
	topWD, topAmount := 0, 0.0
	weekdayNames := []string{"Domingo", "Segunda", "Terça", "Quarta", "Quinta", "Sexta", "Sábado"}
	for i, total := range weekdayTotals {
		avg := 0.0
		if weekdayCount[i] > 0 {
			avg = total / float64(weekdayCount[i])
		}
		if avg > topAmount {
			topAmount = avg
			topWD = i
		}
	}
	analysis.TopWeekdayByExpense = topWD
	analysis.TopWeekdayName = weekdayNames[topWD]

	// Salary effect: check spending spikes in 7 days after income
	if len(incomeDates) > 0 {
		var afterSalarySpend, normalSpend float64
		var afterDays, normalDays int

		for _, tx := range txs {
			if tx.Type == TransactionTypeIncome {
				continue
			}
			isAfterSalary := false
			for _, incDate := range incomeDates {
				diff := tx.Date.Sub(incDate).Hours() / 24
				if diff >= 0 && diff <= 7 {
					isAfterSalary = true
					break
				}
			}
			if isAfterSalary {
				afterSalarySpend += tx.Amount
				afterDays++
			} else {
				normalSpend += tx.Amount
				normalDays++
			}
		}

		avgAfter := 0.0
		avgNormal := 0.0
		if afterDays > 0 {
			avgAfter = afterSalarySpend / float64(afterDays)
		}
		if normalDays > 0 {
			avgNormal = normalSpend / float64(normalDays)
		}

		if avgAfter > avgNormal*1.3 { // 30% spike = salary effect detected
			analysis.SalaryEffectDays = 7
			analysis.SalaryEffectAmount = roundTo2(avgAfter - avgNormal)
		}
	}

	// Money leak categories: month-over-month growth > 10%
	currentMonth := monthKey{now.Year(), int(now.Month())}
	previousMonth := monthKey{now.AddDate(0, -1, 0).Year(), int(now.AddDate(0, -1, 0).Month())}

	for catName, monthly := range catMonthly {
		curr := monthly[currentMonth]
		prev := monthly[previousMonth]
		if prev > 0 && curr > 0 {
			growth := ((curr - prev) / prev) * 100
			if growth > 10 {
				analysis.MoneyLeakCategories = append(analysis.MoneyLeakCategories, MoneyLeakCategory{
					CategoryName:  catName,
					CategoryColor: catColors[catName],
					GrowthPercent: roundTo2(growth),
					CurrentMonth:  roundTo2(curr),
					PreviousMonth: roundTo2(prev),
				})
			}
		}
	}

	// Risk level calculation
	riskScore := 0
	if analysis.WeekendVsWeekdayRatio > 1.5 {
		riskScore++
	}
	if analysis.EmotionalSpendingScore > 40 {
		riskScore++
	}
	if analysis.SalaryEffectAmount > 0 {
		riskScore++
	}
	if len(analysis.MoneyLeakCategories) > 0 {
		riskScore++
	}
	if analysis.ImpulsivePurchaseCount > 10 {
		riskScore++
	}

	switch {
	case riskScore >= 4:
		analysis.RiskLevel = "high"
		analysis.RiskMessage = "Padrão de gasto de alto risco detectado. Considere revisar seus hábitos urgentemente."
	case riskScore >= 2:
		analysis.RiskLevel = "medium"
		analysis.RiskMessage = "Alguns padrões preocupantes identificados. Pequenos ajustes podem melhorar muito sua saúde financeira."
	default:
		analysis.RiskLevel = "low"
		analysis.RiskMessage = "Seus hábitos financeiros estão no caminho certo! Continue assim."
	}

	return analysis, nil
}

func roundTo2(v float64) float64 {
	return math.Round(v*100) / 100
}

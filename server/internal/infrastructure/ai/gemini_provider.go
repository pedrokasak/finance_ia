package ai

import (
	"context"
	"finance-ia/internal/domain/ai"
	"fmt"
	"os"
	"strings"
	"time"

	"google.golang.org/genai"
)

// GeminiProvider implements ai.AIProvider using Google Gemini
type GeminiProvider struct {
	client    *genai.Client
	modelName string
}

// NewGeminiProvider creates a new Google Gemini AI provider
func NewGeminiProvider() (*GeminiProvider, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is required")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiProvider{
		client:    client,
		modelName: "gemini-2.0-flash",
	}, nil
}

// GenerateInsight generates a short, direct financial insight (free plan)
func (g *GeminiProvider) GenerateInsight(ctx ai.FinancialContext) (*ai.AIInsight, error) {
	prompt := buildFreeInsightPrompt(ctx)

	response, err := g.generate(prompt)
	if err != nil {
		return nil, err
	}

	return &ai.AIInsight{
		Type:    ai.InsightTypeTip,
		Title:   "💡 Dica Financeira",
		Content: response,
	}, nil
}

// GenerateDiagnostic generates a comprehensive financial analysis
func (g *GeminiProvider) GenerateDiagnostic(ctx ai.FinancialContext) (*ai.AIInsight, error) {
	prompt := buildHealthDiagnosticPrompt(ctx)
	response, err := g.generate(prompt)
	if err != nil {
		return nil, err
	}
	return &ai.AIInsight{
		Type:    ai.InsightTypeWarning,
		Title:   "⚠️ Diagnóstico de Saúde Financeira",
		Content: response,
	}, nil
}

// GenerateProjection generates a future projection
func (g *GeminiProvider) GenerateProjection(ctx ai.FinancialContext) (*ai.AIInsight, error) {
	prompt := buildProjectionPrompt(ctx)
	response, err := g.generate(prompt)
	if err != nil {
		return nil, err
	}
	return &ai.AIInsight{
		Type:    ai.InsightTypeProjection,
		Title:   "📈 Projeção Financeira",
		Content: response,
	}, nil
}

func (g *GeminiProvider) generate(prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := g.client.Models.GenerateContent(
		ctx,
		g.modelName,
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("gemini: generate content: %w", err)
	}

	if len(result.Candidates) == 0 {
		return "", fmt.Errorf("gemini: no candidates returned")
	}

	var parts []string
	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			parts = append(parts, part.Text)
		}
	}

	return strings.Join(parts, ""), nil
}

func buildFreeInsightPrompt(ctx ai.FinancialContext) string {
	return fmt.Sprintf(`Você é um consultor financeiro pessoal brasileiro. Analise os dados financeiros do usuário e forneça UMA dica curta e direta (máximo 2 frases).

Dados do usuário %s:
- Renda mensal: R$ %.2f
- Gastos totais: R$ %.2f
- Saldo: R$ %.2f
- Taxa de poupança: %.1f%%
- Score de saúde financeira: %d/1000 (%s)

Principais categorias de gasto:
%s

A dica deve ser:
- Curta e direta (máximo 2 frases)
- Baseada nos dados reais
- Em português brasileiro
- Sem jargões financeiros complexos

Responda apenas com a dica, sem introdução.`,
		ctx.UserName,
		ctx.TotalIncome,
		ctx.TotalExpenses,
		ctx.Balance,
		ctx.SavingsRate,
		ctx.HealthScore,
		ctx.HealthLevel,
		formatCategories(ctx.TopCategories),
	)
}

func buildHealthDiagnosticPrompt(ctx ai.FinancialContext) string {
	return fmt.Sprintf(`Você é um consultor financeiro pessoal brasileiro especialista. Faça um diagnóstico completo da saúde financeira do usuário.

Dados do usuário %s em %s:
- Renda: R$ %.2f | Gastos: R$ %.2f | Saldo: R$ %.2f
- Taxa de poupança: %.1f%% | Score: %d/1000 (%s)

Principais gastos:
%s

Tendência mensal:
%s

Forneça um diagnóstico em formato de lista com:
1. Status atual da saúde financeira
2. Principais problemas identificados
3. Pontos positivos

Use linguagem clara e empática. Máximo 200 palavras. Em português brasileiro.`,
		ctx.UserName, ctx.Period,
		ctx.TotalIncome, ctx.TotalExpenses, ctx.Balance,
		ctx.SavingsRate, ctx.HealthScore, ctx.HealthLevel,
		formatCategories(ctx.TopCategories),
		formatTrends(ctx.MonthlyTrends),
	)
}

func buildRecommendationsPrompt(ctx ai.FinancialContext) string {
	return fmt.Sprintf(`Como consultor financeiro, forneça 3 recomendações priorizadas para o usuário %s melhorar sua saúde financeira.

Contexto:
- Renda: R$ %.2f | Gastos: R$ %.2f | Saldo: R$ %.2f
- Taxa de poupança atual: %.1f%% | Meta ideal: 20%%
- Score: %d/1000

Categorias problemáticas:
%s

Para cada recomendação, inclua:
- O que fazer (ação específica)
- Impacto estimado em R$ ou %%

Formato: lista numerada, linguagem direta. Máximo 150 palavras. Em português brasileiro.`,
		ctx.UserName,
		ctx.TotalIncome, ctx.TotalExpenses, ctx.Balance,
		ctx.SavingsRate,
		ctx.HealthScore,
		formatCategories(ctx.TopCategories),
	)
}

func buildProjectionPrompt(ctx ai.FinancialContext) string {
	return fmt.Sprintf(`Como analista financeiro, faça uma projeção do futuro financeiro do usuário %s.

Dados atuais:
- Renda: R$ %.2f | Gastos: R$ %.2f | Poupança mensal: R$ %.2f
- Taxa de poupança: %.1f%%

Considerando o comportamento atual, projete:
1. Saldo em 3 meses
2. Saldo em 6 meses
3. O que acontece se reduzir gastos em 15%%
4. Tempo estimado para sair do ciclo atual

Seja específico com valores em R$. Máximo 150 palavras. Em português brasileiro.`,
		ctx.UserName,
		ctx.TotalIncome, ctx.TotalExpenses, ctx.Balance,
		ctx.SavingsRate,
	)
}

func formatCategories(cats []ai.CategorySpend) string {
	if len(cats) == 0 {
		return "Sem dados de categorias"
	}
	var sb strings.Builder
	for _, c := range cats {
		sb.WriteString(fmt.Sprintf("- %s: R$ %.2f (%.1f%%)\n", c.Name, c.Amount, c.Percentage))
	}
	return sb.String()
}

func formatTrends(trends []ai.MonthTrend) string {
	if len(trends) == 0 {
		return "Sem dados históricos"
	}
	var sb strings.Builder
	for _, t := range trends {
		sb.WriteString(fmt.Sprintf("- %s: Renda R$ %.2f / Gastos R$ %.2f\n", t.Month, t.Income, t.Expenses))
	}
	return sb.String()
}

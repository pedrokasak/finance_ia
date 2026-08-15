package ai

import (
	"bytes"
	"encoding/json"
	"finance-ia/internal/domain/ai"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// GroqProvider implements ai.AIProvider using Groq Cloud API
type GroqProvider struct {
	client *http.Client
	apiKey string
	model  string
}

// NewGroqProvider creates a new Groq AI provider
func NewGroqProvider() (*GroqProvider, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY is required")
	}

	return &GroqProvider{
		client: &http.Client{Timeout: 30 * time.Second},
		apiKey: apiKey,
		model:  "llama-3.3-70b-versatile", // Updated to supported Groq model
	}, nil
}

func (p *GroqProvider) GenerateInsight(ctx ai.FinancialContext) (*ai.AIInsight, error) {
	prompt := buildFreeInsightPrompt(ctx)

	response, err := p.generate(prompt)
	if err != nil {
		return nil, err
	}

	return &ai.AIInsight{
		Type:    ai.InsightTypeTip,
		Title:   "💡 Dica Financeira",
		Content: response,
	}, nil
}

func (p *GroqProvider) GenerateFullAnalysis(ctx ai.FinancialContext) ([]*ai.AIInsight, error) {
	prompts := []struct {
		insightType ai.InsightType
		title       string
		prompt      string
	}{
		{
			insightType: ai.InsightTypeWarning,
			title:       "⚠️ Diagnóstico de Saúde Financeira",
			prompt:      buildHealthDiagnosticPrompt(ctx),
		},
		{
			insightType: ai.InsightTypeTip,
			title:       "🎯 Recomendações Prioritárias",
			prompt:      buildRecommendationsPrompt(ctx),
		},
		{
			insightType: ai.InsightTypeProjection,
			title:       "📈 Projeção Financeira",
			prompt:      buildProjectionPrompt(ctx),
		},
	}

	var insights []*ai.AIInsight
	for _, pr := range prompts {
		response, err := p.generate(pr.prompt)
		if err != nil {
			continue
		}
		insights = append(insights, &ai.AIInsight{
			Type:    pr.insightType,
			Title:   pr.title,
			Content: response,
		})
	}

	return insights, nil
}

func (p *GroqProvider) generate(prompt string) (string, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"

	payload := map[string]interface{}{
		"model": p.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Você é um consultor financeiro pessoal preciso e direto. Você responde estritamente em português.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("groq api error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices returned from groq")
	}

	firstChoice := choices[0].(map[string]interface{})
	message := firstChoice["message"].(map[string]interface{})
	content := message["content"].(string)

	return strings.TrimSpace(content), nil
}

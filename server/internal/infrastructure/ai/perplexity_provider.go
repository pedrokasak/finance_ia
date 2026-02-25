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

// PerplexityProvider implements ai.AIProvider using Perplexity API (OpenAI compatible shape)
type PerplexityProvider struct {
	client *http.Client
	apiKey string
	model  string
}

// NewPerplexityProvider creates a new Perplexity AI provider
func NewPerplexityProvider() (*PerplexityProvider, error) {
	apiKey := os.Getenv("PERPLEXITY_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("PERPLEXITY_API_KEY is required")
	}

	return &PerplexityProvider{
		client: &http.Client{Timeout: 30 * time.Second},
		apiKey: apiKey,
		model:  "sonar-pro", // You can switch to any other model
	}, nil
}

func (p *PerplexityProvider) GenerateInsight(ctx ai.FinancialContext) (*ai.AIInsight, error) {
	prompt := buildFreeInsightPrompt(ctx)

	response, err := p.generate(prompt)
	if err != nil {
		return nil, err
	}

	return &ai.AIInsight{
		Type:    ai.InsightTypeTip,
		Title:   "💡 Dica Financeira (Perplexity)",
		Content: response,
	}, nil
}

func (p *PerplexityProvider) GenerateDiagnostic(ctx ai.FinancialContext) (*ai.AIInsight, error) {
	prompt := buildHealthDiagnosticPrompt(ctx)
	response, err := p.generate(prompt)
	if err != nil {
		return nil, err
	}
	return &ai.AIInsight{
		Type:    ai.InsightTypeWarning,
		Title:   "⚠️ Diagnóstico de Saúde Financeira",
		Content: response,
	}, nil
}

func (p *PerplexityProvider) GenerateProjection(ctx ai.FinancialContext) (*ai.AIInsight, error) {
	prompt := buildProjectionPrompt(ctx)
	response, err := p.generate(prompt)
	if err != nil {
		return nil, err
	}
	return &ai.AIInsight{
		Type:    ai.InsightTypeProjection,
		Title:   "📈 Projeção Financeira",
		Content: response,
	}, nil
}

func (p *PerplexityProvider) generate(prompt string) (string, error) {
	url := "https://api.perplexity.ai/chat/completions"

	payload := map[string]interface{}{
		"model": p.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Você é um consultor financeiro pessoal preciso e direto.",
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
		return "", fmt.Errorf("perplexity api error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices returned from perplexity")
	}

	firstChoice := choices[0].(map[string]interface{})
	message := firstChoice["message"].(map[string]interface{})
	content := message["content"].(string)

	return strings.TrimSpace(content), nil
}

package ai

import (
	"bytes"
	"encoding/json"
	"finance-ia/internal/domain/ai"
	"fmt"
	"io/ioutil"
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

type perplexityRequest struct {
	Model    string              `json:"model"`
	Messages []perplexityMessage `json:"messages"`
}

type perplexityMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type perplexityResponse struct {
	Choices []struct {
		Message perplexityMessage `json:"message"`
	} `json:"choices"`
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

func (p *PerplexityProvider) GenerateFullAnalysis(ctx ai.FinancialContext) ([]*ai.AIInsight, error) {
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

func (p *PerplexityProvider) generate(prompt string) (string, error) {
	url := "https://api.perplexity.ai/chat/completions"

	payload := perplexityRequest{
		Model: p.model,
		Messages: []perplexityMessage{
			{
				Role:    "system",
				Content: "Você é um consultor financeiro pessoal preciso e direto.",
			},
			{
				Role:    "user",
				Content: prompt,
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
		respBody, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("perplexity api error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result perplexityResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from perplexity")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

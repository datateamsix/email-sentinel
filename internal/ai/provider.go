package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// sanitizeAPIError removes potential API keys and sensitive data from error messages
// This prevents accidental exposure of credentials in logs or terminal output
func sanitizeAPIError(errorBody string) string {
	// Limit error message length
	const maxLength = 500
	if len(errorBody) > maxLength {
		errorBody = errorBody[:maxLength] + "... (truncated)"
	}

	// Pattern to detect potential API keys (sequences of 20+ alphanumeric chars)
	apiKeyPattern := regexp.MustCompile(`[a-zA-Z0-9_-]{20,}`)
	sanitized := apiKeyPattern.ReplaceAllString(errorBody, "[REDACTED]")

	// Remove common API key field names with their values
	patterns := []string{
		`"api[_-]?key"\s*:\s*"[^"]+`,
		`"apiKey"\s*:\s*"[^"]+`,
		`"token"\s*:\s*"[^"]+`,
		`"authorization"\s*:\s*"[^"]+`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		sanitized = re.ReplaceAllString(sanitized, `"[REDACTED]":"[REDACTED]"`)
	}

	return sanitized
}

// Provider defines the interface for AI providers
type Provider interface {
	GenerateSummary(ctx context.Context, req SummaryRequest) (*SummaryResponse, int, error)
	Name() string
}

// NewProvider creates a provider instance based on configuration
func NewProvider(cfg *Config) (Provider, error) {
	provider := strings.ToLower(cfg.AISummary.Provider)

	switch provider {
	case "claude":
		return &ClaudeProvider{
			apiKey:      cfg.AISummary.API.Claude.APIKey,
			model:       cfg.AISummary.API.Claude.Model,
			maxTokens:   cfg.AISummary.API.Claude.MaxTokens,
			temperature: cfg.AISummary.API.Claude.Temperature,
			prompt:      cfg.AISummary.Prompt,
		}, nil

	case "openai":
		return &OpenAIProvider{
			apiKey:      cfg.AISummary.API.OpenAI.APIKey,
			model:       cfg.AISummary.API.OpenAI.Model,
			maxTokens:   cfg.AISummary.API.OpenAI.MaxTokens,
			temperature: cfg.AISummary.API.OpenAI.Temperature,
			prompt:      cfg.AISummary.Prompt,
		}, nil

	case "gemini":
		return &GeminiProvider{
			apiKey:      cfg.AISummary.API.Gemini.APIKey,
			model:       cfg.AISummary.API.Gemini.Model,
			maxTokens:   cfg.AISummary.API.Gemini.MaxTokens,
			temperature: cfg.AISummary.API.Gemini.Temperature,
			prompt:      cfg.AISummary.Prompt,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// ====================================
// Claude (Anthropic) Provider
// ====================================

type ClaudeProvider struct {
	apiKey      string
	model       string
	maxTokens   int
	temperature float64
	prompt      PromptConfig
}

func (p *ClaudeProvider) Name() string {
	return "claude"
}

func (p *ClaudeProvider) GenerateSummary(ctx context.Context, req SummaryRequest) (*SummaryResponse, int, error) {
	// Build user prompt
	userPrompt := p.buildPrompt(req)

	// Prepare request payload
	payload := map[string]interface{}{
		"model":      p.model,
		"max_tokens": p.maxTokens,
		"temperature": p.temperature,
		"system":     p.prompt.System,
		"messages": []map[string]string{
			{"role": "user", "content": userPrompt},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, 0, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		sanitized := sanitizeAPIError(string(bodyBytes))
		return nil, 0, fmt.Errorf("API error (status %d): %s", resp.StatusCode, sanitized)
	}

	// Parse response
	var claudeResp struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(bodyBytes, &claudeResp); err != nil {
		return nil, 0, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(claudeResp.Content) == 0 {
		return nil, 0, fmt.Errorf("no content in response")
	}

	// Parse JSON from response text
	var summary SummaryResponse
	if err := json.Unmarshal([]byte(claudeResp.Content[0].Text), &summary); err != nil {
		return nil, 0, fmt.Errorf("failed to parse summary JSON: %w", err)
	}

	totalTokens := claudeResp.Usage.InputTokens + claudeResp.Usage.OutputTokens
	return &summary, totalTokens, nil
}

func (p *ClaudeProvider) buildPrompt(req SummaryRequest) string {
	template := p.prompt.UserTemplate
	template = strings.ReplaceAll(template, "{{.MaxLength}}", fmt.Sprintf("%d", req.MaxLength))
	template = strings.ReplaceAll(template, "{{.Sender}}", req.Sender)
	template = strings.ReplaceAll(template, "{{.Subject}}", req.Subject)

	// Use full body if available, otherwise use snippet
	body := req.Body
	if body == "" {
		body = req.Snippet
	}
	template = strings.ReplaceAll(template, "{{.Body}}", body)

	return template
}

// ====================================
// OpenAI Provider
// ====================================

type OpenAIProvider struct {
	apiKey      string
	model       string
	maxTokens   int
	temperature float64
	prompt      PromptConfig
}

func (p *OpenAIProvider) Name() string {
	return "openai"
}

func (p *OpenAIProvider) GenerateSummary(ctx context.Context, req SummaryRequest) (*SummaryResponse, int, error) {
	userPrompt := p.buildPrompt(req)

	payload := map[string]interface{}{
		"model":       p.model,
		"max_tokens":  p.maxTokens,
		"temperature": p.temperature,
		"messages": []map[string]string{
			{"role": "system", "content": p.prompt.System},
			{"role": "user", "content": userPrompt},
		},
		"response_format": map[string]string{"type": "json_object"},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, 0, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		sanitized := sanitizeAPIError(string(bodyBytes))
		return nil, 0, fmt.Errorf("API error (status %d): %s", resp.StatusCode, sanitized)
	}

	var openaiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			TotalTokens int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(bodyBytes, &openaiResp); err != nil {
		return nil, 0, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(openaiResp.Choices) == 0 {
		return nil, 0, fmt.Errorf("no choices in response")
	}

	var summary SummaryResponse
	if err := json.Unmarshal([]byte(openaiResp.Choices[0].Message.Content), &summary); err != nil {
		return nil, 0, fmt.Errorf("failed to parse summary JSON: %w", err)
	}

	return &summary, openaiResp.Usage.TotalTokens, nil
}

func (p *OpenAIProvider) buildPrompt(req SummaryRequest) string {
	template := p.prompt.UserTemplate
	template = strings.ReplaceAll(template, "{{.MaxLength}}", fmt.Sprintf("%d", req.MaxLength))
	template = strings.ReplaceAll(template, "{{.Sender}}", req.Sender)
	template = strings.ReplaceAll(template, "{{.Subject}}", req.Subject)

	body := req.Body
	if body == "" {
		body = req.Snippet
	}
	template = strings.ReplaceAll(template, "{{.Body}}", body)

	return template
}

// ====================================
// Google Gemini Provider
// ====================================

type GeminiProvider struct {
	apiKey      string
	model       string
	maxTokens   int
	temperature float64
	prompt      PromptConfig
}

func (p *GeminiProvider) Name() string {
	return "gemini"
}

func (p *GeminiProvider) GenerateSummary(ctx context.Context, req SummaryRequest) (*SummaryResponse, int, error) {
	userPrompt := p.buildPrompt(req)

	// Combine system and user prompts for Gemini
	fullPrompt := p.prompt.System + "\n\n" + userPrompt

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": fullPrompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":    p.temperature,
			"maxOutputTokens": p.maxTokens,
			"responseMimeType": "application/json",
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", p.model, p.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, 0, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		sanitized := sanitizeAPIError(string(bodyBytes))
		return nil, 0, fmt.Errorf("API error (status %d): %s", resp.StatusCode, sanitized)
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		UsageMetadata struct {
			TotalTokenCount int `json:"totalTokenCount"`
		} `json:"usageMetadata"`
	}

	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
		return nil, 0, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, 0, fmt.Errorf("no content in response")
	}

	var summary SummaryResponse
	if err := json.Unmarshal([]byte(geminiResp.Candidates[0].Content.Parts[0].Text), &summary); err != nil {
		return nil, 0, fmt.Errorf("failed to parse summary JSON: %w", err)
	}

	return &summary, geminiResp.UsageMetadata.TotalTokenCount, nil
}

func (p *GeminiProvider) buildPrompt(req SummaryRequest) string {
	template := p.prompt.UserTemplate
	template = strings.ReplaceAll(template, "{{.MaxLength}}", fmt.Sprintf("%d", req.MaxLength))
	template = strings.ReplaceAll(template, "{{.Sender}}", req.Sender)
	template = strings.ReplaceAll(template, "{{.Subject}}", req.Subject)

	body := req.Body
	if body == "" {
		body = req.Snippet
	}
	template = strings.ReplaceAll(template, "{{.Body}}", body)

	return template
}

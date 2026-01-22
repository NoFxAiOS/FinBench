package benchmark

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

	"FinBench/market"
)

// ChatRequest represents an OpenAI-compatible chat request
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

// ChatMessage represents a chat message
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents an OpenAI-compatible chat response
type ChatResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// LLMClient is a client for calling LLM APIs
type LLMClient struct {
	config  *ModelConfig
	baseURL string
	client  *http.Client
}

// NewLLMClient creates a new LLM client
func NewLLMClient(config *ModelConfig) *LLMClient {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = GetBaseURL(config.Provider)
	}

	return &LLMClient{
		config:  config,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// Chat sends a chat request and returns the response
func (c *LLMClient) Chat(ctx context.Context, prompt string) (string, error) {
	reqBody := ChatRequest{
		Model: c.config.Model,
		Messages: []ChatMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	url := c.baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var result ChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("unmarshal response: %w, body: %s", err, string(body))
	}

	if result.Error != nil {
		return "", fmt.Errorf("API error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response choices")
	}

	return result.Choices[0].Message.Content, nil
}

// BuildIndicatorPrompt builds a prompt for indicator calculation (English version)
func BuildIndicatorPrompt(klines []market.Kline) string {
	var sb strings.Builder

	sb.WriteString("Below is the K-line (candlestick) data sorted from oldest to newest:\n")
	sb.WriteString("Index | Open | High | Low | Close | Volume\n")
	sb.WriteString("------|------|------|-----|-------|--------\n")

	for i, k := range klines {
		sb.WriteString(fmt.Sprintf("%d | %.2f | %.2f | %.2f | %.2f | %.2f\n",
			i+1, k.Open, k.High, k.Low, k.Close, k.Volume))
	}

	sb.WriteString(fmt.Sprintf(`
Based on the %d candlesticks above, calculate the following technical indicators using standard algorithms:

1. MA20 (20-period Simple Moving Average)
2. EMA12 (12-period Exponential Moving Average)
3. EMA26 (26-period Exponential Moving Average)
4. MACD (EMA12 - EMA26)
5. RSI14 (14-period Relative Strength Index, using Wilder's smoothing method)
6. Bollinger Bands (20-period, 2 standard deviations): upper, middle, lower
7. ATR14 (14-period Average True Range, using Wilder's smoothing method)
8. VolumeMA5 (5-period Volume Moving Average)

Return ONLY a JSON object in the following format, with no additional text:
{
  "ma20": number,
  "ema12": number,
  "ema26": number,
  "macd": number,
  "rsi14": number,
  "boll_upper": number,
  "boll_middle": number,
  "boll_lower": number,
  "atr14": number,
  "volume_ma5": number
}

Requirements:
- Round all values to 2 decimal places
- For EMA, use SMA as initial value with multiplier = 2/(period+1)
- For RSI, use Wilder's smoothing method
- Return ONLY the JSON object, no explanations`, len(klines)))

	return sb.String()
}

// ParseIndicatorResponse parses the LLM response into IndicatorResult
func ParseIndicatorResponse(response string) (*IndicatorResult, error) {
	var result IndicatorResult

	// Try direct parse first
	if err := json.Unmarshal([]byte(response), &result); err == nil {
		return &result, nil
	}

	// Extract JSON from response
	re := regexp.MustCompile(`\{[^{}]*"ma20"[^{}]*\}`)
	match := re.FindString(response)
	if match == "" {
		// Try more lenient matching
		start := strings.Index(response, "{")
		end := strings.LastIndex(response, "}")
		if start != -1 && end != -1 && end > start {
			match = response[start : end+1]
		}
	}

	if match == "" {
		return nil, fmt.Errorf("no JSON found in response")
	}

	if err := json.Unmarshal([]byte(match), &result); err != nil {
		return nil, fmt.Errorf("parse JSON failed: %w", err)
	}

	return &result, nil
}

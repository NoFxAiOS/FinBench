package benchmark

// Supported model providers and their default configurations
// These MUST match the nofx trading system for consistency
// Source: github.com/NoFxAiOS/nofx/mcp/*_client.go
const (
	// Provider identifiers
	ProviderDeepSeek = "deepseek"
	ProviderQwen     = "qwen"
	ProviderOpenAI   = "openai"
	ProviderClaude   = "claude"
	ProviderGemini   = "gemini"
	ProviderGrok     = "grok"
	ProviderKimi     = "kimi"

	// Default API endpoints (from nofx mcp clients)
	DefaultDeepSeekBaseURL = "https://api.deepseek.com"
	DefaultQwenBaseURL     = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	DefaultOpenAIBaseURL   = "https://api.openai.com/v1"
	DefaultClaudeBaseURL   = "https://api.anthropic.com/v1"
	DefaultGeminiBaseURL   = "https://generativelanguage.googleapis.com/v1beta/openai"
	DefaultGrokBaseURL     = "https://api.x.ai/v1"
	DefaultKimiBaseURL     = "https://api.moonshot.ai/v1"

	// Default model versions - EXACTLY as defined in nofx/mcp/*_client.go
	// These are the models used in production trading decisions
	DefaultDeepSeekModel = "deepseek-chat"            // nofx/mcp/deepseek_client.go:10
	DefaultQwenModel     = "qwen3-max"                // nofx/mcp/qwen_client.go:10
	DefaultOpenAIModel   = "gpt-5.2"                  // nofx/mcp/openai_client.go:10
	DefaultClaudeModel   = "claude-opus-4-5-20251101" // nofx/mcp/claude_client.go:12
	DefaultGeminiModel   = "gemini-3-pro-preview"     // nofx/mcp/gemini_client.go:10
	DefaultGrokModel     = "grok-3-latest"            // nofx/mcp/grok_client.go:10
	DefaultKimiModel     = "moonshot-v1-auto"         // nofx/mcp/kimi_client.go:10

	// Default K-line count for benchmark (matches nofx debate engine default)
	// Source: nofx/debate/engine.go:310
	DefaultKlineCount = 50
)

// ModelInfo contains metadata about a model for reporting
type ModelInfo struct {
	Provider    string `json:"provider"`
	Model       string `json:"model"`
	Version     string `json:"version,omitempty"`
	DisplayName string `json:"display_name"`
	BaseURL     string `json:"base_url"`
}

// GetDefaultModels returns the default model configurations
// These are the official benchmark models for FinBench, aligned with nofx
func GetDefaultModels() []ModelInfo {
	return []ModelInfo{
		{
			Provider:    ProviderDeepSeek,
			Model:       DefaultDeepSeekModel,
			DisplayName: "DeepSeek-Chat",
			BaseURL:     DefaultDeepSeekBaseURL,
		},
		{
			Provider:    ProviderQwen,
			Model:       DefaultQwenModel,
			DisplayName: "Qwen3-Max",
			BaseURL:     DefaultQwenBaseURL,
		},
		{
			Provider:    ProviderOpenAI,
			Model:       DefaultOpenAIModel,
			DisplayName: "GPT-5.2",
			BaseURL:     DefaultOpenAIBaseURL,
		},
		{
			Provider:    ProviderClaude,
			Model:       DefaultClaudeModel,
			DisplayName: "Claude-Opus-4.5",
			BaseURL:     DefaultClaudeBaseURL,
		},
		{
			Provider:    ProviderGemini,
			Model:       DefaultGeminiModel,
			DisplayName: "Gemini-3-Pro",
			BaseURL:     DefaultGeminiBaseURL,
		},
		{
			Provider:    ProviderGrok,
			Model:       DefaultGrokModel,
			DisplayName: "Grok-3",
			BaseURL:     DefaultGrokBaseURL,
		},
		{
			Provider:    ProviderKimi,
			Model:       DefaultKimiModel,
			DisplayName: "Moonshot-V1",
			BaseURL:     DefaultKimiBaseURL,
		},
	}
}

// GetBaseURL returns the default base URL for a provider
func GetBaseURL(provider string) string {
	switch provider {
	case ProviderDeepSeek:
		return DefaultDeepSeekBaseURL
	case ProviderQwen:
		return DefaultQwenBaseURL
	case ProviderOpenAI:
		return DefaultOpenAIBaseURL
	case ProviderClaude:
		return DefaultClaudeBaseURL
	case ProviderGemini:
		return DefaultGeminiBaseURL
	case ProviderGrok:
		return DefaultGrokBaseURL
	case ProviderKimi:
		return DefaultKimiBaseURL
	default:
		return DefaultOpenAIBaseURL
	}
}

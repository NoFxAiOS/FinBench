package benchmark

import (
	"time"

	"FinBench/market"
)

// Snapshot represents a point-in-time market data snapshot
type Snapshot struct {
	ID        string         `json:"id"`
	Symbol    string         `json:"symbol"`
	Interval  string         `json:"interval"`
	Timestamp int64          `json:"timestamp"`
	Klines    []market.Kline `json:"klines"`
	// TODO: Add orderbook data in future phases
}

// IndicatorResult holds the calculated indicator values
type IndicatorResult struct {
	MA20     float64 `json:"ma20"`
	EMA12    float64 `json:"ema12"`
	EMA26    float64 `json:"ema26"`
	MACD     float64 `json:"macd"`
	RSI14    float64 `json:"rsi14"`
	BOLLUp   float64 `json:"boll_upper"`
	BOLLMid  float64 `json:"boll_middle"`
	BOLLLow  float64 `json:"boll_lower"`
	ATR14    float64 `json:"atr14"`
	VolumeMA float64 `json:"volume_ma5"`
}

// IndicatorScores holds scores for each indicator
type IndicatorScores struct {
	MA20     float64 `json:"ma20"`
	EMA12    float64 `json:"ema12"`
	EMA26    float64 `json:"ema26"`
	MACD     float64 `json:"macd"`
	RSI14    float64 `json:"rsi14"`
	BOLLUp   float64 `json:"boll_upper"`
	BOLLMid  float64 `json:"boll_middle"`
	BOLLLow  float64 `json:"boll_lower"`
	ATR14    float64 `json:"atr14"`
	VolumeMA float64 `json:"volume_ma5"`
}

// BenchResult holds the benchmark result for a single model run
type BenchResult struct {
	SnapshotID string             `json:"snapshot_id"`
	Model      string             `json:"model"`
	ModelInfo  *ModelInfo         `json:"model_info"`
	RunIndex   int                `json:"run_index"`
	Expected   *IndicatorResult   `json:"expected"`
	Actual     *IndicatorResult   `json:"actual"`
	Errors     map[string]float64 `json:"errors"`
	Scores     *IndicatorScores   `json:"scores"`
	TotalScore float64            `json:"total_score"`
	Latency    time.Duration      `json:"latency"`
	RawOutput  string             `json:"raw_output"`
	Error      string             `json:"error,omitempty"`
}

// BenchConfig holds benchmark configuration
type BenchConfig struct {
	Mode       string        `json:"mode"`        // "static" | "realtime"
	DatasetDir string        `json:"dataset_dir"` // For static mode
	Symbols    []string      `json:"symbols"`     // For realtime mode
	Interval   string        `json:"interval"`    // K-line interval
	KlineCount int           `json:"kline_count"` // Number of K-lines
	Models     []ModelConfig `json:"models"`
	Runs       int           `json:"runs"` // Number of runs per model for statistical analysis
}

// ModelConfig holds configuration for a single LLM
type ModelConfig struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Model    string `json:"model"`
	APIKey   string `json:"api_key"`
	BaseURL  string `json:"base_url,omitempty"`
}

// BenchReport holds the full benchmark report
type BenchReport struct {
	ID          string             `json:"id"`
	Version     string             `json:"version"`
	Timestamp   time.Time          `json:"timestamp"`
	Config      *BenchConfig       `json:"config"`
	Environment *EnvironmentInfo   `json:"environment"`
	Snapshots   []*Snapshot        `json:"snapshots"`
	Results     []*BenchResult     `json:"results"`
	Statistics  []*ModelStatistics `json:"statistics"`
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
}

// EnvironmentInfo holds information about the benchmark environment
type EnvironmentInfo struct {
	FinBenchVersion string    `json:"finbench_version"`
	GoVersion       string    `json:"go_version"`
	Platform        string    `json:"platform"`
	Timestamp       time.Time `json:"timestamp"`
	Timezone        string    `json:"timezone"`
}

// ModelStatistics holds statistical analysis for a model across multiple runs
type ModelStatistics struct {
	Model         string    `json:"model"`
	ModelInfo     ModelInfo `json:"model_info"`
	RunCount      int       `json:"run_count"`
	SuccessCount  int       `json:"success_count"`
	FailureCount  int       `json:"failure_count"`
	AvgScore      float64   `json:"avg_score"`
	MinScore      float64   `json:"min_score"`
	MaxScore      float64   `json:"max_score"`
	StdDev        float64   `json:"std_dev"`
	AvgLatencyMs  float64   `json:"avg_latency_ms"`
	MinLatencyMs  float64   `json:"min_latency_ms"`
	MaxLatencyMs  float64   `json:"max_latency_ms"`
	Consistency   float64   `json:"consistency"` // 100 - (StdDev / AvgScore * 100)
	IndicatorAvgs map[string]float64 `json:"indicator_avgs"`
}

// LeaderboardEntry represents a model's ranking
type LeaderboardEntry struct {
	Rank        int     `json:"rank"`
	Model       string  `json:"model"`
	Provider    string  `json:"provider"`
	ModelID     string  `json:"model_id"`
	AvgScore    float64 `json:"avg_score"`
	StdDev      float64 `json:"std_dev"`
	Consistency float64 `json:"consistency"`
	AvgLatency  float64 `json:"avg_latency_ms"`
	RunCount    int     `json:"run_count"`
}

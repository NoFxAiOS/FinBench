package benchmark

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"sort"
	"sync"
	"time"

	"FinBench/logger"
)

const Version = "1.0.0"

// Engine is the main benchmark orchestrator
type Engine struct {
	config *BenchConfig
}

// NewEngine creates a new benchmark engine
func NewEngine(config *BenchConfig) *Engine {
	if config.Runs <= 0 {
		config.Runs = 1
	}
	return &Engine{config: config}
}

// Run executes the benchmark
func (e *Engine) Run(ctx context.Context) (*BenchReport, error) {
	report := &BenchReport{
		ID:        time.Now().Format("20060102_150405"),
		Version:   Version,
		Timestamp: time.Now(),
		Config:    e.config,
		Environment: &EnvironmentInfo{
			FinBenchVersion: Version,
			GoVersion:       runtime.Version(),
			Platform:        runtime.GOOS + "/" + runtime.GOARCH,
			Timestamp:       time.Now().UTC(),
			Timezone:        time.Now().Location().String(),
		},
	}

	// Step 1: Get snapshots
	var snapshots []*Snapshot
	var err error

	if e.config.Mode == "static" {
		logger.Infof("Loading snapshots from %s", e.config.DatasetDir)
		snapshots, err = LoadSnapshots(e.config.DatasetDir)
		if err != nil {
			return nil, fmt.Errorf("load snapshots: %w", err)
		}
		if len(snapshots) == 0 {
			return nil, fmt.Errorf("no snapshots found in %s", e.config.DatasetDir)
		}
	} else {
		logger.Infof("Capturing realtime snapshots for symbols: %v", e.config.Symbols)
		for _, symbol := range e.config.Symbols {
			snapshot, err := CaptureSnapshot(symbol, e.config.Interval, e.config.KlineCount)
			if err != nil {
				logger.Errorf("Failed to capture %s: %v", symbol, err)
				continue
			}
			snapshots = append(snapshots, snapshot)
			logger.Infof("Captured snapshot: %s", snapshot.ID)

			// Save snapshot for reproducibility
			if err := SaveSnapshot(snapshot, "datasets/snapshots"); err != nil {
				logger.Warnf("Failed to save snapshot: %v", err)
			}
		}
	}

	report.Snapshots = snapshots

	// Step 2: Run benchmarks for each snapshot and model
	var results []*BenchResult
	var mu sync.Mutex

	totalRuns := len(snapshots) * len(e.config.Models) * e.config.Runs
	logger.Infof("Starting benchmark: %d snapshots x %d models x %d runs = %d total runs",
		len(snapshots), len(e.config.Models), e.config.Runs, totalRuns)

	for _, snapshot := range snapshots {
		// Calculate expected results (ground truth)
		expected := CalculateIndicators(snapshot.Klines)
		prompt := BuildIndicatorPrompt(snapshot.Klines)

		logger.Infof("Benchmarking snapshot %s", snapshot.ID)

		for runIdx := 0; runIdx < e.config.Runs; runIdx++ {
			if e.config.Runs > 1 {
				logger.Infof("  Run %d/%d", runIdx+1, e.config.Runs)
			}

			// Run all models concurrently for this snapshot/run
			var wg sync.WaitGroup
			for _, modelCfg := range e.config.Models {
				wg.Add(1)
				go func(mc ModelConfig, run int) {
					defer wg.Done()

					result := e.runSingleBenchmark(ctx, snapshot.ID, &mc, prompt, expected, run)

					mu.Lock()
					results = append(results, result)
					mu.Unlock()

					if result.Error != "" {
						logger.Errorf("    %s: ERROR - %s", mc.Name, result.Error)
					} else {
						logger.Infof("    %s: Score=%.1f Latency=%v", mc.Name, result.TotalScore, result.Latency)
					}
				}(modelCfg, runIdx)
			}
			wg.Wait()

			// Small delay between runs to avoid rate limiting
			if runIdx < e.config.Runs-1 {
				time.Sleep(500 * time.Millisecond)
			}
		}
	}

	report.Results = results

	// Step 3: Calculate statistics for each model
	report.Statistics = e.calculateStatistics(results)

	// Step 4: Calculate leaderboard
	report.Leaderboard = e.calculateLeaderboard(report.Statistics)

	return report, nil
}

// runSingleBenchmark runs a benchmark for a single model on a single snapshot
func (e *Engine) runSingleBenchmark(ctx context.Context, snapshotID string, modelCfg *ModelConfig, prompt string, expected *IndicatorResult, runIndex int) *BenchResult {
	result := &BenchResult{
		SnapshotID: snapshotID,
		Model:      modelCfg.Name,
		ModelInfo: &ModelInfo{
			Provider:    modelCfg.Provider,
			Model:       modelCfg.Model,
			DisplayName: modelCfg.Name,
			BaseURL:     modelCfg.BaseURL,
		},
		RunIndex: runIndex,
		Expected: expected,
	}

	client := NewLLMClient(modelCfg)

	start := time.Now()
	response, err := client.Chat(ctx, prompt)
	result.Latency = time.Since(start)
	result.RawOutput = response

	if err != nil {
		result.Error = err.Error()
		return result
	}

	actual, err := ParseIndicatorResponse(response)
	if err != nil {
		result.Error = fmt.Sprintf("parse response: %v", err)
		return result
	}

	result.Actual = actual
	result.Scores, result.Errors = ScoreIndicators(expected, actual)
	result.TotalScore = CalculateTotalScore(result.Scores)

	return result
}

// calculateStatistics computes statistics for each model
func (e *Engine) calculateStatistics(results []*BenchResult) []*ModelStatistics {
	// Group results by model
	modelResults := make(map[string][]*BenchResult)
	modelInfos := make(map[string]*ModelInfo)

	for _, r := range results {
		modelResults[r.Model] = append(modelResults[r.Model], r)
		if r.ModelInfo != nil {
			modelInfos[r.Model] = r.ModelInfo
		}
	}

	var stats []*ModelStatistics
	for model, rs := range modelResults {
		stat := &ModelStatistics{
			Model:         model,
			RunCount:      len(rs),
			IndicatorAvgs: make(map[string]float64),
		}

		if info, ok := modelInfos[model]; ok {
			stat.ModelInfo = *info
		}

		var scores []float64
		var latencies []float64
		indicatorSums := make(map[string]float64)
		indicatorCounts := make(map[string]int)

		for _, r := range rs {
			if r.Error != "" {
				stat.FailureCount++
				continue
			}
			stat.SuccessCount++
			scores = append(scores, r.TotalScore)
			latencies = append(latencies, float64(r.Latency.Milliseconds()))

			// Aggregate indicator scores
			if r.Scores != nil {
				indicatorSums["ma20"] += r.Scores.MA20
				indicatorSums["ema12"] += r.Scores.EMA12
				indicatorSums["ema26"] += r.Scores.EMA26
				indicatorSums["macd"] += r.Scores.MACD
				indicatorSums["rsi14"] += r.Scores.RSI14
				indicatorSums["boll_upper"] += r.Scores.BOLLUp
				indicatorSums["boll_middle"] += r.Scores.BOLLMid
				indicatorSums["boll_lower"] += r.Scores.BOLLLow
				indicatorSums["atr14"] += r.Scores.ATR14
				indicatorSums["volume_ma5"] += r.Scores.VolumeMA
				for k := range indicatorSums {
					indicatorCounts[k]++
				}
			}
		}

		if len(scores) > 0 {
			stat.AvgScore = average(scores)
			stat.MinScore = min(scores)
			stat.MaxScore = max(scores)
			stat.StdDev = stdDev(scores)
			stat.AvgLatencyMs = average(latencies)
			stat.MinLatencyMs = min(latencies)
			stat.MaxLatencyMs = max(latencies)

			// Calculate consistency (higher is better)
			if stat.AvgScore > 0 {
				stat.Consistency = 100 - (stat.StdDev/stat.AvgScore)*100
				if stat.Consistency < 0 {
					stat.Consistency = 0
				}
			}

			// Calculate indicator averages
			for k, sum := range indicatorSums {
				if count := indicatorCounts[k]; count > 0 {
					stat.IndicatorAvgs[k] = sum / float64(count)
				}
			}
		}

		stats = append(stats, stat)
	}

	// Sort by average score
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].AvgScore > stats[j].AvgScore
	})

	return stats
}

// calculateLeaderboard creates a ranked leaderboard from statistics
func (e *Engine) calculateLeaderboard(stats []*ModelStatistics) []LeaderboardEntry {
	var entries []LeaderboardEntry

	for _, s := range stats {
		entries = append(entries, LeaderboardEntry{
			Model:       s.Model,
			Provider:    s.ModelInfo.Provider,
			ModelID:     s.ModelInfo.Model,
			AvgScore:    s.AvgScore,
			StdDev:      s.StdDev,
			Consistency: s.Consistency,
			AvgLatency:  s.AvgLatencyMs,
			RunCount:    s.RunCount,
		})
	}

	// Sort by score (descending)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].AvgScore > entries[j].AvgScore
	})

	// Assign ranks
	for i := range entries {
		entries[i].Rank = i + 1
	}

	return entries
}

func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	m := values[0]
	for _, v := range values[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

func max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	m := values[0]
	for _, v := range values[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

func stdDev(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	avg := average(values)
	sumSquares := 0.0
	for _, v := range values {
		diff := v - avg
		sumSquares += diff * diff
	}
	return math.Sqrt(sumSquares / float64(len(values)-1))
}

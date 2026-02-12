# FinBench Design Document

A professional benchmark suite for evaluating Large Language Models on financial indicator calculations.

## Objective

Evaluate LLM accuracy in calculating financial technical indicators, establishing an industry-standard benchmark for the financial AI domain.

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      FinBench                           │
├─────────────────────────────────────────────────────────┤
│  Data Layer (from nofx)                                 │
│  ├── provider/coinank  → Real-time K-line, orderbook    │
│  └── market/           → Local indicator calculations   │
├─────────────────────────────────────────────────────────┤
│  Benchmark Layer                                        │
│  ├── benchmark/        → Core benchmark engine          │
│  ├── models/           → Model configurations           │
│  └── scorer/           → Tiered scoring system          │
├─────────────────────────────────────────────────────────┤
│  Output Layer                                           │
│  ├── report/           → JSON/Markdown reports          │
│  └── cmd/              → CLI interface                  │
└─────────────────────────────────────────────────────────┘
```

## Supported Models

FinBench supports 7 major LLM providers, aligned with the nofx trading system.

**Model versions are sourced directly from nofx/mcp/*_client.go**

| Provider | Model ID | Source |
|----------|----------|--------|
| DeepSeek | deepseek-chat | nofx/mcp/deepseek_client.go:10 |
| Qwen | qwen3-max | nofx/mcp/qwen_client.go:10 |
| OpenAI | gpt-5.2 | nofx/mcp/openai_client.go:10 |
| Claude | claude-opus-4-6 | nofx/mcp/claude_client.go:12 |
| Gemini | gemini-3-pro-preview | nofx/mcp/gemini_client.go:10 |
| Grok | grok-3-latest | nofx/mcp/grok_client.go:10 |
| Kimi | moonshot-v1-auto | nofx/mcp/kimi_client.go:10 |

**Default K-line count: 50** (Source: nofx/debate/engine.go:310)

## Phase 1: 10 Core Indicators

| Indicator | Calculation | Full Score Threshold |
|-----------|-------------|---------------------|
| MA20 | Simple Moving Average | Error ≤0.1% |
| EMA12 | Exponential Moving Average | Error ≤0.1% |
| EMA26 | Exponential Moving Average | Error ≤0.1% |
| MACD | EMA12 - EMA26 | Error ≤0.5% |
| RSI14 | Wilder's Smoothing | Error ≤1% |
| BOLL Upper | MA20 + 2σ | Error ≤0.1% |
| BOLL Lower | MA20 - 2σ | Error ≤0.1% |
| ATR14 | Average True Range | Error ≤1% |
| Volume MA5 | 5-period Volume Average | Error ≤0.1% |

## Tiered Scoring System

| Error Range | Score |
|-------------|-------|
| ≤0.1% | 100 |
| 0.1-1% | 80 |
| 1-5% | 60 |
| 5-10% | 40 |
| >10% | 0 |

## Benchmark Flow

```
1. Capture Snapshot
   CoinAnk API → K-lines (30) + Orderbook → Save snapshot with ID

2. Calculate Ground Truth
   market.CalculateEMA/RSI/BOLL... → 10 indicator values

3. Concurrent Model Calls
   ┌─ DeepSeek ──┐
   ├─ Qwen ──────┤
   ├─ GPT-4o ────┤ → Same prompt → Parse JSON response
   ├─ Claude ────┤
   ├─ Gemini ────┤
   ├─ Grok ──────┤
   └─ Kimi ──────┘

4. Score & Generate Report
   Compare with ground truth → Calculate errors → Tiered scoring → Leaderboard
```

## Statistical Analysis

For rigorous benchmarking, multiple runs are supported:

- **Average Score**: Mean score across all runs
- **Standard Deviation**: Score variability measure
- **Consistency**: 100 - (StdDev / AvgScore * 100)
- **Min/Max Scores**: Range of performance
- **Latency Statistics**: Response time analysis

## Data Structures

```go
// Benchmark snapshot
type Snapshot struct {
    ID        string
    Symbol    string      // e.g., BTCUSDT
    Timestamp int64
    Klines    []Kline     // 30 candlesticks
    Orderbook Orderbook   // Market depth (future)
}

// Benchmark result
type BenchResult struct {
    SnapshotID string
    Model      string              // e.g., deepseek-chat
    ModelInfo  ModelInfo           // Provider details
    Scores     map[string]float64  // Per-indicator scores
    TotalScore float64             // Average score
    Latency    time.Duration       // Response time
    RawOutput  string              // Raw LLM response
}

// Model statistics
type ModelStatistics struct {
    Model        string
    RunCount     int
    AvgScore     float64
    StdDev       float64
    Consistency  float64
    AvgLatencyMs float64
}
```

## Dataset Modes

### Static Mode (Reproducible)
```
datasets/
├── snapshots/
│   ├── 20250122_143052_BTCUSDT_1h.json
│   ├── 20250122_143052_ETHUSDT_1h.json
│   └── ...
└── index.json
```

### Realtime Mode (Live Data)
- Captures fresh market data from CoinAnk
- Automatically saves snapshots for future reproducibility
- Tests real-world performance

## CLI Interface

```bash
# Run benchmark with multiple runs for statistical analysis
finbench run -config=config.json -symbols=BTCUSDT,ETHUSDT -runs=3

# Use static dataset for reproducibility
finbench run -mode=static -dataset=datasets/snapshots -config=config.json

# Capture snapshots only
finbench snapshot -symbols=BTCUSDT,ETHUSDT

# List supported models
finbench models
```

## Report Format

Reports are generated in JSON format with:

- Environment info (FinBench version, Go version, platform, timezone)
- Complete configuration used
- All snapshots with raw data
- Individual run results
- Statistical analysis per model
- Final leaderboard ranking

## Future Phases

- **Phase 2**: Reasoning explanation tasks (evaluate logic correctness)
- **Phase 3**: Code generation tasks (run and verify results)
- **Full Version**: 50+ indicators coverage

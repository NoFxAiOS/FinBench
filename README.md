# FinBench

**Financial Technical Indicator Calculation LLM Benchmark**

FinBench is a professional LLM benchmarking tool for evaluating the capabilities of large language models in calculating financial technical indicators. Through standardized evaluation processes, it provides model selection guidance for AI quantitative trading systems.

## Features

- **Multi-Model Support** - Supports 7 major LLM providers
- **10 Indicator Tests** - Covers common technical analysis indicators
- **Statistical Analysis** - Multiple runs for stability data
- **Visual Reports** - Generates professional HTML chart reports
- **Reproducibility** - Supports static dataset mode

## Supported Models

| Provider | Model | API Endpoint |
|----------|-------|--------------|
| DeepSeek | deepseek-chat | api.deepseek.com |
| Qwen | qwen3-max | dashscope.aliyuncs.com |
| OpenAI | gpt-5.2 | api.openai.com |
| Claude | claude-opus-4-5-20251101 | api.anthropic.com |
| Gemini | gemini-3-pro-preview | generativelanguage.googleapis.com |
| Grok | grok-3-latest | api.x.ai |
| Kimi | moonshot-v1-auto | api.moonshot.ai |

## Evaluation Indicators

| Indicator | Description | Score Threshold |
|-----------|-------------|-----------------|
| MA20 | 20-period Simple Moving Average | ≤0.1% full score |
| EMA12 | 12-period Exponential Moving Average | ≤0.1% full score |
| EMA26 | 26-period Exponential Moving Average | ≤0.1% full score |
| MACD | EMA12 - EMA26 | ≤0.5% full score |
| RSI14 | 14-period Relative Strength Index | ≤1% full score |
| BOLL | Bollinger Bands (Upper/Middle/Lower) | ≤0.1% full score |
| ATR14 | 14-period Average True Range | ≤1% full score |
| VolumeMA5 | 5-period Volume Moving Average | ≤0.1% full score |

### Scoring Rules

| Error Range | Score |
|-------------|-------|
| ≤ 0.1% | 100 |
| 0.1% - 1% | 80 |
| 1% - 5% | 60 |
| 5% - 10% | 40 |
| > 10% | 0 |

## Quick Start

### Installation

```bash
git clone https://github.com/NoFxAiOS/FinBench.git
cd FinBench
go build -o finbench ./cmd/finbench
```

### Configuration

Copy the configuration template and fill in API Keys:

```bash
cp config.template.json config.json
```

Edit `config.json`:

```json
{
  "models": [
    {
      "name": "DeepSeek-Chat",
      "provider": "deepseek",
      "model": "deepseek-chat",
      "api_key": "your-api-key-here"
    }
  ]
}
```

### Running Benchmarks

```bash
# Standard benchmark (10 runs)
./finbench run -config=config.json -symbols=BTCUSDT -runs=10 -output=report.json

# Quick benchmark (3 runs)
./finbench run -config=config.json -symbols=BTCUSDT -runs=3

# Multi-symbol benchmark
./finbench run -config=config.json -symbols=BTCUSDT,ETHUSDT -runs=10

# Using static dataset (reproducible)
./finbench run -mode=static -dataset=datasets/snapshots -config=config.json
```

### Generating Visual Reports

```bash
python3 scripts/generate_report.py report.json finbench_report.html
open finbench_report.html
```

## Command Reference

```bash
# View supported models
./finbench models

# Capture market data snapshot
./finbench snapshot -symbols=BTCUSDT,ETHUSDT -output=datasets/snapshots

# View help
./finbench help
```

### Run Command Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `-config` | Configuration file path | Required |
| `-mode` | Benchmark mode (realtime/static) | realtime |
| `-symbols` | Trading pairs (comma-separated) | BTCUSDT |
| `-interval` | Kline interval | 1h |
| `-klines` | Number of klines | 50 |
| `-runs` | Runs per model | 1 |
| `-output` | Output report path | - |

## Report Examples

HTML reports include:

- **Leaderboard** - Overall model rankings
- **Score Bar Chart** - Visual comparison
- **Radar Chart** - Indicator capability analysis
- **Heatmap** - Indicator score matrix
- **Latency Comparison** - Response speed analysis
- **Detail Cards** - Per-model statistics

## Evaluation Methodology

```
┌─────────────────────────────────────────┐
│         Unified Data Source (Klines)    │
│        CoinAnk API / Static Snapshots   │
└─────────────────┬───────────────────────┘
                  │
        ┌─────────┴─────────┐
        ▼                   ▼
┌───────────────┐   ┌───────────────┐
│ Local Calc    │   │  LLM Calc     │
│ (Ground Truth)│   │ (Under Test)  │
└───────┬───────┘   └───────┬───────┘
        │                   │
        └─────────┬─────────┘
                  ▼
        ┌───────────────────┐
        │   Compare & Score │
        │   Error → Score   │
        └───────────────────┘
```

## Project Structure

```
FinBench/
├── cmd/finbench/       # CLI entry point
├── benchmark/          # Benchmark engine
│   ├── engine.go       # Main engine
│   ├── llm.go          # LLM client
│   ├── models.go       # Model configuration
│   ├── scorer.go       # Scoring logic
│   └── calculator.go   # Indicator calculation
├── market/             # Market data
│   ├── kline.go        # Kline fetching
│   └── indicators.go   # Indicator implementation
├── provider/           # Data providers
│   └── coinank/        # CoinAnk API
├── scripts/            # Utility scripts
│   └── generate_report.py  # Report generator
├── datasets/           # Datasets
│   └── snapshots/      # Snapshot data
└── docs/               # Documentation
```

## Security Notice

**Do not commit files containing API Keys**

The following files are ignored by `.gitignore`:
- `config.json` / `config*.json` (except template)
- `*_report.json` / `report.json`
- `finbench_report.html`

## License

MIT License

## Related Projects

- [nofx](https://github.com/NoFxAiOS/nofx) - AI Quantitative Trading System

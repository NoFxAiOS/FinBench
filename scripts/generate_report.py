#!/usr/bin/env python3
"""
FinBench HTML Report Generator
Generates professional visualization reports from benchmark JSON data
"""

import json
import sys
from datetime import datetime
from pathlib import Path

def load_report(json_path: str) -> dict:
    """Load benchmark report from JSON file"""
    with open(json_path, 'r') as f:
        return json.load(f)

def generate_html_report(data: dict) -> str:
    """Generate complete HTML report with charts"""

    # Extract data
    leaderboard = data.get('leaderboard', [])
    statistics = data.get('statistics', [])
    config = data.get('config', {})
    environment = data.get('environment', {})
    timestamp = data.get('timestamp', datetime.now().isoformat())

    # Filter out failed models (score = 0)
    active_stats = [s for s in statistics if s.get('avg_score', 0) > 0]
    active_leaderboard = [l for l in leaderboard if l.get('avg_score', 0) > 0]

    # Prepare chart data
    model_names = [s['model'] for s in active_stats]
    avg_scores = [s['avg_score'] for s in active_stats]
    consistencies = [s['consistency'] for s in active_stats]
    latencies = [s['avg_latency_ms'] for s in active_stats]

    # Indicator names
    indicators = ['ma20', 'ema12', 'ema26', 'macd', 'rsi14', 'boll_upper', 'boll_middle', 'boll_lower', 'atr14', 'volume_ma5']
    indicator_labels = ['MA20', 'EMA12', 'EMA26', 'MACD', 'RSI14', 'BOLL_Up', 'BOLL_Mid', 'BOLL_Low', 'ATR14', 'Vol_MA5']

    # Build indicator data for each model
    indicator_data = {}
    for stat in active_stats:
        model = stat['model']
        indicator_data[model] = [stat.get('indicator_avgs', {}).get(ind, 0) for ind in indicators]

    # Generate colors for models
    colors = [
        'rgba(54, 162, 235, 0.8)',   # blue
        'rgba(255, 99, 132, 0.8)',   # red
        'rgba(75, 192, 192, 0.8)',   # teal
        'rgba(255, 206, 86, 0.8)',   # yellow
        'rgba(153, 102, 255, 0.8)',  # purple
        'rgba(255, 159, 64, 0.8)',   # orange
        'rgba(199, 199, 199, 0.8)',  # gray
        'rgba(83, 102, 255, 0.8)',   # indigo
    ]

    border_colors = [c.replace('0.8', '1') for c in colors]

    # Build radar datasets
    radar_datasets = []
    for i, model in enumerate(model_names):
        radar_datasets.append({
            'label': model,
            'data': indicator_data.get(model, [0]*10),
            'backgroundColor': colors[i % len(colors)].replace('0.8', '0.2'),
            'borderColor': border_colors[i % len(border_colors)],
            'pointBackgroundColor': border_colors[i % len(border_colors)],
            'pointBorderColor': '#fff',
            'pointHoverBackgroundColor': '#fff',
            'pointHoverBorderColor': border_colors[i % len(border_colors)]
        })

    # Build heatmap data
    heatmap_data = []
    for i, model in enumerate(model_names):
        for j, ind in enumerate(indicators):
            score = indicator_data.get(model, [0]*10)[j]
            heatmap_data.append({'x': j, 'y': i, 'v': score})

    # Failed models
    failed_models = [s['model'] for s in statistics if s.get('avg_score', 0) == 0]

    html = f'''<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FinBench Benchmark Report</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/chartjs-chart-matrix@1.1.1/dist/chartjs-chart-matrix.min.js"></script>
    <style>
        :root {{
            --primary: #2563eb;
            --success: #10b981;
            --warning: #f59e0b;
            --danger: #ef4444;
            --dark: #1f2937;
            --light: #f3f4f6;
            --card-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
        }}

        * {{
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }}

        body {{
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }}

        .container {{
            max-width: 1400px;
            margin: 0 auto;
        }}

        .header {{
            background: white;
            border-radius: 16px;
            padding: 30px 40px;
            margin-bottom: 24px;
            box-shadow: var(--card-shadow);
            text-align: center;
        }}

        .header h1 {{
            font-size: 2.5rem;
            color: var(--dark);
            margin-bottom: 10px;
        }}

        .header .subtitle {{
            color: #6b7280;
            font-size: 1.1rem;
        }}

        .header .meta {{
            display: flex;
            justify-content: center;
            gap: 30px;
            margin-top: 20px;
            flex-wrap: wrap;
        }}

        .header .meta-item {{
            background: var(--light);
            padding: 8px 16px;
            border-radius: 8px;
            font-size: 0.9rem;
        }}

        .card {{
            background: white;
            border-radius: 16px;
            padding: 24px;
            margin-bottom: 24px;
            box-shadow: var(--card-shadow);
        }}

        .card-title {{
            font-size: 1.4rem;
            color: var(--dark);
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 2px solid var(--light);
            display: flex;
            align-items: center;
            gap: 10px;
        }}

        .card-title::before {{
            content: '';
            width: 4px;
            height: 24px;
            background: var(--primary);
            border-radius: 2px;
        }}

        .grid {{
            display: grid;
            gap: 24px;
        }}

        .grid-2 {{
            grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
        }}

        .grid-3 {{
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
        }}

        /* Leaderboard Table */
        .leaderboard {{
            width: 100%;
            border-collapse: collapse;
        }}

        .leaderboard th {{
            background: var(--dark);
            color: white;
            padding: 14px 16px;
            text-align: left;
            font-weight: 600;
        }}

        .leaderboard th:first-child {{
            border-radius: 8px 0 0 0;
        }}

        .leaderboard th:last-child {{
            border-radius: 0 8px 0 0;
        }}

        .leaderboard td {{
            padding: 14px 16px;
            border-bottom: 1px solid var(--light);
        }}

        .leaderboard tr:hover {{
            background: #f9fafb;
        }}

        .leaderboard .rank {{
            font-weight: bold;
            font-size: 1.2rem;
        }}

        .rank-1 {{ color: #fbbf24; }}
        .rank-2 {{ color: #9ca3af; }}
        .rank-3 {{ color: #d97706; }}

        .score-bar {{
            background: var(--light);
            border-radius: 4px;
            height: 24px;
            position: relative;
            overflow: hidden;
        }}

        .score-bar-fill {{
            height: 100%;
            border-radius: 4px;
            display: flex;
            align-items: center;
            padding-left: 8px;
            color: white;
            font-weight: 600;
            font-size: 0.85rem;
            transition: width 0.5s ease;
        }}

        .score-excellent {{ background: linear-gradient(90deg, #10b981, #34d399); }}
        .score-good {{ background: linear-gradient(90deg, #3b82f6, #60a5fa); }}
        .score-fair {{ background: linear-gradient(90deg, #f59e0b, #fbbf24); }}
        .score-poor {{ background: linear-gradient(90deg, #ef4444, #f87171); }}

        /* Chart containers */
        .chart-container {{
            position: relative;
            height: 350px;
        }}

        .chart-container-large {{
            position: relative;
            height: 450px;
        }}

        /* Model cards */
        .model-card {{
            background: white;
            border-radius: 12px;
            padding: 20px;
            box-shadow: var(--card-shadow);
            border-left: 4px solid var(--primary);
        }}

        .model-card h4 {{
            font-size: 1.1rem;
            color: var(--dark);
            margin-bottom: 12px;
        }}

        .model-card .stats {{
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: 10px;
        }}

        .model-card .stat {{
            background: var(--light);
            padding: 10px;
            border-radius: 8px;
        }}

        .model-card .stat-label {{
            font-size: 0.75rem;
            color: #6b7280;
            text-transform: uppercase;
        }}

        .model-card .stat-value {{
            font-size: 1.2rem;
            font-weight: 600;
            color: var(--dark);
        }}

        /* Failed models */
        .failed-models {{
            background: #fef2f2;
            border: 1px solid #fecaca;
            border-radius: 8px;
            padding: 16px;
            margin-top: 20px;
        }}

        .failed-models h4 {{
            color: var(--danger);
            margin-bottom: 8px;
        }}

        /* Summary stats */
        .summary-stats {{
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 16px;
            margin-bottom: 24px;
        }}

        .summary-stat {{
            background: white;
            border-radius: 12px;
            padding: 20px;
            text-align: center;
            box-shadow: var(--card-shadow);
        }}

        .summary-stat .value {{
            font-size: 2.5rem;
            font-weight: 700;
            color: var(--primary);
        }}

        .summary-stat .label {{
            color: #6b7280;
            margin-top: 8px;
        }}

        /* Footer */
        .footer {{
            text-align: center;
            padding: 20px;
            color: white;
            opacity: 0.8;
            font-size: 0.9rem;
        }}

        /* Indicator heatmap legend */
        .heatmap-legend {{
            display: flex;
            justify-content: center;
            gap: 20px;
            margin-top: 16px;
            flex-wrap: wrap;
        }}

        .legend-item {{
            display: flex;
            align-items: center;
            gap: 6px;
            font-size: 0.85rem;
        }}

        .legend-color {{
            width: 20px;
            height: 20px;
            border-radius: 4px;
        }}

        /* Responsive */
        @media (max-width: 768px) {{
            .header h1 {{
                font-size: 1.8rem;
            }}
            .grid-2, .grid-3 {{
                grid-template-columns: 1fr;
            }}
            .chart-container {{
                height: 300px;
            }}
        }}
    </style>
</head>
<body>
    <div class="container">
        <!-- Header -->
        <div class="header">
            <h1>FinBench Benchmark Report</h1>
            <p class="subtitle">Financial Technical Indicator Calculation LLM Benchmark</p>
            <div class="meta">
                <span class="meta-item">Date: {timestamp[:10]}</span>
                <span class="meta-item">Runs per Model: {config.get('runs', 10)}</span>
                <span class="meta-item">Models Tested: {len(active_stats)}</span>
                <span class="meta-item">Symbols: {', '.join(config.get('symbols', ['BTCUSDT']))}</span>
                <span class="meta-item">Interval: {config.get('interval', '1h')}</span>
            </div>
        </div>

        <!-- Summary Stats -->
        <div class="summary-stats">
            <div class="summary-stat">
                <div class="value">{len(active_stats)}</div>
                <div class="label">Models Tested</div>
            </div>
            <div class="summary-stat">
                <div class="value">{max(avg_scores):.1f}</div>
                <div class="label">Highest Score</div>
            </div>
            <div class="summary-stat">
                <div class="value">{sum(avg_scores)/len(avg_scores):.1f}</div>
                <div class="label">Average Score</div>
            </div>
            <div class="summary-stat">
                <div class="value">{min(latencies):.0f}ms</div>
                <div class="label">Fastest Response</div>
            </div>
        </div>

        <!-- Leaderboard -->
        <div class="card">
            <h2 class="card-title">Leaderboard</h2>
            <table class="leaderboard">
                <thead>
                    <tr>
                        <th>Rank</th>
                        <th>Model</th>
                        <th>Provider</th>
                        <th>Avg Score</th>
                        <th>Consistency</th>
                        <th>Avg Latency</th>
                        <th>Success Rate</th>
                    </tr>
                </thead>
                <tbody>
'''

    # Add leaderboard rows
    for entry in active_leaderboard:
        rank = entry['rank']
        rank_class = f'rank-{rank}' if rank <= 3 else ''
        rank_emoji = {1: '🥇', 2: '🥈', 3: '🥉'}.get(rank, '')

        score = entry['avg_score']
        if score >= 70:
            score_class = 'score-excellent'
        elif score >= 60:
            score_class = 'score-good'
        elif score >= 50:
            score_class = 'score-fair'
        else:
            score_class = 'score-poor'

        # Find success rate from statistics
        stat = next((s for s in statistics if s['model'] == entry['model']), None)
        success_rate = (stat['success_count'] / stat['run_count'] * 100) if stat else 100

        html += f'''
                    <tr>
                        <td class="rank {rank_class}">{rank_emoji} #{rank}</td>
                        <td><strong>{entry['model']}</strong></td>
                        <td>{entry['provider']}</td>
                        <td>
                            <div class="score-bar">
                                <div class="score-bar-fill {score_class}" style="width: {score}%">{score:.1f}</div>
                            </div>
                        </td>
                        <td>{entry['consistency']:.1f}%</td>
                        <td>{entry['avg_latency_ms']:.0f}ms</td>
                        <td>{success_rate:.0f}%</td>
                    </tr>
'''

    html += '''
                </tbody>
            </table>
'''

    # Add failed models section if any
    if failed_models:
        html += f'''
            <div class="failed-models">
                <h4>Failed Models</h4>
                <p>{', '.join(failed_models)} - Failed to complete benchmark due to API errors or network issues</p>
            </div>
'''

    html += '''
        </div>

        <!-- Charts Row 1 -->
        <div class="grid grid-2">
            <div class="card">
                <h2 class="card-title">Score Comparison</h2>
                <div class="chart-container">
                    <canvas id="scoreChart"></canvas>
                </div>
            </div>
            <div class="card">
                <h2 class="card-title">Latency Comparison</h2>
                <div class="chart-container">
                    <canvas id="latencyChart"></canvas>
                </div>
            </div>
        </div>

        <!-- Radar Chart -->
        <div class="card">
            <h2 class="card-title">Indicator Capability Radar</h2>
            <div class="chart-container-large">
                <canvas id="radarChart"></canvas>
            </div>
        </div>

        <!-- Charts Row 2 -->
        <div class="grid grid-2">
            <div class="card">
                <h2 class="card-title">Consistency Comparison</h2>
                <div class="chart-container">
                    <canvas id="consistencyChart"></canvas>
                </div>
            </div>
            <div class="card">
                <h2 class="card-title">Score Distribution</h2>
                <div class="chart-container">
                    <canvas id="distributionChart"></canvas>
                </div>
            </div>
        </div>

        <!-- Indicator Heatmap -->
        <div class="card">
            <h2 class="card-title">Indicator Score Heatmap</h2>
            <div class="chart-container-large">
                <canvas id="heatmapChart"></canvas>
            </div>
            <div class="heatmap-legend">
                <div class="legend-item"><span class="legend-color" style="background: #ef4444"></span> 0-20 Very Poor</div>
                <div class="legend-item"><span class="legend-color" style="background: #f97316"></span> 20-40 Poor</div>
                <div class="legend-item"><span class="legend-color" style="background: #eab308"></span> 40-60 Fair</div>
                <div class="legend-item"><span class="legend-color" style="background: #84cc16"></span> 60-80 Good</div>
                <div class="legend-item"><span class="legend-color" style="background: #22c55e"></span> 80-100 Excellent</div>
            </div>
        </div>

        <!-- Model Detail Cards -->
        <div class="card">
            <h2 class="card-title">Model Details</h2>
            <div class="grid grid-3">
'''

    # Add model detail cards
    for i, stat in enumerate(active_stats):
        border_color = border_colors[i % len(border_colors)].replace('rgba', 'rgb').replace(', 1)', ')')
        html += f'''
                <div class="model-card" style="border-left-color: {border_color}">
                    <h4>{stat['model']}</h4>
                    <div class="stats">
                        <div class="stat">
                            <div class="stat-label">Avg Score</div>
                            <div class="stat-value">{stat['avg_score']:.1f}</div>
                        </div>
                        <div class="stat">
                            <div class="stat-label">Std Dev</div>
                            <div class="stat-value">{stat['std_dev']:.2f}</div>
                        </div>
                        <div class="stat">
                            <div class="stat-label">Max Score</div>
                            <div class="stat-value">{stat['max_score']:.1f}</div>
                        </div>
                        <div class="stat">
                            <div class="stat-label">Min Score</div>
                            <div class="stat-value">{stat['min_score']:.1f}</div>
                        </div>
                        <div class="stat">
                            <div class="stat-label">Avg Latency</div>
                            <div class="stat-value">{stat['avg_latency_ms']:.0f}ms</div>
                        </div>
                        <div class="stat">
                            <div class="stat-label">Consistency</div>
                            <div class="stat-value">{stat['consistency']:.1f}%</div>
                        </div>
                    </div>
                </div>
'''

    html += '''
            </div>
        </div>

        <!-- Indicator Analysis -->
        <div class="card">
            <h2 class="card-title">Indicator Difficulty Analysis</h2>
            <div class="chart-container">
                <canvas id="indicatorDifficultyChart"></canvas>
            </div>
        </div>

        <!-- Footer -->
        <div class="footer">
            <p>Generated by FinBench v''' + environment.get('finbench_version', '1.0.0') + f''' | {timestamp[:19]}</p>
            <p>Platform: {environment.get('platform', 'Unknown')} | Go {environment.get('go_version', 'Unknown')}</p>
        </div>
    </div>

    <script>
        // Chart.js configuration
        Chart.defaults.font.family = "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif";

        // Data
        const modelNames = {json.dumps(model_names)};
        const avgScores = {json.dumps(avg_scores)};
        const consistencies = {json.dumps(consistencies)};
        const latencies = {json.dumps(latencies)};
        const indicatorLabels = {json.dumps(indicator_labels)};
        const colors = {json.dumps(colors)};
        const borderColors = {json.dumps(border_colors)};

        // Score Chart
        new Chart(document.getElementById('scoreChart'), {{
            type: 'bar',
            data: {{
                labels: modelNames,
                datasets: [{{
                    label: 'Average Score',
                    data: avgScores,
                    backgroundColor: colors,
                    borderColor: borderColors,
                    borderWidth: 2,
                    borderRadius: 8
                }}]
            }},
            options: {{
                responsive: true,
                maintainAspectRatio: false,
                plugins: {{
                    legend: {{ display: false }}
                }},
                scales: {{
                    y: {{
                        beginAtZero: true,
                        max: 100,
                        title: {{ display: true, text: 'Score' }}
                    }}
                }}
            }}
        }});

        // Latency Chart
        new Chart(document.getElementById('latencyChart'), {{
            type: 'bar',
            data: {{
                labels: modelNames,
                datasets: [{{
                    label: 'Average Latency (ms)',
                    data: latencies,
                    backgroundColor: 'rgba(239, 68, 68, 0.7)',
                    borderColor: 'rgba(239, 68, 68, 1)',
                    borderWidth: 2,
                    borderRadius: 8
                }}]
            }},
            options: {{
                responsive: true,
                maintainAspectRatio: false,
                plugins: {{
                    legend: {{ display: false }}
                }},
                scales: {{
                    y: {{
                        beginAtZero: true,
                        title: {{ display: true, text: 'Latency (ms)' }}
                    }}
                }}
            }}
        }});

        // Radar Chart
        const radarDatasets = {json.dumps(radar_datasets)};
        new Chart(document.getElementById('radarChart'), {{
            type: 'radar',
            data: {{
                labels: indicatorLabels,
                datasets: radarDatasets
            }},
            options: {{
                responsive: true,
                maintainAspectRatio: false,
                scales: {{
                    r: {{
                        beginAtZero: true,
                        max: 100,
                        ticks: {{ stepSize: 20 }}
                    }}
                }},
                plugins: {{
                    legend: {{
                        position: 'bottom'
                    }}
                }}
            }}
        }});

        // Consistency Chart
        new Chart(document.getElementById('consistencyChart'), {{
            type: 'doughnut',
            data: {{
                labels: modelNames,
                datasets: [{{
                    data: consistencies,
                    backgroundColor: colors,
                    borderColor: borderColors,
                    borderWidth: 2
                }}]
            }},
            options: {{
                responsive: true,
                maintainAspectRatio: false,
                plugins: {{
                    legend: {{
                        position: 'bottom'
                    }},
                    title: {{
                        display: true,
                        text: 'Consistency Distribution (%)'
                    }}
                }}
            }}
        }});

        // Distribution Chart (Box-like visualization using bar with error bars simulation)
        new Chart(document.getElementById('distributionChart'), {{
            type: 'bar',
            data: {{
                labels: modelNames,
                datasets: [{{
                    label: 'Score Range',
                    data: {json.dumps([{'min': s['min_score'], 'max': s['max_score'], 'avg': s['avg_score']} for s in active_stats])}.map(d => d.avg),
                    backgroundColor: colors.map(c => c.replace('0.8', '0.5')),
                    borderColor: borderColors,
                    borderWidth: 2,
                    borderRadius: 8
                }}]
            }},
            options: {{
                responsive: true,
                maintainAspectRatio: false,
                plugins: {{
                    legend: {{ display: false }},
                    tooltip: {{
                        callbacks: {{
                            label: function(context) {{
                                const stats = {json.dumps([{'min': s['min_score'], 'max': s['max_score'], 'avg': s['avg_score'], 'std': s['std_dev']} for s in active_stats])};
                                const s = stats[context.dataIndex];
                                return [`Average: ${{s.avg.toFixed(1)}}`, `Range: ${{s.min.toFixed(1)}} - ${{s.max.toFixed(1)}}`, `Std Dev: ${{s.std.toFixed(2)}}`];
                            }}
                        }}
                    }}
                }},
                scales: {{
                    y: {{
                        beginAtZero: true,
                        max: 100,
                        title: {{ display: true, text: 'Score' }}
                    }}
                }}
            }}
        }});

        // Indicator Difficulty Chart
        const indicatorAvgs = indicatorLabels.map((label, idx) => {{
            const sum = {json.dumps([[indicator_data.get(m, [0]*10)[i] for m in model_names] for i in range(10)])};
            return sum[idx].reduce((a, b) => a + b, 0) / sum[idx].length;
        }});

        new Chart(document.getElementById('indicatorDifficultyChart'), {{
            type: 'bar',
            data: {{
                labels: indicatorLabels,
                datasets: [{{
                    label: 'Average Score',
                    data: indicatorAvgs,
                    backgroundColor: indicatorAvgs.map(v => {{
                        if (v >= 80) return 'rgba(34, 197, 94, 0.7)';
                        if (v >= 60) return 'rgba(132, 204, 22, 0.7)';
                        if (v >= 40) return 'rgba(234, 179, 8, 0.7)';
                        if (v >= 20) return 'rgba(249, 115, 22, 0.7)';
                        return 'rgba(239, 68, 68, 0.7)';
                    }}),
                    borderWidth: 0,
                    borderRadius: 8
                }}]
            }},
            options: {{
                indexAxis: 'y',
                responsive: true,
                maintainAspectRatio: false,
                plugins: {{
                    legend: {{ display: false }},
                    title: {{
                        display: true,
                        text: 'Average Score by Indicator (Lower = Harder)'
                    }}
                }},
                scales: {{
                    x: {{
                        beginAtZero: true,
                        max: 100
                    }}
                }}
            }}
        }});

        // Heatmap Chart
        const heatmapData = {json.dumps(heatmap_data)};
        new Chart(document.getElementById('heatmapChart'), {{
            type: 'matrix',
            data: {{
                datasets: [{{
                    label: 'Score',
                    data: heatmapData,
                    backgroundColor: function(ctx) {{
                        const v = ctx.dataset.data[ctx.dataIndex].v;
                        if (v >= 80) return 'rgba(34, 197, 94, 0.8)';
                        if (v >= 60) return 'rgba(132, 204, 22, 0.8)';
                        if (v >= 40) return 'rgba(234, 179, 8, 0.8)';
                        if (v >= 20) return 'rgba(249, 115, 22, 0.8)';
                        return 'rgba(239, 68, 68, 0.8)';
                    }},
                    borderColor: 'white',
                    borderWidth: 2,
                    width: (ctx) => (ctx.chart.chartArea || {{}}).width / {len(indicators)} - 4,
                    height: (ctx) => (ctx.chart.chartArea || {{}}).height / {len(model_names)} - 4
                }}]
            }},
            options: {{
                responsive: true,
                maintainAspectRatio: false,
                plugins: {{
                    legend: {{ display: false }},
                    tooltip: {{
                        callbacks: {{
                            title: () => '',
                            label: function(ctx) {{
                                const d = ctx.dataset.data[ctx.dataIndex];
                                return `${{modelNames[d.y]}} - ${{indicatorLabels[d.x]}}: ${{d.v.toFixed(1)}}`;
                            }}
                        }}
                    }}
                }},
                scales: {{
                    x: {{
                        type: 'category',
                        labels: indicatorLabels,
                        offset: true,
                        grid: {{ display: false }}
                    }},
                    y: {{
                        type: 'category',
                        labels: modelNames,
                        offset: true,
                        grid: {{ display: false }}
                    }}
                }}
            }}
        }});
    </script>
</body>
</html>
'''

    return html

def main():
    if len(sys.argv) < 2:
        print("Usage: python generate_report.py <benchmark_report.json> [output.html]")
        sys.exit(1)

    json_path = sys.argv[1]
    output_path = sys.argv[2] if len(sys.argv) > 2 else 'finbench_report.html'

    print(f"Loading benchmark data from {json_path}...")
    data = load_report(json_path)

    print("Generating HTML report...")
    html = generate_html_report(data)

    print(f"Writing report to {output_path}...")
    with open(output_path, 'w', encoding='utf-8') as f:
        f.write(html)

    print(f"Report generated: {output_path}")
    print(f"Open in browser to view the interactive report.")

if __name__ == '__main__':
    main()

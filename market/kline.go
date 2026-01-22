package market

import (
	"context"
	"fmt"
	"time"

	"FinBench/provider/coinank/coinank_api"
	"FinBench/provider/coinank/coinank_enum"
)

// Kline represents a candlestick data point
type Kline struct {
	OpenTime  int64   `json:"open_time"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
	CloseTime int64   `json:"close_time"`
}

// GetKlines fetches kline data from CoinAnk API
func GetKlines(symbol, interval string, limit int) ([]Kline, error) {
	coinankInterval, err := parseInterval(interval)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	ts := time.Now().UnixMilli()

	coinankKlines, err := coinank_api.Kline(ctx, symbol, coinank_enum.Binance, ts, coinank_enum.To, limit, coinankInterval)
	if err != nil {
		return nil, fmt.Errorf("CoinAnk API error: %w", err)
	}

	klines := make([]Kline, len(coinankKlines))
	for i, ck := range coinankKlines {
		klines[i] = Kline{
			OpenTime:  ck.StartTime,
			Open:      ck.Open,
			High:      ck.High,
			Low:       ck.Low,
			Close:     ck.Close,
			Volume:    ck.Volume,
			CloseTime: ck.EndTime,
		}
	}

	return klines, nil
}

func parseInterval(interval string) (coinank_enum.Interval, error) {
	switch interval {
	case "1m":
		return coinank_enum.Minute1, nil
	case "3m":
		return coinank_enum.Minute3, nil
	case "5m":
		return coinank_enum.Minute5, nil
	case "15m":
		return coinank_enum.Minute15, nil
	case "30m":
		return coinank_enum.Minute30, nil
	case "1h":
		return coinank_enum.Hour1, nil
	case "2h":
		return coinank_enum.Hour2, nil
	case "4h":
		return coinank_enum.Hour4, nil
	case "6h":
		return coinank_enum.Hour6, nil
	case "8h":
		return coinank_enum.Hour8, nil
	case "12h":
		return coinank_enum.Hour12, nil
	case "1d":
		return coinank_enum.Day1, nil
	case "3d":
		return coinank_enum.Day3, nil
	case "1w":
		return coinank_enum.Week1, nil
	default:
		return "", fmt.Errorf("unsupported interval: %s", interval)
	}
}

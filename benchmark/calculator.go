package benchmark

import (
	"FinBench/market"
)

// CalculateIndicators calculates all indicators from klines (the "expected" values)
func CalculateIndicators(klines []market.Kline) *IndicatorResult {
	result := &IndicatorResult{}

	if len(klines) >= 20 {
		result.MA20 = market.CalculateSMA(klines, 20)
	}

	if len(klines) >= 12 {
		result.EMA12 = market.CalculateEMA(klines, 12)
	}

	if len(klines) >= 26 {
		result.EMA26 = market.CalculateEMA(klines, 26)
		result.MACD = market.CalculateMACD(klines)
	}

	if len(klines) > 14 {
		result.RSI14 = market.CalculateRSI(klines, 14)
		result.ATR14 = market.CalculateATR(klines, 14)
	}

	if len(klines) >= 20 {
		result.BOLLUp, result.BOLLMid, result.BOLLLow = market.CalculateBOLL(klines, 20, 2.0)
	}

	if len(klines) >= 5 {
		result.VolumeMA = market.CalculateVolumeMA(klines, 5)
	}

	return result
}

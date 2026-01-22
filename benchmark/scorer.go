package benchmark

import (
	"math"
)

// ScoreFromError calculates score based on error percentage using tiered scoring
// Error thresholds: ≤0.1% = 100, 0.1-1% = 80, 1-5% = 60, 5-10% = 40, >10% = 0
func ScoreFromError(errorPct float64) float64 {
	errorPct = math.Abs(errorPct)

	switch {
	case errorPct <= 0.1:
		return 100
	case errorPct <= 1:
		return 80
	case errorPct <= 5:
		return 60
	case errorPct <= 10:
		return 40
	default:
		return 0
	}
}

// CalculateError calculates the percentage error between expected and actual values
func CalculateError(expected, actual float64) float64 {
	if expected == 0 {
		if actual == 0 {
			return 0
		}
		return 100 // Expected is 0 but actual is not
	}
	return math.Abs(expected-actual) / math.Abs(expected) * 100
}

// ScoreIndicators compares expected and actual results, returns scores and errors
func ScoreIndicators(expected, actual *IndicatorResult) (*IndicatorScores, map[string]float64) {
	errors := make(map[string]float64)
	scores := &IndicatorScores{}

	// Calculate errors
	errors["ma20"] = CalculateError(expected.MA20, actual.MA20)
	errors["ema12"] = CalculateError(expected.EMA12, actual.EMA12)
	errors["ema26"] = CalculateError(expected.EMA26, actual.EMA26)
	errors["macd"] = CalculateError(expected.MACD, actual.MACD)
	errors["rsi14"] = CalculateError(expected.RSI14, actual.RSI14)
	errors["boll_upper"] = CalculateError(expected.BOLLUp, actual.BOLLUp)
	errors["boll_middle"] = CalculateError(expected.BOLLMid, actual.BOLLMid)
	errors["boll_lower"] = CalculateError(expected.BOLLLow, actual.BOLLLow)
	errors["atr14"] = CalculateError(expected.ATR14, actual.ATR14)
	errors["volume_ma5"] = CalculateError(expected.VolumeMA, actual.VolumeMA)

	// Calculate scores
	scores.MA20 = ScoreFromError(errors["ma20"])
	scores.EMA12 = ScoreFromError(errors["ema12"])
	scores.EMA26 = ScoreFromError(errors["ema26"])
	scores.MACD = ScoreFromError(errors["macd"])
	scores.RSI14 = ScoreFromError(errors["rsi14"])
	scores.BOLLUp = ScoreFromError(errors["boll_upper"])
	scores.BOLLMid = ScoreFromError(errors["boll_middle"])
	scores.BOLLLow = ScoreFromError(errors["boll_lower"])
	scores.ATR14 = ScoreFromError(errors["atr14"])
	scores.VolumeMA = ScoreFromError(errors["volume_ma5"])

	return scores, errors
}

// CalculateTotalScore calculates the average score across all indicators
func CalculateTotalScore(scores *IndicatorScores) float64 {
	total := scores.MA20 + scores.EMA12 + scores.EMA26 + scores.MACD +
		scores.RSI14 + scores.BOLLUp + scores.BOLLMid + scores.BOLLLow +
		scores.ATR14 + scores.VolumeMA
	return total / 10
}

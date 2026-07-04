package scalper

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/trading"
	"testing"
)

func defaultEntryParams() entryParams {
	return entryParams{
		PullbackLookback: 3,
		RsiMomentumLevel: 50,
		RsiOverbought:    65,
		RsiOversold:      35,
	}
}

// newCandle builds one candle; the caller assembles the slice in the project
// order: index 0 is the newest candle.
func newCandle(open, high, low, close float64) domainStructs.DomainCandle {
	return domainStructs.DomainCandle{
		Open:   open,
		High:   high,
		Low:    low,
		Close:  close,
		Volume: 1,
	}
}

func TestDetectBiasLong(t *testing.T) {
	// Rising closes oldest→newest so the linear regression agrees with the EMAs.
	candles := make([]domainStructs.DomainCandle, 5)
	price := 110.0
	for i := range candles {
		candles[i] = newCandle(price-1, price+1, price-2, price)
		price -= 2
	}
	emaFast := []float64{105, 104, 103, 102, 101}
	emaSlow := []float64{100, 100, 100, 100, 100}

	bias, _ := detectBias(candles, emaFast, emaSlow, false)
	if bias != trading.TrendLong {
		t.Errorf("expected bias %s, got %s", trading.TrendLong, bias)
	}
}

func TestDetectBiasShort(t *testing.T) {
	candles := make([]domainStructs.DomainCandle, 5)
	price := 90.0
	for i := range candles {
		candles[i] = newCandle(price+1, price+2, price-1, price)
		price += 2
	}
	emaFast := []float64{95, 96, 97, 98, 99}
	emaSlow := []float64{100, 100, 100, 100, 100}

	bias, _ := detectBias(candles, emaFast, emaSlow, false)
	if bias != trading.TrendShort {
		t.Errorf("expected bias %s, got %s", trading.TrendShort, bias)
	}
}

func TestDetectBiasUnknownWhenPriceBelowSlowEma(t *testing.T) {
	// Fast EMA above slow EMA but the close dipped below the slow EMA: no regime.
	candles := []domainStructs.DomainCandle{
		newCandle(100, 101, 98, 99),
		newCandle(101, 102, 100, 100),
		newCandle(102, 103, 101, 101),
	}
	emaFast := []float64{105, 105, 105}
	emaSlow := []float64{100, 100, 100}

	bias, _ := detectBias(candles, emaFast, emaSlow, false)
	if bias != trading.TrendUnknown {
		t.Errorf("expected bias %s, got %s", trading.TrendUnknown, bias)
	}
}

func TestDetectBiasVetoedByOppositeTrend(t *testing.T) {
	// EMAs say long but the closes fall in a straight line: linear regression
	// is bearish and must veto the regime when confirmation is required.
	candles := make([]domainStructs.DomainCandle, 20)
	price := 101.0
	for i := range candles {
		// oldest candle gets the highest price: falling series oldest→newest
		candles[i] = newCandle(price+1, price+2, price-1, price)
		price += 1
	}
	emaFast := make([]float64, 20)
	emaSlow := make([]float64, 20)
	for i := range emaFast {
		emaFast[i] = 90
		emaSlow[i] = 80
	}

	bias, lrTrend := detectBias(candles, emaFast, emaSlow, true)
	if lrTrend != trading.TrendShort {
		t.Fatalf("expected linear regression trend %s, got %s", trading.TrendShort, lrTrend)
	}
	if bias != trading.TrendUnknown {
		t.Errorf("expected vetoed bias %s, got %s", trading.TrendUnknown, bias)
	}
}

// longSignalFixture builds closed entry candles with a finished pullback:
// price rode above the EMA, dipped to it, and the trigger candle closed back
// above with the RSI flipping through the momentum level.
func longSignalFixture() (candles []domainStructs.DomainCandle, ema []float64, rsi []float64) {
	candles = []domainStructs.DomainCandle{
		newCandle(100.5, 102, 100.4, 101.5), // trigger: bullish, closes above ema
		newCandle(100.2, 100.8, 99.5, 100.3), // pullback: low 99.5 <= ema 100.1
		newCandle(101, 101.5, 100.4, 100.6),
		newCandle(101.5, 102, 101, 101.2),
		newCandle(101, 102, 100.8, 101.8),
	}
	ema = []float64{100.2, 100.1, 100.1, 100.0, 99.9}
	rsi = []float64{55, 45, 48, 52, 60}
	return
}

func TestDetectEntrySignalLong(t *testing.T) {
	candles, ema, rsi := longSignalFixture()

	side := detectEntrySignal(trading.TrendLong, candles, ema, rsi, defaultEntryParams())
	if side != domainStructs.PositionSideLong {
		t.Errorf("expected side %s, got %s", domainStructs.PositionSideLong, side)
	}
}

func TestDetectEntrySignalNoBiasNoSignal(t *testing.T) {
	candles, ema, rsi := longSignalFixture()

	side := detectEntrySignal(trading.TrendUnknown, candles, ema, rsi, defaultEntryParams())
	if side != "" {
		t.Errorf("expected no signal without bias, got %s", side)
	}
}

func TestDetectEntrySignalLongRejectedWithoutPullback(t *testing.T) {
	candles, ema, rsi := longSignalFixture()
	// Remove the pullback: all lows stay above the EMA.
	for i := range candles {
		if i > 0 {
			candles[i].Low = ema[i] + 1
		}
	}

	side := detectEntrySignal(trading.TrendLong, candles, ema, rsi, defaultEntryParams())
	if side != "" {
		t.Errorf("expected no signal without pullback, got %s", side)
	}
}

func TestDetectEntrySignalLongRejectedWithoutRsiFlip(t *testing.T) {
	candles, ema, _ := longSignalFixture()
	// RSI stayed above the momentum level the whole time: the flip is not fresh.
	rsi := []float64{55, 54, 53, 56, 60}

	side := detectEntrySignal(trading.TrendLong, candles, ema, rsi, defaultEntryParams())
	if side != "" {
		t.Errorf("expected no signal without rsi flip, got %s", side)
	}
}

func TestDetectEntrySignalLongRejectedWhenOverbought(t *testing.T) {
	candles, ema, rsi := longSignalFixture()
	rsi[0] = 70 // above the overbought level: too late to enter

	side := detectEntrySignal(trading.TrendLong, candles, ema, rsi, defaultEntryParams())
	if side != "" {
		t.Errorf("expected no signal when overbought, got %s", side)
	}
}

func TestDetectEntrySignalLongRejectedOnBearishTrigger(t *testing.T) {
	candles, ema, rsi := longSignalFixture()
	candles[0] = newCandle(101.5, 102, 100.4, 100.5) // bearish trigger candle

	side := detectEntrySignal(trading.TrendLong, candles, ema, rsi, defaultEntryParams())
	if side != "" {
		t.Errorf("expected no signal on bearish trigger candle, got %s", side)
	}
}

func TestDetectEntrySignalShort(t *testing.T) {
	candles := []domainStructs.DomainCandle{
		newCandle(99.5, 99.6, 98, 98.5),      // trigger: bearish, closes below ema
		newCandle(99.8, 100.5, 99.2, 99.7),   // pullback: high 100.5 >= ema 99.9
		newCandle(99, 99.6, 98.5, 99.4),
		newCandle(98.5, 99, 98, 98.8),
		newCandle(99, 99.2, 98, 98.2),
	}
	ema := []float64{99.8, 99.9, 99.9, 100.0, 100.1}
	rsi := []float64{45, 55, 52, 48, 40}

	side := detectEntrySignal(trading.TrendShort, candles, ema, rsi, defaultEntryParams())
	if side != domainStructs.PositionSideShort {
		t.Errorf("expected side %s, got %s", domainStructs.PositionSideShort, side)
	}
}

func TestDetectEntrySignalNotEnoughData(t *testing.T) {
	candles := []domainStructs.DomainCandle{newCandle(100, 101, 99, 100.5)}
	ema := []float64{100}
	rsi := []float64{55}

	side := detectEntrySignal(trading.TrendLong, candles, ema, rsi, defaultEntryParams())
	if side != "" {
		t.Errorf("expected no signal with not enough data, got %s", side)
	}
}

func TestIsVolatilityEnough(t *testing.T) {
	// ATR 0.1% of price with a 0.05% floor: enough.
	if !isVolatilityEnough(100, 100000, 0.05) {
		t.Error("expected volatility to be enough")
	}
	// ATR 0.01% of price with a 0.05% floor: dead market.
	if isVolatilityEnough(10, 100000, 0.05) {
		t.Error("expected volatility to be not enough")
	}
	if isVolatilityEnough(10, 0, 0.05) {
		t.Error("expected not enough volatility for zero price")
	}
}

func TestIsTpWorthFees(t *testing.T) {
	// TP = 2 x ATR 200 = 0.4% of price, fee floor 0.11% x 3 = 0.33%: worth it.
	if !isTpWorthFees(200, 2, 100000, 0.11, 3) {
		t.Error("expected take profit to cover the fees")
	}
	// TP = 1.5 x ATR 39 = 0.093% of price, fee floor 0.33%: fees eat the profit.
	if isTpWorthFees(39, 1.5, 62900, 0.11, 3) {
		t.Error("expected take profit to not cover the fees")
	}
	if isTpWorthFees(200, 2, 0, 0.11, 3) {
		t.Error("expected not worth fees for zero price")
	}
}

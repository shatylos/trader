package trading

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"math"
	"testing"
)

func TestCalcATRConstantRange(t *testing.T) {
	// Every candle has the same high-low range 2 and closes in the middle,
	// so every true range is 2 and the ATR must be 2.
	candles := make([]domainStructs.DomainCandle, 8)
	for i := range candles {
		candles[i] = domainStructs.DomainCandle{
			Open:  10,
			High:  11,
			Low:   9,
			Close: 10,
		}
	}

	atr, err := CalcATR(candles, 3)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if math.Abs(atr-2) > 1e-9 {
		t.Errorf("expected atr 2, got %f", atr)
	}
}

func TestCalcATRUsesPreviousCloseGap(t *testing.T) {
	// Candles gap up: high-low range is 1 but the distance to the previous
	// close is 10, so the true range must be driven by the gap.
	candles := make([]domainStructs.DomainCandle, 6)
	price := 100.0
	// build oldest→newest, store newest first
	for i := len(candles) - 1; i >= 0; i-- {
		candles[i] = domainStructs.DomainCandle{
			Open:  price,
			High:  price + 1,
			Low:   price,
			Close: price + 1,
		}
		price += 10
	}

	atr, err := CalcATR(candles, 3)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if atr <= 5 {
		t.Errorf("expected atr driven by the gap to previous close, got %f", atr)
	}
}

func TestCalcATRNotEnoughCandles(t *testing.T) {
	candles := make([]domainStructs.DomainCandle, 3)

	_, err := CalcATR(candles, 3)
	if err == nil {
		t.Error("expected error for not enough candles, got nil")
	}
}

func TestCalcATRInvalidPeriod(t *testing.T) {
	candles := make([]domainStructs.DomainCandle, 5)

	_, err := CalcATR(candles, 0)
	if err == nil {
		t.Error("expected error for zero period, got nil")
	}
}

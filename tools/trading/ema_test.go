package trading

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"math"
	"testing"
)

// candlesFromCloses builds candles from close prices given oldest→newest and
// returns them in the project order: index 0 is the newest candle.
func candlesFromCloses(closes ...float64) []domainStructs.DomainCandle {
	candles := make([]domainStructs.DomainCandle, len(closes))
	for i, close := range closes {
		idx := len(closes) - 1 - i
		candles[idx] = domainStructs.DomainCandle{
			Open:   close,
			High:   close,
			Low:    close,
			Close:  close,
			Volume: 1,
		}
	}
	return candles
}

func TestCalcEMASeriesConstantPrice(t *testing.T) {
	candles := candlesFromCloses(10, 10, 10, 10, 10, 10)

	emas, err := CalcEMASeries(candles, 3)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if len(emas) != len(candles) {
		t.Fatalf("expected series length %d, got %d", len(candles), len(emas))
	}
	// A constant price series must produce an EMA equal to the price.
	for i := 0; i <= len(candles)-3; i++ {
		if math.Abs(emas[i]-10) > 1e-9 {
			t.Errorf("expected ema 10 at index %d, got %f", i, emas[i])
		}
	}
	// Positions older than the seed have no EMA.
	if emas[len(emas)-1] != 0 {
		t.Errorf("expected 0 for the oldest position without ema, got %f", emas[len(emas)-1])
	}
}

func TestCalcEMASeriesKnownValues(t *testing.T) {
	// closes oldest→newest: 1, 2, 3, 4, 5; period 3, multiplier 0.5
	// seed = (1+2+3)/3 = 2
	// next = (4-2)*0.5 + 2 = 3
	// next = (5-3)*0.5 + 3 = 4
	candles := candlesFromCloses(1, 2, 3, 4, 5)

	emas, err := CalcEMASeries(candles, 3)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if math.Abs(emas[2]-2) > 1e-9 {
		t.Errorf("expected seed ema 2, got %f", emas[2])
	}
	if math.Abs(emas[1]-3) > 1e-9 {
		t.Errorf("expected ema 3, got %f", emas[1])
	}
	if math.Abs(emas[0]-4) > 1e-9 {
		t.Errorf("expected ema 4, got %f", emas[0])
	}
}

func TestCalcEMAFollowsTrend(t *testing.T) {
	candles := candlesFromCloses(10, 11, 12, 13, 14, 15, 16, 17, 18, 19)

	ema, err := CalcEMA(candles, 5)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	// On a steadily rising series the EMA must lag below the last close but
	// stay above the value from the middle of the window.
	if ema >= 19 || ema <= 14 {
		t.Errorf("expected ema between 14 and 19, got %f", ema)
	}
}

func TestCalcEMANotEnoughCandles(t *testing.T) {
	candles := candlesFromCloses(1, 2)

	_, err := CalcEMASeries(candles, 3)
	if err == nil {
		t.Error("expected error for not enough candles, got nil")
	}
}

func TestCalcEMAInvalidPeriod(t *testing.T) {
	candles := candlesFromCloses(1, 2, 3)

	_, err := CalcEMASeries(candles, 0)
	if err == nil {
		t.Error("expected error for zero period, got nil")
	}
}

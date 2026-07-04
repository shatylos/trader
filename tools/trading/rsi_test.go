package trading

import (
	"math"
	"testing"
)

func TestCalcRSISeriesAllGains(t *testing.T) {
	candles := candlesFromCloses(1, 2, 3, 4, 5, 6, 7, 8)

	rsis, err := CalcRSISeries(candles, 3)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if len(rsis) != len(candles) {
		t.Fatalf("expected series length %d, got %d", len(candles), len(rsis))
	}
	// A series with only gains must produce RSI 100.
	if math.Abs(rsis[0]-100) > 1e-9 {
		t.Errorf("expected rsi 100, got %f", rsis[0])
	}
}

func TestCalcRSISeriesAllLosses(t *testing.T) {
	candles := candlesFromCloses(8, 7, 6, 5, 4, 3, 2, 1)

	rsis, err := CalcRSISeries(candles, 3)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	// A series with only losses must produce RSI 0.
	if math.Abs(rsis[0]-0) > 1e-9 {
		t.Errorf("expected rsi 0, got %f", rsis[0])
	}
}

func TestCalcRSISeriesFlat(t *testing.T) {
	candles := candlesFromCloses(5, 5, 5, 5, 5, 5)

	rsis, err := CalcRSISeries(candles, 3)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	// A flat series has no direction: neutral RSI 50.
	if math.Abs(rsis[0]-50) > 1e-9 {
		t.Errorf("expected rsi 50, got %f", rsis[0])
	}
}

func TestCalcRSISeriesBalancedMoves(t *testing.T) {
	// Alternating equal gains and losses: avgGain converges to avgLoss, RSI near 50.
	candles := candlesFromCloses(10, 11, 10, 11, 10, 11, 10, 11, 10, 11, 10, 11)

	rsis, err := CalcRSISeries(candles, 4)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if rsis[0] < 30 || rsis[0] > 70 {
		t.Errorf("expected rsi around 50 for balanced moves, got %f", rsis[0])
	}
}

func TestCalcRSISeriesNotEnoughCandles(t *testing.T) {
	candles := candlesFromCloses(1, 2, 3)

	_, err := CalcRSISeries(candles, 3)
	if err == nil {
		t.Error("expected error for not enough candles, got nil")
	}
}

func TestCalcRSISeriesUndefinedPositionsAreZero(t *testing.T) {
	candles := candlesFromCloses(1, 2, 3, 4, 5, 6)

	rsis, err := CalcRSISeries(candles, 3)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	// The oldest `period` positions have no RSI.
	for i := len(rsis) - 3; i < len(rsis); i++ {
		if rsis[i] != 0 {
			t.Errorf("expected 0 at undefined position %d, got %f", i, rsis[i])
		}
	}
}

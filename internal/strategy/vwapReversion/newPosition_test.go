package vwapReversion

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"testing"
)

// flatVolumeCandles builds candles where typical price == close so the VWAP and
// sigma are easy to reason about in the assertions below.
func flatVolumeCandles(closes []float64) []domainStructs.DomainCandle {
	candles := make([]domainStructs.DomainCandle, len(closes))
	for i, c := range closes {
		candles[i] = domainStructs.DomainCandle{
			Time:   int64(i),
			High:   c,
			Low:    c,
			Open:   c,
			Close:  c,
			Volume: 1,
		}
	}
	return candles
}

func TestCalculateVwapChartLong(t *testing.T) {
	v := VwapReversion{}
	v.config = Config{
		EntrySigmaMult: 2,
		SlSigmaMult:    3,
		PricePrecision: 4,
	}

	// Equal volume, symmetric prices around 100: VWAP = 100, sigma = sqrt(2).
	chart := v.calculateVwapChart(flatVolumeCandles([]float64{98, 99, 100, 101, 102}))

	if chart.Vwap != 100 {
		t.Errorf("Vwap = %f; expects 100", chart.Vwap)
	}
	// sigma = sqrt( (4+1+0+1+4)/5 ) = sqrt(2) = 1.4142...
	expectedSigma := 1.4142
	if chart.StdDev != expectedSigma {
		t.Errorf("StdDev = %f; expects %f", chart.StdDev, expectedSigma)
	}
	// lower entry band = VWAP - 2*sigma = 100 - 2.8284 = 97.1716
	expectedLower := 97.1716
	if chart.LowerBand != expectedLower {
		t.Errorf("LowerBand = %f; expects %f", chart.LowerBand, expectedLower)
	}
	if chart.TakeProfit != 100 {
		t.Errorf("TakeProfit = %f; expects 100 (VWAP)", chart.TakeProfit)
	}

	entry, tp, sl := v.entryTpSlForSide(chart, domainStructs.PositionSideLong)
	if entry != expectedLower {
		t.Errorf("long entry = %f; expects %f", entry, expectedLower)
	}
	if tp != 100 {
		t.Errorf("long tp = %f; expects 100", tp)
	}
	// SL = VWAP - 3*sigma = 100 - 4.2426 = 95.7574
	expectedSl := 95.7574
	if sl != expectedSl {
		t.Errorf("long sl = %f; expects %f", sl, expectedSl)
	}
}

func TestEntryTpSlForSideShort(t *testing.T) {
	v := VwapReversion{}
	v.config = Config{
		EntrySigmaMult: 2,
		SlSigmaMult:    3,
		PricePrecision: 4,
	}
	chart := v.calculateVwapChart(flatVolumeCandles([]float64{98, 99, 100, 101, 102}))

	entry, tp, sl := v.entryTpSlForSide(chart, domainStructs.PositionSideShort)
	// upper entry band = VWAP + 2*sigma = 102.8284
	expectedUpper := 102.8284
	if entry != expectedUpper {
		t.Errorf("short entry = %f; expects %f", entry, expectedUpper)
	}
	if tp != 100 {
		t.Errorf("short tp = %f; expects 100", tp)
	}
	// SL = VWAP + 3*sigma = 104.2426
	expectedSl := 104.2426
	if sl != expectedSl {
		t.Errorf("short sl = %f; expects %f", sl, expectedSl)
	}
}

func TestCalculateQty(t *testing.T) {
	v := VwapReversion{}
	v.config = Config{
		RiskPercent:  10,
		QtyPrecision: 4,
	}

	// balance 500, risk 10% = 50. Distance |80000 - 70000| = 10000.
	// qty = 50 / 10000 = 0.005
	qty, err := v.calculateQty(500.0, 80000, 70000)
	if err != nil {
		t.Errorf("Error calculateQty: %s", err.Error())
	}
	expected := 0.005
	if qty != expected {
		t.Errorf("calculateQty = %f; expects %f", qty, expected)
	}
}

func TestCalculateQtyZeroBalance(t *testing.T) {
	v := VwapReversion{}
	v.config = Config{RiskPercent: 10, QtyPrecision: 4}

	_, err := v.calculateQty(0, 80000, 70000)
	if err == nil {
		t.Errorf("calculateQty with zero balance should return an error")
	}
}

func TestCalculateQtyZeroDistance(t *testing.T) {
	v := VwapReversion{}
	v.config = Config{RiskPercent: 10, QtyPrecision: 4}

	_, err := v.calculateQty(500, 80000, 80000)
	if err == nil {
		t.Errorf("calculateQty with zero price distance should return an error")
	}
}

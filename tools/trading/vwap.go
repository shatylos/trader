package trading

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/math"
)

type VWAP struct {
	Candles []domainStructs.DomainCandle
	AncVWAP float64
}

func CreateVWAP(candles []domainStructs.DomainCandle) (vwap VWAP) {
	vwap.Candles = candles
	var pvSum float64
	var volSum float64
	for _, c := range candles {
		price := math.Div(c.High+c.Low+c.Close, 3.0)
		pvSum += math.Mul(price, c.Volume)
		volSum += c.Volume
	}

	if volSum == 0 {
		vwap.AncVWAP = 0
		return
	}

	vwap.AncVWAP = math.Div(pvSum, volSum)
	return
}

func (v *VWAP) CalcDeviation(deviation float64) (upperBand float64, lowerBand float64) {
	var varianceSum float64
	var volSum float64

	for _, c := range v.Candles {
		// V * (P - VWAP)^2
		price := math.Div(c.High+c.Low+c.Close, 3.0)
		diff := price - v.AncVWAP
		varianceSum += math.Mul(c.Volume, math.Mul(diff, diff))
		volSum += c.Volume
	}

	if volSum == 0 {
		upperBand = v.AncVWAP
		lowerBand = v.AncVWAP
		return
	}

	// σ = sqrt( Σ(V * (P - VWAP)^2) / ΣV )
	variance := math.Div(varianceSum, volSum)
	stdDev := math.Sqrt(variance)

	// Bands
	offset := deviation * stdDev
	upperBand = v.AncVWAP + offset
	lowerBand = v.AncVWAP - offset

	return
}

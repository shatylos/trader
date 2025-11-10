package _struct

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	"time"
)

type HeapTimeframe struct {
	Config     HeapConfig
	CandleTime time.Time
	Candles    []domainStructs.DomainCandle
	HeapStatus entity.HeapStatus
	prevValues struct {
		IsSidewaysState bool
		Zone            string
		Trend           string
		TrendSlope      float64
		HeapStatus      entity.HeapStatus
	}
}

type HeapConfig struct {
	Resolution                  string
	QtyPercent                  float64
	CandleReview                int64
	CandleCacheDuration         time.Duration
	SidewaysMinCandlesAmount    int64
	MaxNumberOrdersToBuy        int64
	SidewaysPercentToPrice      float64
	SidewaysPremiumCoefficient  float64
	SidewaysDiscountCoefficient float64
	MinPercentRangeToSell       float64
	MinPercentRangeToBuyMore    float64
	DurationToMoveToHeap        time.Duration
	QtyPercentOnMaxPrice        float64
	QtyPercentOnMinPrice        float64
}

func (t *HeapTimeframe) GetConfig() TimeframeConfig {
	return t.Config
}

func (t *HeapTimeframe) Resolution() string {
	return t.Config.Resolution
}

func (t *HeapTimeframe) SetCandles(candles []domainStructs.DomainCandle) {
	t.Candles = candles
	t.CandleTime = time.Now()
}

func (t *HeapTimeframe) GetCandles() []domainStructs.DomainCandle {
	return t.Candles
}

func (t *HeapTimeframe) GetCandleTime() time.Time {
	return t.CandleTime
}

func (t *HeapTimeframe) MinPercentRangeToSell() float64 {
	return t.Config.MinPercentRangeToSell
}

func (t *HeapTimeframe) IsHeap() bool {
	return true
}

func (t HeapConfig) GetCandleReview() int64 {
	return t.CandleReview
}

func (t HeapConfig) GetCandleCacheDuration() time.Duration {
	return t.CandleCacheDuration
}

func (t *HeapTimeframe) IsStatusChanged() (isChanged bool) {
	if t.HeapStatus.IsSidewaysState != t.prevValues.IsSidewaysState {
		isChanged = true
		t.prevValues.IsSidewaysState = t.HeapStatus.IsSidewaysState
	}
	if t.HeapStatus.Zone != t.prevValues.Zone {
		isChanged = true
		t.prevValues.Zone = t.HeapStatus.Zone
	}
	if t.HeapStatus.Trend != t.prevValues.Trend {
		isChanged = true
		t.prevValues.Trend = t.HeapStatus.Trend
	}
	if t.HeapStatus.TrendSlope != t.prevValues.TrendSlope {
		isChanged = true
		t.prevValues.TrendSlope = t.HeapStatus.TrendSlope
	}
	if t.HeapStatus.Price != t.prevValues.HeapStatus.Price {
		isChanged = true
		t.prevValues.HeapStatus.Price = t.HeapStatus.Price
	}
	if t.HeapStatus.Qty != t.prevValues.HeapStatus.Qty {
		isChanged = true
		t.prevValues.HeapStatus.Qty = t.HeapStatus.Qty
	}
	return
}

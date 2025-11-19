package _struct

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	"time"
)

type HeapTimeframe struct {
	Config         HeapConfig
	CandleTime     time.Time
	Candles        []domainStructs.DomainCandle
	HeapStatus     entity.HeapStatus
	prevHeapStatus entity.HeapStatus
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

func (t *HeapTimeframe) GetTrendSlope() float64 {
	return t.HeapStatus.TrendSlope
}

func (t HeapConfig) GetCandleReview() int64 {
	return t.CandleReview
}

func (t HeapConfig) GetCandleCacheDuration() time.Duration {
	return t.CandleCacheDuration
}

func (t *HeapTimeframe) IsStatusChanged() (isChanged bool) {
	if t.HeapStatus.IsSidewaysState != t.prevHeapStatus.IsSidewaysState {
		isChanged = true
		t.prevHeapStatus.IsSidewaysState = t.HeapStatus.IsSidewaysState
	}
	if t.HeapStatus.Zone != t.prevHeapStatus.Zone {
		isChanged = true
		t.prevHeapStatus.Zone = t.HeapStatus.Zone
	}
	if t.HeapStatus.Trend != t.prevHeapStatus.Trend {
		isChanged = true
		t.prevHeapStatus.Trend = t.HeapStatus.Trend
	}
	if t.HeapStatus.TrendSlope != t.prevHeapStatus.TrendSlope {
		isChanged = true
		t.prevHeapStatus.TrendSlope = t.HeapStatus.TrendSlope
	}
	if t.HeapStatus.Price != t.prevHeapStatus.Price {
		isChanged = true
		t.prevHeapStatus.Price = t.HeapStatus.Price
	}
	if t.HeapStatus.Qty != t.prevHeapStatus.Qty {
		isChanged = true
		t.prevHeapStatus.Qty = t.HeapStatus.Qty
	}
	return
}

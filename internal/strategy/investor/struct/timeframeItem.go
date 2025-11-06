package _struct

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"time"
)

type TimeframeItem struct {
	Config          TimeframeItemConfig
	CandleTime      time.Time
	Candles         []domainStructs.DomainCandle
	IsSidewaysState bool
	Zone            string
	Trend           string
	TrendSlope      float64
	prevValues      struct {
		IsSidewaysState bool
		Zone            string
		Trend           string
		TrendSlope      float64
	}
}

type TimeframeItemConfig struct {
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
}

func (t *TimeframeItem) GetConfig() TimeframeConfig {
	return t.Config
}

func (t *TimeframeItem) Resolution() string {
	return t.Config.Resolution
}

func (t *TimeframeItem) SetCandles(candles []domainStructs.DomainCandle) {
	t.Candles = candles
	t.CandleTime = time.Now()
}

func (t *TimeframeItem) GetCandles() []domainStructs.DomainCandle {
	return t.Candles
}

func (t *TimeframeItem) GetCandleTime() time.Time {
	return t.CandleTime
}

func (t *TimeframeItem) MinPercentRangeToSell() float64 {
	return t.Config.MinPercentRangeToSell
}

func (t *TimeframeItem) IsHeap() bool {
	return false
}

func (t TimeframeItemConfig) GetCandleReview() int64 {
	return t.CandleReview
}

func (t TimeframeItemConfig) GetCandleCacheDuration() time.Duration {
	return t.CandleCacheDuration
}

func (t *TimeframeItem) IsStatusChanged() (isChanged bool) {
	if t.IsSidewaysState != t.prevValues.IsSidewaysState {
		isChanged = true
		t.prevValues.IsSidewaysState = t.IsSidewaysState
	}
	if t.Zone != t.prevValues.Zone {
		isChanged = true
		t.prevValues.Zone = t.Zone
	}
	if t.Trend != t.prevValues.Trend {
		isChanged = true
		t.prevValues.Trend = t.Trend
	}
	if t.TrendSlope != t.prevValues.TrendSlope {
		isChanged = true
		t.prevValues.TrendSlope = t.TrendSlope
	}
	return
}

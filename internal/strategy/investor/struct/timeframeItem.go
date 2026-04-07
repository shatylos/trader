package _struct

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"time"
)

type TimeframeItem struct {
	Config          TimeframeItemConfig
	CandleTime      time.Time
	Candles         []domainStructs.DomainCandle
	TrendSlope      float64
	HigherTimeframe Timeframe
	IsSidewaysState bool
	SidewaysFrom    time.Time
	Zone            string
	Trend           string
	TradeStateMsg   string
	prevValues      struct {
		IsSidewaysState bool
		SidewaysFrom    time.Time
		Zone            string
		Trend           string
		TrendSlope      float64
		TradeStateMsg   string
	}
}

type TimeframeItemConfig struct {
	Resolution               string
	CanOpenNewOrder          bool
	QtyPercent               float64
	CandleReview             int64
	CandleCacheDuration      time.Duration
	SidewaysMinCandlesAmount int64
	SidewaysPercentToPrice   float64
	DurationToMoveToHeap     time.Duration
	VwapDeviationsBuy        []float64
	VwapDeviationsSell       []float64
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

func (t *TimeframeItem) IsHeap() bool {
	return false
}

func (t *TimeframeItem) GetTrendSlope() float64 {
	return t.TrendSlope
}

func (t *TimeframeItem) GetHigherTrendSlope() float64 {
	if t.HigherTimeframe == nil {
		return 0
	}
	return t.HigherTimeframe.GetTrendSlope()
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
	if !t.SidewaysFrom.Equal(t.prevValues.SidewaysFrom) {
		isChanged = true
		t.prevValues.SidewaysFrom = t.SidewaysFrom
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
	if t.TradeStateMsg != t.prevValues.TradeStateMsg {
		isChanged = true
		t.prevValues.TradeStateMsg = t.TradeStateMsg
	}
	return
}

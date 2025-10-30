package _struct

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"time"
)

type Timeframe struct {
	Config          TimeframeConfig
	GetCandleTime   time.Time
	Candles         []domainStructs.DomainCandle
	IsSidewaysState bool
	Zone            string
	prevValues      struct {
		IsSidewaysState bool
		Zone            string
	}
}
type TimeframeConfig struct {
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
	IsHeap                      bool
	HeapConfig                  *HeapConfig
	DurationToMoveToHeap        time.Duration
}

type HeapConfig struct {
	QtyPercentOnMaxPrice float64
	QtyPercentOnMinPrice float64
}

func (t *Timeframe) IsStatusChanged() (isChanged bool) {
	if t.IsSidewaysState != t.prevValues.IsSidewaysState {
		isChanged = true
		t.prevValues.IsSidewaysState = t.IsSidewaysState
	}
	if t.Zone != t.prevValues.Zone {
		isChanged = true
		t.prevValues.Zone = t.Zone
	}
	return
}

package _struct

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"time"
)

type Timeframe struct {
	Config        TimeframeConfig
	GetCandleTime time.Time
	Candles       []domainStructs.DomainCandle
}
type TimeframeConfig struct {
	Resolution                  string
	QtyPercent                  float64
	CandleReview                int64
	CandleCacheSeconds          int64
	SidewaysMinCandlesAmount    int64
	MaxNumberOrdersToBuy        int64
	SidewaysPercentToPrice      float64
	SidewaysPremiumCoefficient  float64
	SidewaysDiscountCoefficient float64
	MinPercentRangeToSell       float64
	MinPercentRangeToBuyMore    float64
	IsHeap                      bool
	HeapConfig                  *HeapConfig
}

type HeapConfig struct {
	QtyPercentOnMaxPrice float64
	QtyPercentOnMinPrice float64
}

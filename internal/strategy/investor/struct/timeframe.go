package _struct

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"time"
)

type Timeframe interface {
	GetConfig() TimeframeConfig
	Resolution() string
	SetCandles(candles []domainStructs.DomainCandle)
	GetCandles() []domainStructs.DomainCandle
	GetCandleTime() time.Time
	MinPercentRangeToSell() float64
	IsHeap() bool
}

type TimeframeConfig interface {
	GetCandleReview() int64
	GetCandleCacheDuration() time.Duration
}

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
	GetTrendSlope() float64
	GetCandleTime() time.Time
	IsHeap() bool
}

type TimeframeConfig interface {
	GetCandleReview() int64
	GetCandleCacheDuration() time.Duration
}

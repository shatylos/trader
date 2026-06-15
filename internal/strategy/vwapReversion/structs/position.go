package structs

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"time"
)

type Position struct {
	Id               *string                      `bson:"_id,omitempty"`
	Chart            VwapChart                    `bson:"Chart"`
	Trend            string                       `bson:"Trend"`
	LtTrend          string                       `bson:"LtTrend"`
	StTrend          string                       `bson:"StTrend"`
	Order            domainStructs.DomainOrder    `bson:"Order"`
	CreatedTime      time.Time                    `bson:"CreatedTime"`
	UpdatedTime      time.Time                    `bson:"UpdatedTime"`
	ClosedTime       time.Time                    `bson:"ClosedTime"`
	Status           string                       `bson:"Status"`
	ProviderPosition domainStructs.DomainPosition `bson:"ProviderPosition"`
	BalanceBefore    float64                      `bson:"BalanceBefore"`
	BalanceAfter     float64                      `bson:"BalanceAfter"`
	TotalClosePnl    float64                      `bson:"TotalClosePnl"`
}

// VwapChart holds the volume-weighted average price reference together with the
// deviation bands and the resulting entry / take-profit / stop-loss levels for a
// single mean-reversion setup.
type VwapChart struct {
	Vwap       float64
	UpperBand  float64
	LowerBand  float64
	StdDev     float64
	EntryPrice float64
	TakeProfit float64
	StopLoss   float64
	Qty        float64
}

const StatusNew = "NEW"
const StatusActive = "ACTIVE"
const StatusClosed = "CLOSED"

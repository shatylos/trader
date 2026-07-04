package structs

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"time"
)

type Position struct {
	Id               *string                      `bson:"_id,omitempty"`
	Chart            ScalpChart                   `bson:"Chart"`
	Bias             string                       `bson:"Bias"`
	LtTrend          string                       `bson:"LtTrend"`
	Signal           string                       `bson:"Signal"`
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

// ScalpChart holds the indicator snapshot the entry decision was made on,
// together with the resulting entry / take-profit / stop-loss levels of a
// single scalp setup.
type ScalpChart struct {
	LtfEma     float64
	HtfEmaFast float64
	HtfEmaSlow float64
	Rsi        float64
	Atr        float64
	EntryPrice float64
	TakeProfit float64
	StopLoss   float64
	Qty        float64
}

const StatusNew = "NEW"
const StatusActive = "ACTIVE"
const StatusClosed = "CLOSED"

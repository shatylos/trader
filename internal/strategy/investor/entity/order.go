package entity

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/math"
	"time"
)

type Order struct {
	domainStructs.DomainOrder `bson:",inline"`
	Timeframe                 string                     `bson:"Timeframe"`
	CreatedTime               time.Time                  `bson:"CreatedTime"`
	UpdatedTime               time.Time                  `bson:"UpdatedTime"`
	WalletBefore              domainStructs.DomainWallet `bson:"WalletBefore"`
	WalletAfter               domainStructs.DomainWallet `bson:"WalletAfter"`
	ConfigKey                 int                        `bson:"ConfigKey"`
	AverageBuyPrice           float64                    `bson:"AverageBuyPrice"`
	QtyInTrade                float64                    `bson:"QtyInTrade"`
	RealizedPNL               float64                    `bson:"RealizedPNL"`
	StateApplied              bool                       `bson:"StateApplied"`
	Moved                     bool                       `bson:"Moved"`
}

func (o *Order) Amount() float64 {
	return math.Mul(o.Qty, o.Price)
}

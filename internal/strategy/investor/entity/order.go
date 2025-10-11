package entity

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/math"
	"time"
)

type Order struct {
	domainStructs.DomainOrder `bson:",inline"`
	DealId                    string                     `bson:"DealId"`
	Timeframe                 string                     `bson:"Timeframe"`
	CreatedTime               time.Time                  `bson:"CreatedTime"`
	UpdatedTime               time.Time                  `bson:"UpdatedTime"`
	WalletBefore              domainStructs.DomainWallet `bson:"WalletBefore"`
	WalletAfter               domainStructs.DomainWallet `bson:"WalletAfter"`
}

func (o *Order) Amount() float64 {
	return math.Mul(o.Qty, o.Price)
}

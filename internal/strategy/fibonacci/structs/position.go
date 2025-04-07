package structs

import (
	"github.com/shatylos/trader/internal/domain/structs"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
)

type Position struct {
	Id               *string                      `bson:"_id,omitempty"`
	FibonacciChart   FibonacciChart               `bson:"FibonacciChart"`
	Trend            string                       `bson:"Trend"`
	Orders           PositionOrders               `bson:"Orders"`
	CreatedTime      int64                        `bson:"CreatedTime"`
	UpdatedTime      int64                        `bson:"UpdatedTime"`
	Status           string                       `bson:"Status"`
	ProviderPosition domainStructs.DomainPosition `bson:"ProviderPosition"`
}
type FibonacciChart struct {
	EntryPoint1    float64
	EntryPoint2    float64
	EntryPoint3    float64
	SourceMaxPrice float64
	SourceMinPrice float64
	StopLoss       float64
	TakeProfit1    float64
	TakeProfit2    float64
	TakeProfit3    float64
	FullQty        float64
}

type PositionOrders struct {
	Order1 structs.DomainOrder
	Order2 structs.DomainOrder
	Order3 structs.DomainOrder
}

const StatusNew = "NEW"
const StatusActive = "ACTIVE"
const StatusClosed = "CLOSED"

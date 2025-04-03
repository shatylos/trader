package structs

type DomainOrder struct {
	Id          *string `bson:"_id,omitempty"`
	CreatedTime int64   `bson:"CreatedTime"`
	OrderId     string  `bson:"OrderId"`
	OrderStatus string  `bson:"OrderStatus"`
	OrderType   string  `bson:"OrderType"`
	Price       float64 `bson:"Price"`
	Qty         float64 `bson:"Qty"`
	ReduceOnly  bool    `bson:"ReduceOnly"`
	Side        string  `bson:"Side"`
	Symbol      string  `bson:"Symbol"`
	TimeInForce string  `bson:"TimeInForce"`
	UpdatedTime int64   `bson:"UpdatedTime"`
	TakeProfit  int64   `bson:"TakeProfit"`
	StopLoss    int64   `bson:"StopLoss"`
}

const OrderStatusOpen = "Open"
const OrderStatusFilled = "Filled"
const OrderStatusCancelled = "Cancelled"

package structs

type DomainOrder struct {
	Id           *string `bson:"_id,omitempty"`
	CreatedTime  int64   `bson:"CreatedTime"`
	OrderId      string  `bson:"OrderId"`
	OrderStatus  string  `bson:"OrderStatus"`
	OrderType    string  `bson:"OrderType"`
	Price        float64 `bson:"Price"`
	Qty          float64 `bson:"Qty"`
	ReduceOnly   bool    `bson:"ReduceOnly"`
	Side         string  `bson:"Side"`
	Symbol       string  `bson:"Symbol"`
	TimeInForce  string  `bson:"TimeInForce"`
	UpdatedTime  int64   `bson:"UpdatedTime"`
	TpModifyTime int64   `bson:"TpModifyTime"`
}

type OrderStatusesStruct struct {
	New             string
	Open            string
	Filled          string
	PartiallyFilled string
	Canceled        string
}

var OrderStatuses = OrderStatusesStruct{
	New:             "NEW",
	Open:            "OPEN",
	Filled:          "FILLED",
	PartiallyFilled: "PARTIALLY_FILLED",
	Canceled:        "CANCELED",
}

type OrderTypesStruct struct {
	Limit  string
	Market string
}

var OrderTypes = OrderTypesStruct{
	Limit:  "LIMIT",
	Market: "MARKET",
}

const OrderSideBuy = "BUY"
const OrderSideSell = "SELL"

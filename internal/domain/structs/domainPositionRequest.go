package structs

type DomainPositionRequest struct {
	Leverage   int64
	Price      float64
	Qty        float64
	ReduceOnly bool
	Side       string
	StopLoss   float64
	Symbol     string
	TakeProfit float64
	Type       string
	// TpOrderType makes the attached take profit execute as this order type
	// when set to PositionTypes.Limit (maker fee); empty means market order.
	TpOrderType string
}

type TpSlRequest struct {
	CoinPare   string
	TakeProfit float64
	StopLoss   float64
	Side       string
}

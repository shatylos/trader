package structs

type DomainPositionRequest struct {
	Leverage    int64
	PositionId  string
	Price       float64
	Qty         float64
	ReduceOnly  bool
	Side        string
	StopLoss    float64
	Symbol      string
	TakeProfit  float64
	TimeInForce string
	Type        string
}

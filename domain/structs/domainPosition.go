package structs

type DomainPosition struct {
	Leverage      int64
	PositionId    string
	Price         float64
	Qty           float64
	RealizedPnl   float64
	Side          string
	StopLoss      float64
	Symbol        string
	TakeProfit    float64
	Type          string
	UnrealizedPnl float64
}

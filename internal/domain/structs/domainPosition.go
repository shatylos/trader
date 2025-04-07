package structs

type DomainPosition struct {
	Leverage      int64
	PositionId    string
	AvgPrice      float64
	MarkPrice     float64
	Size          float64
	RealizedPnl   float64
	TotalPnl      float64
	Side          string
	StopLoss      float64
	Symbol        string
	TakeProfit    float64
	Type          string
	UnrealizedPnl float64
}

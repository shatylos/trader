package structs

type DomainPosition struct {
	EntryPrice       float64
	Leverage         int64
	LiquidationPrice float64
	Margin           float64
	Pair             string
	Quantity         float64
	RealizedPnl      float64
	StopLoss         float64
	TakeProfit       float64
	Type             string
	UnrealizedPnl    float64
	Value            float64
}

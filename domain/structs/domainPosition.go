package structs

type DomainPosition struct {
	Margin          float64
	Type            string
	Quantity        float64
	Leverage        int64
	BasePrice       float64
	Pair            string
	FundingQuantity float64
	TakeProfit      float64
	StopLoss        float64
	MarginCurrency  string
	Roe             float64
	PositionId      int64
	UnrealizedPnl   float64
	RealizedPnl     float64
	PnlCurrency     string
	FundingCurrency string
}

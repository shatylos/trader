package structs

type DomainPosition struct {
	Leverage   int64
	AvgPrice   float64
	MarkPrice  float64
	Size       float64
	TotalPnl   float64
	Side       string
	StopLoss   float64
	Symbol     string
	TakeProfit float64
}

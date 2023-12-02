package structs

type HistoryOrder struct {
	Id                        int64
	DomainOrderId             string
	FilledPrice               float64
	FilledQty                 float64
	Side                      string
	CreatedTime               int64
	UpdatedTime               int64
	MainCurrencyAmountBefore  float64
	TradeCurrencyAmountBefore float64
	AveragePrice              float64
}

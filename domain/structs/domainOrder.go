package structs

type DomainOrder struct {
	CreatedTime int64
	OrderId     string
	OrderStatus string
	OrderType   string
	Price       float64
	Qty         float64
	ReduceOnly  bool
	Side        string
	Symbol      string
	TimeInForce string
	UpdatedTime string
}

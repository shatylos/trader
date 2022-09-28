package structs

type DomainOrder struct {
	CreatedTime string
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

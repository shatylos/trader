package structs

type DomainOrderRequest struct {
	OrderId     string
	Price       float64
	Qty         float64
	ReduceOnly  bool
	Side        string
	Symbol      string
	TimeInForce string
	Type        string
}

package entity

type HeapStatus struct {
	Qty             float64
	Price           float64
	Deal            *Deal
	LastOrderHeap   *Order
	LastOrderMoved  *Order
	PremiumDiscount float64
	IsSidewaysState bool
	Zone            string
	Trend           string
	TrendSlope      float64
}

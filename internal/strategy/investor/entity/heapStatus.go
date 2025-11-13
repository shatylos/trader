package entity

type HeapStatus struct {
	Qty             float64
	PurposeQty      float64
	PurposeQtyEqv   float64
	QtyExcess       float64
	QtyExcessEqv    float64
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

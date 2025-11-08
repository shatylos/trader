package entity

type HeapStatus struct {
	Qty            float64
	Price          float64
	Deal           *Deal
	LastOrderHeap  *Order
	LastOrderMoved *Order
}

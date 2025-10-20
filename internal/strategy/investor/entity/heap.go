package entity

type Heap struct {
	Qty            float64
	Price          float64
	Deal           *Deal
	LastOrderHeap  *Order
	LastOrderMoved *Order
}

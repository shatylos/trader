package entity

type Heap struct {
	Qty            float64
	Price          float64
	LastOrderHeap  *Order
	LastOrderMoved *Order
}

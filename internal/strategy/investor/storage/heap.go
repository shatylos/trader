package storage

import (
	"context"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/math"
)

func (s *Storage) UpdateHeapStatus(ctx context.Context, heapTimeFrame *_struct.HeapTimeframe) (err error) {

	var dealRelations []*entity.DealRelation
	dealRelations, err = s.GetDealRelationsOnHeap(ctx)
	if err != nil {
		return
	}

	buyQty := 0.0
	buyAmount := 0.0
	sellQty := 0.0
	sellAmount := 0.0

	var lastOrderHeap *entity.Order
	var lastOrderMoved *entity.Order

	for _, dealRelation := range dealRelations {
		for _, order := range dealRelation.Orders {
			if order.Timeframe == heapTimeFrame.Config.Resolution {
				if lastOrderHeap == nil || order.CreatedTime.After(lastOrderHeap.CreatedTime) {
					lastOrderHeap = order
				}
			} else {
				if lastOrderMoved == nil || order.CreatedTime.After(lastOrderMoved.CreatedTime) {
					lastOrderMoved = order
				}
			}
			if order.Side == structs.OrderSideBuy {
				buyQty += buyQty
				buyAmount += order.Amount()
			}
			if order.Side == structs.OrderSideSell {
				sellQty += sellQty
				sellAmount += order.Amount()
			}
		}
	}

	qty := buyQty - sellQty
	price := 0.0
	if buyQty > 0 {
		price = math.Div(buyAmount, buyQty)
	}

	var deal *entity.Deal
	deal, err = GetActiveDealByTimeframe(ctx, s, heapTimeFrame)
	if err != nil {
		return
	}

	heapTimeFrame.HeapStatus.Qty = qty
	heapTimeFrame.HeapStatus.Price = price
	heapTimeFrame.HeapStatus.Deal = deal
	heapTimeFrame.HeapStatus.LastOrderHeap = lastOrderHeap
	heapTimeFrame.HeapStatus.LastOrderMoved = lastOrderMoved

	return
}

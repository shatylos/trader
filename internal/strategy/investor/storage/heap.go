package storage

import (
	"context"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
)

func (s *Storage) UpdateHeap(ctx context.Context, heapParam *entity.Heap, timeFrameItem *_struct.Timeframe) (heap *entity.Heap, err error) {
	heap = heapParam
	if heap == nil {
		heap, err = s.getHeap(ctx, timeFrameItem)
		if err != nil {
			return
		}
	}
	return
}

func (s *Storage) getHeap(ctx context.Context, timeFrameItem *_struct.Timeframe) (heapPointer *entity.Heap, err error) {

	heapTimeframe, ok := ctx.Value(_struct.CtxHeapTimeframeKey).(string)
	if !ok {
		msg := "HeapTimeframe is not accessible from context"
		logger.Error(msg)
		err = tools.AppError{Message: msg}
		return
	}

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
			if order.Timeframe == heapTimeframe {
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
	deal, err = s.GetActiveDealByTimeframe(ctx, timeFrameItem)
	if err != nil {
		return
	}

	heap := entity.Heap{
		Qty:            qty,
		Price:          price,
		Deal:           deal,
		LastOrderHeap:  lastOrderHeap,
		LastOrderMoved: lastOrderMoved,
	}

	heapPointer = &heap

	return
}

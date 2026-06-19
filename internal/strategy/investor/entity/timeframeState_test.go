package entity

import (
	"github.com/shatylos/trader/internal/domain/structs"
	"testing"
)

func makeFilledOrder(side string, qty, price float64) *Order {
	order := &Order{Timeframe: "1h"}
	order.Side = side
	order.OrderStatus = structs.OrderStatuses.Filled
	order.Qty = qty
	order.Price = price
	return order
}

func TestApplyFilledOrderBuySellCycle(t *testing.T) {
	state := &TimeframeState{Timeframe: "1h"}
	minQty := 0.0001

	buy1 := makeFilledOrder(structs.OrderSideBuy, 1, 100)
	state.ApplyFilledOrder(buy1, minQty)
	if state.QtyInTrade != 1 || state.AverageBuyPrice != 100 {
		t.Fatalf("after first buy: qty %g avg %g", state.QtyInTrade, state.AverageBuyPrice)
	}

	buy2 := makeFilledOrder(structs.OrderSideBuy, 1, 200)
	state.ApplyFilledOrder(buy2, minQty)
	if state.QtyInTrade != 2 || state.AverageBuyPrice != 150 {
		t.Fatalf("after second buy: qty %g avg %g", state.QtyInTrade, state.AverageBuyPrice)
	}
	if buy2.QtyInTrade != 2 || buy2.AverageBuyPrice != 150 || !buy2.StateApplied {
		t.Fatalf("buy order stored values: qty %g avg %g applied %v", buy2.QtyInTrade, buy2.AverageBuyPrice, buy2.StateApplied)
	}

	sell1 := makeFilledOrder(structs.OrderSideSell, 1, 180)
	state.ApplyFilledOrder(sell1, minQty)
	if sell1.RealizedPNL != 30 {
		t.Fatalf("sell realized PNL: %g", sell1.RealizedPNL)
	}
	if state.QtyInTrade != 1 || state.AverageBuyPrice != 150 || state.RealizedPNL != 30 {
		t.Fatalf("after first sell: qty %g avg %g realized %g", state.QtyInTrade, state.AverageBuyPrice, state.RealizedPNL)
	}

	state.LastBuyOrder = buy2
	state.LastSellOrder = sell1

	// the second sell closes the cycle
	sell2 := makeFilledOrder(structs.OrderSideSell, 1, 160)
	state.ApplyFilledOrder(sell2, minQty)
	if state.QtyInTrade != 0 || state.AverageBuyPrice != 0 || state.RealizedPNL != 0 {
		t.Fatalf("after cycle close: qty %g avg %g realized %g", state.QtyInTrade, state.AverageBuyPrice, state.RealizedPNL)
	}
	if state.LastBuyOrder != nil || state.LastSellOrder != nil {
		t.Fatalf("after cycle close last orders must be nil")
	}
	if sell2.QtyInTrade != 0 || sell2.AverageBuyPrice != 0 {
		t.Fatalf("closing sell stored values: qty %g avg %g", sell2.QtyInTrade, sell2.AverageBuyPrice)
	}
	if sell2.RealizedPNL != 10 {
		t.Fatalf("closing sell realized PNL: %g", sell2.RealizedPNL)
	}
}

func TestApplyMovedSellDefersLossWithZeroPNL(t *testing.T) {
	minQty := 0.0001

	// child holds 2 coins bought at avg 200
	child := &TimeframeState{Timeframe: "30m"}
	child.ApplyFilledOrder(makeFilledOrder(structs.OrderSideBuy, 2, 200), minQty)

	// move 1 coin to parent at a price below avg buy (would be a loss if sold)
	movedSell := makeFilledOrder(structs.OrderSideSell, 1, 180)
	movedSell.Moved = true
	child.ApplyMovedSell(movedSell, minQty)

	if movedSell.RealizedPNL != 0 {
		t.Fatalf("moved sell must have 0 PNL, got %g", movedSell.RealizedPNL)
	}
	if child.RealizedPNL != 0 {
		t.Fatalf("child realized PNL must stay 0 after move, got %g", child.RealizedPNL)
	}
	if child.QtyInTrade != 1 || child.AverageBuyPrice != 200 {
		t.Fatalf("child after move: qty %g avg %g", child.QtyInTrade, child.AverageBuyPrice)
	}
	if !movedSell.StateApplied || movedSell.QtyInTrade != 1 || movedSell.AverageBuyPrice != 200 {
		t.Fatalf("moved sell stored values: applied %v qty %g avg %g", movedSell.StateApplied, movedSell.QtyInTrade, movedSell.AverageBuyPrice)
	}

	// parent receives the moved qty as a normal buy and later books the real loss
	parent := &TimeframeState{Timeframe: "4h"}
	movedBuy := makeFilledOrder(structs.OrderSideBuy, 1, 180)
	movedBuy.Moved = true
	parent.ApplyFilledOrder(movedBuy, minQty)
	if parent.QtyInTrade != 1 || parent.AverageBuyPrice != 180 {
		t.Fatalf("parent after moved buy: qty %g avg %g", parent.QtyInTrade, parent.AverageBuyPrice)
	}

	// parent sells lower: the deferred loss surfaces here
	parentSell := makeFilledOrder(structs.OrderSideSell, 1, 170)
	parent.ApplyFilledOrder(parentSell, minQty)
	if parentSell.RealizedPNL != -10 {
		t.Fatalf("parent realized loss expected -10, got %g", parentSell.RealizedPNL)
	}
}

func TestApplyStoredOrderRestoresState(t *testing.T) {
	buy := makeFilledOrder(structs.OrderSideBuy, 2, 150)
	buy.AverageBuyPrice = 150
	buy.QtyInTrade = 2
	buy.StateApplied = true

	sell := makeFilledOrder(structs.OrderSideSell, 1, 180)
	sell.AverageBuyPrice = 150
	sell.QtyInTrade = 1
	sell.RealizedPNL = 30
	sell.StateApplied = true

	state := &TimeframeState{Timeframe: "1h"}
	state.ApplyStoredOrder(buy)
	state.ApplyStoredOrder(sell)

	if state.QtyInTrade != 1 || state.AverageBuyPrice != 150 || state.RealizedPNL != 30 {
		t.Fatalf("restored state: qty %g avg %g realized %g", state.QtyInTrade, state.AverageBuyPrice, state.RealizedPNL)
	}

	// stored zero values mean the cycle was closed by this order
	closingSell := makeFilledOrder(structs.OrderSideSell, 1, 160)
	closingSell.RealizedPNL = 10
	closingSell.StateApplied = true
	state.LastBuyOrder = buy
	state.LastSellOrder = sell
	state.ApplyStoredOrder(closingSell)

	if state.QtyInTrade != 0 || state.AverageBuyPrice != 0 || state.RealizedPNL != 0 {
		t.Fatalf("restored state after close: qty %g avg %g realized %g", state.QtyInTrade, state.AverageBuyPrice, state.RealizedPNL)
	}
	if state.LastBuyOrder != nil || state.LastSellOrder != nil {
		t.Fatalf("after restored cycle close last orders must be nil")
	}
}

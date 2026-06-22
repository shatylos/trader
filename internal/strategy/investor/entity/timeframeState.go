package entity

import (
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/math"
)

// TimeframeState keeps the calculated trading values of one timeframe.
// It is never saved to DB: it is restored from the order history on the
// first run and kept in memory after that.
type TimeframeState struct {
	Timeframe       string
	AverageBuyPrice float64
	QtyInTrade      float64
	PriceToBuy      float64
	PriceToSell     float64
	QtyToBuy        float64
	QtyToSell       float64
	NumOrderToBuy   int
	NumOrderToSell  int
	RealizedPNL     float64
	UnrealizedPNL   float64
	LastBuyOrder    *Order
	LastSellOrder   *Order
	LastFilledOrder *Order
	ActiveOrder     *Order
	IsReset         bool
}

// ApplyFilledOrder updates the state with a filled order and writes the
// resulting values into the order, so the state can be restored from DB.
func (s *TimeframeState) ApplyFilledOrder(order *Order, minQty float64) {
	if order.Side == structs.OrderSideBuy {
		newQty := s.QtyInTrade + order.Qty
		if newQty > 0 {
			spentAmount := math.Mul(s.QtyInTrade, s.AverageBuyPrice) + order.Amount()
			s.AverageBuyPrice = math.Div(spentAmount, newQty)
		}
		s.QtyInTrade = newQty
	}
	if order.Side == structs.OrderSideSell {
		order.RealizedPNL = order.Amount() - math.Mul(order.Qty, s.AverageBuyPrice)
		s.RealizedPNL += order.RealizedPNL
		s.QtyInTrade -= order.Qty
		if s.QtyInTrade <= minQty {
			s.closeCycle()
		}
	}
	order.AverageBuyPrice = s.AverageBuyPrice
	order.QtyInTrade = s.QtyInTrade
	order.StateApplied = true
	s.LastFilledOrder = order
}

// ApplyMovedSell closes the moved qty out of the child cycle with exactly 0
// PNL (the loss is deferred to the parent timeframe, not realized here). The
// resulting state values are written into the order so it restores from DB the
// same way as a live fill.
func (s *TimeframeState) ApplyMovedSell(order *Order, minQty float64) {
	order.RealizedPNL = 0
	s.QtyInTrade -= order.Qty
	if s.QtyInTrade <= minQty {
		s.closeCycle()
	}
	order.AverageBuyPrice = s.AverageBuyPrice
	order.QtyInTrade = s.QtyInTrade
	order.StateApplied = true
	s.LastFilledOrder = order
}

// ApplyStoredOrder restores the state from the values saved in a filled order.
func (s *TimeframeState) ApplyStoredOrder(order *Order) {
	s.AverageBuyPrice = order.AverageBuyPrice
	s.QtyInTrade = order.QtyInTrade
	if order.Side == structs.OrderSideSell {
		s.RealizedPNL += order.RealizedPNL
	}
	if order.QtyInTrade == 0 && order.AverageBuyPrice == 0 { // @TODO: Check condition. Maybe change to less than min to buy
		s.closeCycle()
	}
	s.LastFilledOrder = order
}

// closeCycle drops the buy/sell progression when all coins of the cycle are
// sold (a dust remainder below the min qty is discarded from accounting).
func (s *TimeframeState) closeCycle() {
	s.QtyInTrade = 0
	s.AverageBuyPrice = 0
	s.RealizedPNL = 0
	s.LastBuyOrder = nil
	s.LastSellOrder = nil
}

func (s *TimeframeState) Reset() {
	s.closeCycle()
	s.UnrealizedPNL = 0
	s.PriceToBuy = 0
	s.PriceToSell = 0
	s.QtyToBuy = 0
	s.QtyToSell = 0
	s.NumOrderToBuy = 0
	s.NumOrderToSell = 0
	s.IsReset = true
}

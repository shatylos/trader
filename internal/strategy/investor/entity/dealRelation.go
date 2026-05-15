package entity

import (
	"fmt"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/trading"
	"time"
)

type DealRelation struct {
	Deal            *Deal
	Orders          []*Order
	AverageBuyPrice float64
	PriceToBuy      float64
	PriceToSell     float64
	QtyToBuy        float64
	QtyToSell       float64
	NumOrderToBuy   int
	NumOrderToSell  int
	RevenueMainCur  float64
	RevenueTradeCur float64
	RealizedPNL     float64
	UnrealizedPNL   float64
	QtyInTrade      float64
}

func (d *DealRelation) GetTotalAmountBefore(mainCurrency, tradeCurrency string) (amount float64) {
	if d.Deal.ClosedTime.IsZero() {
		return
	}
	monthKey := fmt.Sprintf("%d-%d", d.Deal.ClosedTime.Year(), d.Deal.ClosedTime.Month())
	var takenOrderTime time.Time
	for _, order := range d.Orders {
		orderMonthKey := fmt.Sprintf("%d-%d", order.CreatedTime.Year(), order.CreatedTime.Month())
		if orderMonthKey != monthKey {
			continue
		}
		if takenOrderTime.IsZero() || takenOrderTime.After(order.CreatedTime) {
			mainCurr := trading.CurrencyAmountTotal(&order.WalletBefore, mainCurrency)
			tradeCurr := trading.CurrencyAmountTotal(&order.WalletBefore, tradeCurrency)
			amount = mainCurr + trading.TradeCurrencyToMain(tradeCurr, order.Price)
		}
	}
	return
}

func (d *DealRelation) CalcQtyInTrade() (qtyInTrade float64) {
	var boughtQty, soldQty float64

	for _, order := range d.Orders {
		if order.OrderStatus != structs.OrderStatuses.Filled {
			continue
		}

		if order.Side == structs.OrderSideBuy {
			boughtQty += order.Qty
		}
		if order.Side == structs.OrderSideSell {
			soldQty += order.Qty
		}
	}
	qtyInTrade = boughtQty - soldQty
	return
}

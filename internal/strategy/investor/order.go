package investor

import (
	"context"
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/tgNotifier"
	"github.com/shatylos/trader/tools/trading"
	"strconv"
	"time"
)

func (i *Investor) doBuy(ctx context.Context, timeFrameItem *_struct.TimeframeItem, state *entity.TimeframeState) (providerOrderId string, err error) {
	err = i.updateWalletInfo()
	if err != nil {
		err = apperrors.Wrap(err, "error update wallet info")
		return
	}

	walletBefore := *i.State.Wallet

	var qty, price float64
	qty, price = i.getQtyAndPriceToBuy(state)

	available := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.MainCurrency)
	if qty > math.Div(available, i.State.CurrentPrice) {
		err = apperrors.New("insufficient balance to buy %g%s. Available %g%s", qty, i.Config.TradeCurrency, available, i.Config.MainCurrency)
		return
	}

	if i.Config.Verbose {
		logger.Info(fmt.Sprintf("try to open limit order to buy %g%s. Price is %g", qty, i.Config.TradeCurrency, price))
	}

	providerOrderId, err = i.provider.OpenOrder(domainStructs.DomainOrderRequest{
		OrderId:     strconv.FormatInt(time.Now().UnixNano(), 10),
		Qty:         qty,
		Price:       price,
		ReduceOnly:  false,
		Side:        domainStructs.OrderSideBuy,
		Symbol:      i.Config.CoinPare,
		TimeInForce: "GTC",
		Type:        domainStructs.OrderTypes.Limit,
	})
	if err != nil {
		err = apperrors.Wrap(err, "error open buy order. qty: %g, price: %g, symbol: %s", qty, price, i.Config.CoinPare)
		return
	}

	msg := fmt.Sprintf("[%s] Placed order to buy %g %s by price %g for timeframe %s", i.Config.Id, qty, i.Config.TradeCurrency, price, timeFrameItem.Config.Resolution)
	logger.Info(msg)

	var domainOrder domainStructs.DomainOrder
	domainOrder, err = i.provider.GetOrder(providerOrderId)
	if err != nil {
		err = apperrors.Wrap(err, "error get order. ID: %s", providerOrderId)
		return
	}

	order := entity.Order{
		Timeframe:    timeFrameItem.Config.Resolution,
		DomainOrder:  domainOrder,
		WalletBefore: walletBefore,
		ConfigKey:    state.NumOrderToBuy - 1,
	}
	err = i.Storage.SaveOrder(ctx, &order)
	if err != nil {
		err = apperrors.Wrap(err, "error save order")
		return
	}

	state.ActiveOrder = &order
	state.LastBuyOrder = &order

	if order.OrderStatus == domainStructs.OrderStatuses.Filled {
		err = i.updateOrder(ctx, state, &order, timeFrameItem)
		if err != nil {
			err = apperrors.Wrap(err, "error update order")
			return
		}
	}

	return
}

func (i *Investor) doSell(ctx context.Context, timeFrameItem *_struct.TimeframeItem, state *entity.TimeframeState) (
	providerOrderId string, err error) {

	err = i.updateWalletInfo()
	if err != nil {
		err = apperrors.Wrap(err, "error update wallet info")
		return
	}

	walletBefore := *i.State.Wallet

	qty, price := i.getQtyAndPriceToSell(state)

	available := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.TradeCurrency)
	if qty > available {
		err = apperrors.New("insufficient balance to sell %g%s. Available %g%s", qty, i.Config.TradeCurrency, available, i.Config.TradeCurrency)
		return
	}

	if i.Config.Verbose {
		logger.Info(fmt.Sprintf("try to open limit order to sell %g%s. Price is %g",
			qty, i.Config.TradeCurrency, price))
	}

	providerOrderId, err = i.provider.OpenOrder(domainStructs.DomainOrderRequest{
		OrderId:     strconv.FormatInt(time.Now().UnixNano(), 10),
		Qty:         qty,
		Price:       price,
		ReduceOnly:  false,
		Side:        domainStructs.OrderSideSell,
		Symbol:      i.Config.CoinPare,
		TimeInForce: "GTC",
		Type:        domainStructs.OrderTypes.Limit,
	})
	if err != nil {
		err = apperrors.Wrap(err, "error open sell order. qty: %g, price: %g, symbol: %s", qty, price, i.Config.CoinPare)
		return
	}

	msg := fmt.Sprintf("[%s] Placed order to sell %g %s by price %g for timeframe %s", i.Config.Id, qty, i.Config.TradeCurrency, price, timeFrameItem.Config.Resolution)
	logger.Info(msg)

	var domainOrder domainStructs.DomainOrder
	domainOrder, err = i.provider.GetOrder(providerOrderId)
	if err != nil {
		err = apperrors.Wrap(err, "error get order. ID: %s", providerOrderId)
		return
	}

	order := entity.Order{
		Timeframe:    timeFrameItem.Config.Resolution,
		DomainOrder:  domainOrder,
		WalletBefore: walletBefore,
		ConfigKey:    state.NumOrderToSell - 1,
	}
	err = i.Storage.SaveOrder(ctx, &order)
	if err != nil {
		err = apperrors.Wrap(err, "error save order")
		return
	}

	state.ActiveOrder = &order
	state.LastSellOrder = &order

	if order.OrderStatus == domainStructs.OrderStatuses.Filled {
		err = i.updateOrder(ctx, state, &order, timeFrameItem)
		if err != nil {
			err = apperrors.Wrap(err, "error update order")
			return
		}
	}

	return
}

func (i *Investor) updateOrder(ctx context.Context, state *entity.TimeframeState, order *entity.Order, timeFrame _struct.Timeframe) (err error) {

	var updatedOrder domainStructs.DomainOrder
	updatedOrder, err = i.provider.GetOrder(order.OrderId)
	if err != nil {
		err = apperrors.Wrap(err, "error get order. ID: %s", order.OrderId)
		return
	}

	err = i.updateWalletInfo()
	if err != nil {
		err = apperrors.Wrap(err, "error update wallet info")
		return
	}

	if updatedOrder.OrderStatus == domainStructs.OrderStatuses.Filled {

		if updatedOrder.Side != domainStructs.OrderSideBuy && updatedOrder.Side != domainStructs.OrderSideSell {
			err = apperrors.New("unexpected order side %s", updatedOrder.Side)
			return
		}

		order.OrderStatus = updatedOrder.OrderStatus
		order.Price = updatedOrder.Price
		order.Qty = updatedOrder.Qty
		order.WalletAfter = *i.State.Wallet

		state.ApplyFilledOrder(order, i.Config.MinQty)
		state.ActiveOrder = nil

		err = i.Storage.SaveOrder(ctx, order)
		if err != nil {
			err = apperrors.Wrap(err, "error save order. DomainOrderID: %v, OrderID: %s", order.Id, order.OrderId)
			return
		}
		msg := fmt.Sprintf("[%s] Filled the %s order. Price: %.*f %s Qty: %.*f %s. For timeframe %s", i.Config.Id, order.Side, i.Config.PricePrecision, order.DomainOrder.Price, i.Config.MainCurrency, i.Config.QtyPrecision, order.DomainOrder.Qty, i.Config.TradeCurrency, timeFrame.Resolution())
		logger.Success(msg)
		if i.Config.TelegramNotifier {
			tgNotifier.Notify(msg)
		}
	}

	if (updatedOrder.OrderStatus == domainStructs.OrderStatuses.New ||
		updatedOrder.OrderStatus == domainStructs.OrderStatuses.Open) &&
		!i.State.Wallet.IsEqual(&order.WalletBefore) {

		order.WalletBefore = *i.State.Wallet
		err = i.Storage.SaveOrder(ctx, order)
		if err != nil {
			err = apperrors.Wrap(err, "error save order. DomainOrderID: %v, OrderID: %s", order.Id, order.OrderId)
			return
		}
	}

	return
}

func (i *Investor) doCancel(ctx context.Context, state *entity.TimeframeState, order *entity.Order) (err error) {
	err = i.provider.CancelOrder(order.OrderId, i.Config.CoinPare)
	if err != nil {
		err = apperrors.Wrap(err, "error cancel order. OrderID: %s", order.OrderId)
		return
	}

	order.OrderStatus = domainStructs.OrderStatuses.Canceled
	err = i.Storage.SaveOrder(ctx, order)
	if err != nil {
		err = apperrors.Wrap(err, "error save order. OrderID: %s", order.OrderId)
		return
	}

	state.ActiveOrder = nil

	err = i.updateWalletInfo()
	if err != nil {
		err = apperrors.Wrap(err, "error update wallet info")
		return
	}

	logger.Warning(fmt.Sprintf("cancelled %s order for timeframe %s", order.Side, order.Timeframe))

	return
}

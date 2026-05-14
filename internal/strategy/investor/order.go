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

func (i *Investor) doBuy(ctx context.Context, timeFrameItem *_struct.TimeframeItem, dealRelation *entity.DealRelation) (providerOrderId string, err error) {
	deal := dealRelation.Deal

	if deal.Id == nil {
		err = apperrors.New("deal struct must be saved before do buy")
		return
	}

	err = i.updateWalletInfo()
	if err != nil {
		err = apperrors.Wrap(err, "error update wallet info")
		return
	}

	var walletBefore, walletAfter domainStructs.DomainWallet
	walletBefore = *i.State.Wallet

	var qty, price float64
	qty, price = i.getQtyAndPriceToBuy(dealRelation)

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
		DealId:       *deal.Id,
		Timeframe:    timeFrameItem.Config.Resolution,
		DomainOrder:  domainOrder,
		WalletBefore: walletBefore,
		WalletAfter:  walletAfter,
		ConfigKey:    dealRelation.NumOrderToBuy - 1,
	}
	err = i.Storage.SaveOrder(ctx, &order)
	if err != nil {
		err = apperrors.Wrap(err, "error save order")
		return
	}

	if order.OrderStatus == domainStructs.OrderStatuses.Filled {
		err = i.updateOrder(ctx, dealRelation, &order, timeFrameItem)
		if err != nil {
			err = apperrors.Wrap(err, "error update order")
			return
		}
	}

	return
}

func (i *Investor) doSell(ctx context.Context, timeFrame _struct.Timeframe, dealRelation *entity.DealRelation) (
	providerOrderId string, err error) {
	deal := dealRelation.Deal

	if deal.Id == nil {
		err = apperrors.New("deal struct must be saved before do sell")
		return
	}

	err = i.updateWalletInfo()
	if err != nil {
		err = apperrors.Wrap(err, "error update wallet info")
		return
	}

	var walletBefore, walletAfter domainStructs.DomainWallet
	walletBefore = *i.State.Wallet

	qty, price := i.getQtyAndPriceToSell(dealRelation)

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

	msg := fmt.Sprintf("[%s] Placed order to sell %g %s by price %g for timeframe %s", i.Config.Id, qty, i.Config.TradeCurrency, price, timeFrame.Resolution())
	logger.Info(msg)

	var domainOrder domainStructs.DomainOrder
	domainOrder, err = i.provider.GetOrder(providerOrderId)
	if err != nil {
		err = apperrors.Wrap(err, "error get order. ID: %s", providerOrderId)
		return
	}

	order := entity.Order{
		DealId:       *deal.Id,
		Timeframe:    timeFrame.Resolution(),
		DomainOrder:  domainOrder,
		WalletBefore: walletBefore,
		WalletAfter:  walletAfter,
		ConfigKey:    dealRelation.NumOrderToSell - 1,
	}
	err = i.Storage.SaveOrder(ctx, &order)
	if err != nil {
		err = apperrors.Wrap(err, "error save order")
		return
	}

	if order.OrderStatus == domainStructs.OrderStatuses.Filled {
		err = i.updateOrder(ctx, dealRelation, &order, timeFrame)
		if err != nil {
			err = apperrors.Wrap(err, "error update order")
			return
		}
	}

	return
}

func (i *Investor) updateOrder(ctx context.Context, dealRelation *entity.DealRelation, order *entity.Order, timeFrame _struct.Timeframe) (err error) {

	deal := dealRelation.Deal
	err = i.updateWalletInfo()
	if err != nil {
		err = apperrors.Wrap(err, "error update wallet info")
		return
	}

	var updatedOrder domainStructs.DomainOrder
	updatedOrder, err = i.provider.GetOrder(order.OrderId)
	if err != nil {
		err = apperrors.Wrap(err, "error get order. ID: %s", order.OrderId)
		return
	}

	if updatedOrder.OrderStatus == domainStructs.OrderStatuses.Filled {

		err = i.updateWalletInfo()
		if err != nil {
			err = apperrors.Wrap(err, "error update wallet info")
			return
		}

		if updatedOrder.Side == domainStructs.OrderSideBuy {
			deal.Status = entity.DealStatusActive
		} else if updatedOrder.Side == domainStructs.OrderSideSell {
			if dealRelation.CalcQtyInTrade() <= i.Config.MinQty {
				deal.SetClose()
			}
		} else {
			err = apperrors.New("unexpected order side %s", updatedOrder.Side)
			return
		}

		order.OrderStatus = updatedOrder.OrderStatus
		order.Price = updatedOrder.Price
		order.Qty = updatedOrder.Qty
		order.WalletAfter = *i.State.Wallet
		err = i.Storage.SaveOrder(ctx, order)
		if err != nil {
			err = apperrors.Wrap(err, "error save order. DomainOrderID: %v, OrderID: %s", order.Id, order.OrderId)
			return
		}
		err = i.Storage.SaveDeal(ctx, deal)
		if err != nil {
			err = apperrors.Wrap(err, "error save deal. DealID: %v", deal.Id)
			return
		}
		msg := fmt.Sprintf("[%s] Filled the %s order. Qty: %g %s for timeframe %s", i.Config.Id, order.Side, order.DomainOrder.Qty, i.Config.TradeCurrency, timeFrame.Resolution())
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

func (i *Investor) doCancel(ctx context.Context, order *entity.Order) (err error) {
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

	err = i.updateWalletInfo()
	if err != nil {
		err = apperrors.Wrap(err, "error update wallet info")
		return
	}

	logger.Warning(fmt.Sprintf("cancelled %s order for timeframe %s", order.Side, order.Timeframe))

	return
}

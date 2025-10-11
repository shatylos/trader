package investor

import (
	"context"
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/tgNotifier"
	"strconv"
	"time"
)

func (i *Investor) doBuy(ctx context.Context, timeFrameItem *Timeframe, deal *entity.Deal) (providerOrderId string, err error) {

	if deal.Id == nil {
		msg := "deal struct must be saved before do buy"
		logger.Error(msg)
		err = tools.AppError{Message: msg}
		return
	}

	err = i.updateWalletInfo()
	if err != nil {
		return
	}

	var walletBefore, walletAfter domainStructs.DomainWallet
	walletBefore = *i.Wallet

	var qty, price float64
	qty, price, err = i.calculateQtyToBuy(timeFrameItem)
	if err != nil {
		return
	}

	if i.config.Verbose {
		logger.Info(fmt.Sprintf("Try to open limit order to buy %g%s. Price is %g", qty, i.config.TradeCurrency, price))
	}

	providerOrderId, err = i.provider.OpenOrder(domainStructs.DomainOrderRequest{
		OrderId:     strconv.FormatInt(time.Now().UnixNano(), 10),
		Qty:         qty,
		Price:       price,
		ReduceOnly:  false,
		Side:        domainStructs.OrderSideBuy,
		Symbol:      i.config.CoinPare,
		TimeInForce: "GTC",
		Type:        domainStructs.OrderTypes.Limit,
	})
	if err != nil {
		return
	}

	msg := fmt.Sprintf("[%s] Placed order to buy %g %s by price %g for timeframe %s", i.config.Id, qty, i.config.TradeCurrency, price, timeFrameItem.Config.Resolution)
	logger.Info(msg)

	var domainOrder domainStructs.DomainOrder
	domainOrder, err = i.provider.GetOrder(providerOrderId)
	if err != nil {
		return
	}

	order := entity.Order{
		DealId:       *deal.Id,
		Timeframe:    timeFrameItem.Config.Resolution,
		DomainOrder:  domainOrder,
		WalletBefore: walletBefore,
		WalletAfter:  walletAfter,
	}
	err = i.Storage.SaveOrder(ctx, &order)
	if err != nil {
		return
	}

	if order.OrderStatus == domainStructs.OrderStatuses.Filled {
		err = i.updateOrder(ctx, deal, &order, timeFrameItem)
		if err != nil {
			return
		}
	}

	return
}

func (i *Investor) calculateQtyToBuy(timeFrameItem *Timeframe) (qty float64, currentPrice float64, err error) {
	mainCurrencyAvailable := i.getMainCurrencyAvailable()
	if mainCurrencyAvailable == 0 {
		return
	}

	currentPrice = timeFrameItem.Candles[0].Close
	qtyPercent := timeFrameItem.Config.QtyPercent
	minQty := i.config.MinQty
	doIncreaseQtyToMinQty := i.config.DoIncreaseQtyToMinQty

	//qty = mainCurrencyAvailable / 100 * qtyPercent / currentPrice
	qty = math.Div(math.Mul(math.Div(mainCurrencyAvailable, 100), qtyPercent), currentPrice)

	if qty < minQty {
		if doIncreaseQtyToMinQty {
			qty = minQty
		} else {
			qty = 0
		}
	}
	qty = i.addCommission(qty, i.config.CommissionBuy)
	qty = math.RoundCell(qty, i.config.QtyPrecision)

	return
}

func (i *Investor) calculateQtyToSell(qty float64) (qtyResult float64) {
	//tradeAmountAvailable := currencyAmountAvailable(i.Wallet, i.config.TradeCurrency)
	//remainingBalance := tradeAmountAvailable - qty
	if qty > 0 {
		//minCoinReserve := qty / 100 * i.config.MinCoinReservePercent
		//minCoinReserve := math.Mul(math.Div(qty, 100), i.config.MinCoinReservePercent)
		// @TODO: Remove min_coin_reserve_percent and make the reduce always by config
		//if remainingBalance < minCoinReserve {
		qty = i.removeCommission(qty, i.config.CommissionBuy)
		//}
	}
	qtyResult = math.RoundFloor(qty, i.config.QtyPrecision)
	return
}

func (i *Investor) addCommission(qty, commission float64) (calculatedQty float64) {
	calculatedQty = math.Mul(math.Div(qty, 100), 100+commission)
	return
}
func (i *Investor) removeCommission(qty, commission float64) (calculatedQty float64) {
	calculatedQty = math.Mul(math.Div(qty, 100+commission), 100)
	return
}

func (i *Investor) doSell(ctx context.Context, timeFrameItem *Timeframe, deal *entity.Deal, qty, price float64) (
	providerOrderId string, err error) {

	if deal.Id == nil {
		msg := "deal struct must be saved before do sell"
		logger.Error(msg)
		err = tools.AppError{Message: msg}
		return
	}

	err = i.updateWalletInfo()
	if err != nil {
		return
	}

	var walletBefore, walletAfter domainStructs.DomainWallet
	walletBefore = *i.Wallet

	qty = i.calculateQtyToSell(qty)

	if i.config.Verbose {
		logger.Info(fmt.Sprintf("Try to open limit order to sell %g%s. Price is %g",
			qty, i.config.TradeCurrency, price))
	}

	providerOrderId, err = i.provider.OpenOrder(domainStructs.DomainOrderRequest{
		OrderId:     strconv.FormatInt(time.Now().UnixNano(), 10),
		Qty:         qty,
		Price:       price,
		ReduceOnly:  false,
		Side:        domainStructs.OrderSideSell,
		Symbol:      i.config.CoinPare,
		TimeInForce: "GTC",
		Type:        domainStructs.OrderTypes.Limit,
	})
	if err != nil {
		return
	}

	msg := fmt.Sprintf("[%s] Placed order to sell %g %s by price %g for timeframe %s", i.config.Id, qty, i.config.TradeCurrency, price, timeFrameItem.Config.Resolution)
	logger.Info(msg)

	var domainOrder domainStructs.DomainOrder
	domainOrder, err = i.provider.GetOrder(providerOrderId)
	if err != nil {
		return
	}

	order := entity.Order{
		DealId:       *deal.Id,
		Timeframe:    timeFrameItem.Config.Resolution,
		DomainOrder:  domainOrder,
		WalletBefore: walletBefore,
		WalletAfter:  walletAfter,
	}
	err = i.Storage.SaveOrder(ctx, &order)
	if err != nil {
		return
	}

	if order.OrderStatus == domainStructs.OrderStatuses.Filled {
		err = i.updateOrder(ctx, deal, &order, timeFrameItem)
		if err != nil {
			return
		}
	}

	return
}

func (i *Investor) updateOrder(ctx context.Context, deal *entity.Deal, order *entity.Order, timeFrameItem *Timeframe) (err error) {

	err = i.updateWalletInfo()
	if err != nil {
		return
	}

	var updatedOrder domainStructs.DomainOrder
	updatedOrder, err = i.provider.GetOrder(order.OrderId)
	if err != nil {
		return
	}

	if updatedOrder.OrderStatus == domainStructs.OrderStatuses.Filled {

		err = i.updateWalletInfo()
		if err != nil {
			return
		}

		if updatedOrder.Side == domainStructs.OrderSideBuy {
			deal.Status = entity.DealStatusActive
		} else if updatedOrder.Side == domainStructs.OrderSideSell {
			deal.SetClose()
		} else {
			msg := fmt.Sprintf("Unexpected order side %s", updatedOrder.Side)
			logger.Error(msg)
			err = tools.AppError{Message: msg}
			return
		}

		order.OrderStatus = updatedOrder.OrderStatus
		order.Price = updatedOrder.Price
		order.Qty = updatedOrder.Qty
		order.WalletAfter = *i.Wallet
		err = i.Storage.SaveOrder(ctx, order)
		if err != nil {
			return
		}
		err = i.Storage.SaveDeal(ctx, deal)
		if err != nil {
			return
		}
		msg := fmt.Sprintf("[%s] Filled the %s order. Qty: %g %s for timeframe %s", i.config.Id, order.Side, order.DomainOrder.Qty, i.config.TradeCurrency, timeFrameItem.Config.Resolution)
		logger.Success(msg)
		if i.config.TelegramNotifier {
			tgNotifier.Notify(msg)
		}
	}

	if (updatedOrder.OrderStatus == domainStructs.OrderStatuses.New ||
		updatedOrder.OrderStatus == domainStructs.OrderStatuses.Open) &&
		// @TODO: Implement isEqual function for wallet
		i.Wallet.UpdatedTime.After(order.WalletBefore.UpdatedTime) {

		order.WalletBefore = *i.Wallet
		err = i.Storage.SaveOrder(ctx, order)
		if err != nil {
			return
		}
	}

	return
}

func (i *Investor) doCancel(ctx context.Context, order *entity.Order) (err error) {
	err = i.provider.CancelOrder(order.OrderId, i.config.CoinPare)
	if err != nil {
		return
	}

	order.OrderStatus = domainStructs.OrderStatuses.Canceled
	err = i.Storage.SaveOrder(ctx, order)
	if err != nil {
		return
	}

	err = i.updateWalletInfo()
	if err != nil {
		return
	}

	logger.Warning(fmt.Sprintf("Cancelled %s order for timeframe %s", order.Side, order.Timeframe))

	return
}

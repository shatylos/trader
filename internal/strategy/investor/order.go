package investor

import (
	"context"
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/storage"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"strconv"
	"time"
)

func (i *Investor) doBuy(ctx context.Context, timeFrameItem *Timeframe, deal *storage.Deal) (order *storage.Order, providerOrderId string, err error) {

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

	var qty float64
	qty, err = i.calculateQtyToBuy(timeFrameItem)
	if err != nil {
		return
	}

	logger.Info(fmt.Sprintf("Try to buy %g%s. By market. Expected price is %g", qty, i.config.TradeCurrency, timeFrameItem.Candles[0].Close))

	providerOrderId, err = i.provider.OpenOrder(domainStructs.DomainOrderRequest{
		OrderId:     strconv.FormatInt(time.Now().UnixNano(), 10),
		Qty:         qty,
		ReduceOnly:  false,
		Side:        domainStructs.OrderSideBuy,
		Symbol:      i.config.CoinPare,
		TimeInForce: "GTC",
		Type:        domainStructs.OrderTypes.Market,
	})
	if err != nil {
		return
	}

	var domainOrder domainStructs.DomainOrder
	domainOrder, err = i.provider.GetOrder(providerOrderId)
	if err != nil {
		return
	}

	err = i.updateWalletInfo()
	if err != nil {
		return
	}
	walletAfter = *i.Wallet

	storageOrder := storage.Order{
		DealId:       *deal.Id,
		Timeframe:    timeFrameItem.Config.Resolution,
		DomainOrder:  domainOrder,
		WalletBefore: walletBefore,
		WalletAfter:  walletAfter,
	}
	err = i.Storage.SaveOrder(ctx, &storageOrder)
	if err != nil {
		return
	}
	order = &storageOrder

	deal.Status = storage.DealStatusActive
	err = i.Storage.SaveDeal(ctx, deal)
	if err != nil {
		return
	}

	return
}

func (i *Investor) calculateQtyToBuy(timeFrameItem *Timeframe) (qty float64, err error) {
	mainCurrencyAvailable := i.getMainCurrencyAvailable()
	if mainCurrencyAvailable == 0 {
		return
	}

	currentPrice := timeFrameItem.Candles[0].Close
	qtyPercent := timeFrameItem.Config.QtyPercent
	minQty := timeFrameItem.Config.MinQty
	doIncreaseQtyToMinQty := timeFrameItem.Config.DoIncreaseQtyToMinQty

	//qty = mainCurrencyAvailable / 100 * qtyPercent / currentPrice
	qty = math.Div(math.Mul(math.Div(mainCurrencyAvailable, 100), qtyPercent), currentPrice)
	qty = math.Round(qty, i.config.QtyPrecision)

	if qty < minQty {
		if doIncreaseQtyToMinQty {
			qty = minQty
		} else {
			qty = 0
		}
	}

	return
}

func (i *Investor) doSell(ctx context.Context, timeFrameItem *Timeframe, deal *storage.Deal, qty float64) (order *storage.Order, providerOrderId string, err error) {
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

	logger.Info(fmt.Sprintf("Try to sel %g%s. By market. Expected price is %g", qty, i.config.TradeCurrency, timeFrameItem.Candles[0].Close))

	qty = math.Round(qty, i.config.QtyPrecision)
	providerOrderId, err = i.provider.OpenOrder(domainStructs.DomainOrderRequest{
		OrderId:     strconv.FormatInt(time.Now().UnixNano(), 10),
		Qty:         qty,
		ReduceOnly:  false,
		Side:        domainStructs.OrderSideSell,
		Symbol:      i.config.CoinPare,
		TimeInForce: "GTC",
		Type:        domainStructs.OrderTypes.Market,
	})
	if err != nil {
		return
	}

	var domainOrder domainStructs.DomainOrder
	domainOrder, err = i.provider.GetOrder(providerOrderId)
	if err != nil {
		return
	}

	err = i.updateWalletInfo()
	if err != nil {
		return
	}
	walletAfter = *i.Wallet

	storageOrder := storage.Order{
		DealId:       *deal.Id,
		Timeframe:    timeFrameItem.Config.Resolution,
		DomainOrder:  domainOrder,
		WalletBefore: walletBefore,
		WalletAfter:  walletAfter,
	}
	err = i.Storage.SaveOrder(ctx, &storageOrder)
	if err != nil {
		return
	}
	order = &storageOrder

	deal.Status = storage.DealStatusClosed
	deal.ClosedTime = time.Now()
	err = i.Storage.SaveDeal(ctx, deal)
	if err != nil {
		return
	}

	return
}

package investor

import (
	"context"
	"fmt"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/tgNotifier"
)

func (i *Investor) handlePremium(ctx context.Context, state *entity.TimeframeState, timeFrameItem *_struct.TimeframeItem) (err error) {
	currentPrice := timeFrameItem.Candles[0].Close

	if state.PriceToSell == 0 {
		timeFrameItem.TradeStateMsg = "Handle premium. Will not sell"
		return
	}

	if state.ActiveOrder != nil {
		timeFrameItem.TradeStateMsg = "There is an active order. Wait for fill the order"
		return
	}

	if currentPrice < state.PriceToSell {
		timeFrameItem.TradeStateMsg = fmt.Sprintf("Handle premium. Expect higher price (%.2f) to sell", state.PriceToSell)
		return
	}

	if state.QtyInTrade == 0 {
		timeFrameItem.TradeStateMsg = "Handle premium. Qty in trade is 0. Will not sell"
		return
	}

	if state.QtyInTrade < i.Config.MinQty {
		err = i.resetTimeframeState(ctx, state)
		if err != nil {
			err = apperrors.Wrap(err, "error reset timeframe state")
			return
		}
		timeFrameItem.TradeStateMsg = "Handle premium. Not enough qty to sell. State was reset."
		return
	}

	var providerOrderId string
	providerOrderId, err = i.doSell(ctx, timeFrameItem, state)
	if err != nil {
		if providerOrderId != "" {
			i.Config.Enabled = false
			msg := fmt.Sprintf("[%s] Order for sell was created at the provider's side but error happened. Configuration disabled.", i.Config.Id)
			logger.Warning(msg)
			tgNotifier.Notify(msg)
		}
		err = apperrors.Wrap(err, "error do sell")
		return
	}

	return
}

func (i *Investor) handleDiscount(ctx context.Context, state *entity.TimeframeState, timeFrameItem *_struct.TimeframeItem) (err error) {

	if !timeFrameItem.Config.CanOpenNewOrder {
		timeFrameItem.TradeStateMsg = "Open new order disabled in config"
		return
	}

	if state.PriceToBuy == 0 {
		timeFrameItem.TradeStateMsg = "Handle discount. Will not buy more"
		return
	}

	currentPrice := timeFrameItem.Candles[0].Close
	if currentPrice > state.PriceToBuy {
		timeFrameItem.TradeStateMsg = fmt.Sprintf("Handle discount. Expect lower price (%.2f) to buy", state.PriceToBuy)
		return
	}

	if state.ActiveOrder != nil {
		timeFrameItem.TradeStateMsg = "There is an active order"
		return
	}

	var providerOrderId string
	providerOrderId, err = i.doBuy(ctx, timeFrameItem, state)
	if err != nil {
		if providerOrderId != "" {
			i.Config.Enabled = false
			msg := fmt.Sprintf("[%s] Order for buy was created at the provider's side but error happened. Configuration disabled.", i.Config.Id)
			logger.Warning(msg)
			tgNotifier.Notify(msg)
		}
		err = apperrors.Wrap(err, "error do buy")
		return
	}

	return
}

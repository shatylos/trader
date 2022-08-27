package bybit

import (
	"bitbucket.org/shatylos/trader/domain/constant"
	"bitbucket.org/shatylos/trader/domain/domains/bybit/request"
	"bitbucket.org/shatylos/trader/domain/structs"
)

type DomainBybitMargin struct {
	IsDemo bool
}

func (d *DomainBybitMargin) GetType() int64 {
	return constant.DomainTypeMargin
}

func (d *DomainBybitMargin) GetWallet() (*structs.DomainWallet, error) {
	walletBalances, er := request.GetWalletBalance(d.IsDemoMode())
	if er != nil {
		return nil, er
	}

	var availableCoins []structs.DomainWalletCoinItem
	var reservedCoins []structs.DomainWalletCoinItem

	for coinCode, walletBalance := range *walletBalances {
		if walletBalance.AvailableBalance == 0 && walletBalance.UsedMargin == 0 {
			continue
		}
		availableCoins = append(availableCoins, structs.DomainWalletCoinItem{
			Coin:   coinCode,
			Amount: walletBalance.AvailableBalance,
		})
		reservedCoins = append(reservedCoins, structs.DomainWalletCoinItem{
			Coin:   coinCode,
			Amount: walletBalance.UsedMargin,
		})
	}

	result := structs.DomainWallet{
		Available: availableCoins,
		Reserved:  reservedCoins,
	}

	return &result, nil
}

func (d *DomainBybitMargin) IsDemoMode() bool {
	return d.IsDemo
}

func (d *DomainBybitMargin) LoadCandleHistory(symbol string, resolution string, from int64, limit int64) ([]structs.DomainCandle, error) {

	candles, err := request.GetKlineList(symbol, resolution, from, limit, d.IsDemo)
	if err != nil {
		return nil, err
	}

	candlesResult := make([]structs.DomainCandle, len(candles))

	for i, candle := range candles {
		candlesResult[i] = structs.DomainCandle{
			Time:   int64(candle.OpenTime),
			High:   candle.High,
			Low:    candle.Low,
			Open:   candle.Open,
			Close:  candle.Close,
			Volume: candle.Volume,
		}
	}

	return candlesResult, nil
}

func (d *DomainBybitMargin) GetPositionList(coinPare string) ([]structs.DomainPosition, error) {

	positions, err := request.GetPositionList(coinPare, d.IsDemo)
	if err != nil {
		return nil, err
	}
	resultPositions := make([]structs.DomainPosition, len(positions))

	for i, position := range positions {
		resultPosition := structs.DomainPosition{
			Price:         position.EntryPrice,
			Leverage:      int64(position.Leverage),
			Symbol:        position.Symbol,
			Qty:           position.Size,
			RealizedPnl:   position.RealisedPnl,
			StopLoss:      position.StopLoss,
			TakeProfit:    position.TakeProfit,
			Side:          position.Side,
			UnrealizedPnl: position.UnrealisedPnl,
		}
		resultPositions[i] = resultPosition
	}

	return resultPositions, nil
}

func (d *DomainBybitMargin) OpenPosition(positionRequest structs.DomainPositionRequest) (string, error) {

	orderRequest := request.OrderRequest{
		CloseOnTrigger: false,
		OrderLinkId:    positionRequest.PositionId,
		OrderType:      positionRequest.Type,
		Price:          positionRequest.Price,
		Qty:            positionRequest.Qty,
		ReduceOnly:     false,
		Side:           positionRequest.Side,
		StopLoss:       positionRequest.StopLoss,
		Symbol:         positionRequest.Symbol,
		TakeProfit:     positionRequest.TakeProfit,
		TimeInForce:    "FillOrKill",
	}

	order, err := request.CreateOrder(orderRequest, d.IsDemo)
	if err != nil {
		return "", err
	}

	return order.OrderId, nil
}

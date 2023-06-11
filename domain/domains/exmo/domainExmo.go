package exmo

import (
	"bitbucket.org/shatylos/trader/domain/constant"
	"bitbucket.org/shatylos/trader/domain/domains/exmo/request"
	"bitbucket.org/shatylos/trader/domain/structs"
	tradingConstant "bitbucket.org/shatylos/trader/trading/constant"
)

type DomainExmo struct {
	isDemo bool
}

func (d *DomainExmo) GetType() int64 {
	return constant.DomainTypeSpot
}

func (d *DomainExmo) GetWallet() (*structs.DomainWallet, error) {
	userInfo, er := request.GetUserInfo()
	if er != nil {
		return nil, er
	}

	var availableCoins []structs.DomainWalletCoinItem
	for coinCode, coinAmount := range userInfo.Balances {
		availableCoins = append(availableCoins, structs.DomainWalletCoinItem{
			Coin:   coinCode,
			Amount: coinAmount,
		})
	}

	var reservedCoins []structs.DomainWalletCoinItem
	for coinCode, coinAmount := range userInfo.Reserved {
		reservedCoins = append(reservedCoins, structs.DomainWalletCoinItem{
			Coin:   coinCode,
			Amount: coinAmount,
		})
	}

	result := structs.DomainWallet{
		Available: availableCoins,
		Reserved:  reservedCoins,
	}

	return &result, nil
}

func (d *DomainExmo) IsDemoMode() bool {
	return false
}

func (d *DomainExmo) LoadCandleHistory(symbol string, resolution string, from int64, limit int64) ([]structs.DomainCandle, error) {
	to := from + (tradingConstant.ResolToSec[resolution] * limit)
	candles, err := request.LoadCandleHistory(symbol, resolution, from, to)
	if err != nil {
		return nil, err
	}

	candlesResult := make([]structs.DomainCandle, len(candles))

	for i, candle := range candles {
		candlesResult[i] = structs.DomainCandle{
			Time:   candle.T / 1000,
			High:   candle.H,
			Low:    candle.L,
			Open:   candle.O,
			Close:  candle.C,
			Volume: candle.V,
		}
	}

	return candlesResult, nil
}

func (d *DomainExmo) GetOpenOrderList(coinPare string) ([]structs.DomainOrder, error) {
	panic("not implemented")
}

func (d *DomainExmo) GetPositionList(coinPare string) ([]structs.DomainPosition, error) {
	panic("Not implemented")
}

func (d *DomainExmo) OpenPosition(positionRequest structs.DomainPositionRequest) (string, error) {
	panic("Not implemented")
}

func (d *DomainExmo) OpenOrder(orderRequest structs.DomainOrderRequest) (string, error) {
	panic("Not implemented")
}

func (d *DomainExmo) CancelOrder(orderId string) error {
	panic("Not implemented")
}

func (d *DomainExmo) GetHistoryOrders(limit int64) ([]structs.DomainOrder, error) {
	panic("Not implemented")
}

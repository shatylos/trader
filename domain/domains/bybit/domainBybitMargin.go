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

func (d *DomainBybitMargin) LoadCandleHistory(symbol string, resolution string, from int64, to int64) ([]structs.DomainCandle, error) {
	panic("Not implemented")
	// No need to map symbols or resolution. We use the same symbols like exmo

	//candles, err := request.LoadCandleHistory(symbol, resolution, from, to)
	//if err != nil {
	//	return nil, err
	//}
	//
	//candlesResult := make([]structs.DomainCandle, len(candles))
	//
	//for i, candle := range candles {
	//	candlesResult[i] = structs.DomainCandle{
	//		Time:  candle.T / 1000,
	//		High:  candle.H,
	//		Low:   candle.L,
	//		Open:  candle.O,
	//		Close: candle.C,
	//	}
	//}
	//
	//return candlesResult, nil
}

func (d *DomainBybitMargin) GetPositionList() ([]structs.DomainPosition, error) {
	panic("Not implemented")
}

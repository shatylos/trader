package bybit

import (
	"bitbucket.org/shatylos/trader/domain/constant"
	"bitbucket.org/shatylos/trader/domain/domains/bybit/request"
	"bitbucket.org/shatylos/trader/domain/structs"
	"bitbucket.org/shatylos/trader/utils"
	"fmt"
	"strconv"
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
		high, er := strconv.ParseFloat(candle.High, 64)
		if er != nil {
			return nil, utils.AppError{Message: fmt.Sprintf("Can not parse float value in candle.High. Source value is \"%s\"", candle.High)}
		}
		low, er := strconv.ParseFloat(candle.Low, 64)
		if er != nil {
			return nil, utils.AppError{Message: fmt.Sprintf("Can not parse float value in candle.Low. Source value is \"%s\"", candle.Low)}
		}
		open, er := strconv.ParseFloat(candle.Open, 64)
		if er != nil {
			return nil, utils.AppError{Message: fmt.Sprintf("Can not parse float value in candle.Open. Source value is \"%s\"", candle.Open)}
		}
		_close, er := strconv.ParseFloat(candle.Close, 64)
		if er != nil {
			return nil, utils.AppError{Message: fmt.Sprintf("Can not parse float value in candle.Close. Source value is \"%s\"", candle.Close)}
		}
		volume, er := strconv.ParseFloat(candle.Volume, 64)
		if er != nil {
			return nil, utils.AppError{Message: fmt.Sprintf("Can not parse float value in candle.Volume. Source value is \"%s\"", candle.Volume)}
		}

		candlesResult[i] = structs.DomainCandle{
			Time:   int64(candle.OpenTime),
			High:   high,
			Low:    low,
			Open:   open,
			Close:  _close,
			Volume: volume,
		}
	}

	return candlesResult, nil
}

func (d *DomainBybitMargin) GetPositionList() ([]structs.DomainPosition, error) {
	panic("Not implemented")
}

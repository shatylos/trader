package exmo

import (
	"bitbucket.org/shatylos/trader/domain/constant"
	"bitbucket.org/shatylos/trader/domain/domains/exmo/request"
	"bitbucket.org/shatylos/trader/domain/structs"
)

type DomainExmoMargin struct {
	isDemo bool
}

func (d *DomainExmoMargin) GetType() int64 {
	return constant.DomainTypeMargin
}

func (d *DomainExmoMargin) GetWallet() (*structs.DomainWallet, error) {
	marginWalletList, er := request.GetMarginWalletList()
	if er != nil {
		return nil, er
	}

	var availableCoins []structs.DomainWalletCoinItem
	var reservedCoins []structs.DomainWalletCoinItem

	for coinCode, coinValues := range marginWalletList.Wallets {
		availableCoins = append(availableCoins, structs.DomainWalletCoinItem{
			Coin:   coinCode,
			Amount: coinValues.Free,
		})
		reservedCoins = append(reservedCoins, structs.DomainWalletCoinItem{
			Coin:   coinCode,
			Amount: coinValues.Used,
		})
	}

	result := structs.DomainWallet{
		Available: availableCoins,
		Reserved:  reservedCoins,
	}

	return &result, nil
}

func (d *DomainExmoMargin) IsDemoMode() bool {
	return false
}

func (d *DomainExmoMargin) LoadCandleHistory(symbol string, resolution string, from int64, to int64) ([]structs.DomainCandle, error) {
	// No need to map symbols or resolution. We use the same symbols like exmo

	candles, err := request.LoadCandleHistory(symbol, resolution, from, to)
	if err != nil {
		return nil, err
	}

	candlesResult := make([]structs.DomainCandle, len(candles))

	for i, candle := range candles {
		candlesResult[i] = structs.DomainCandle{
			Time:  candle.T / 1000,
			High:  candle.H,
			Low:   candle.L,
			Open:  candle.O,
			Close: candle.C,
		}
	}

	return candlesResult, nil
}

func (d *DomainExmoMargin) GetPositionList() ([]structs.DomainPosition, error) {
	marginPositions, err := request.GetMarginPositionList()
	if err != nil {
		return nil, err
	}

	result := make([]structs.DomainPosition, len(marginPositions))

	for i, position := range marginPositions {
		result[i] = structs.DomainPosition{
			Margin:          position.Margin,
			Type:            position.Type,
			Quantity:        position.Quantity,
			Leverage:        position.Leverage,
			BasePrice:       position.BasePrice,
			Pair:            position.Pair,
			FundingQuantity: position.FundingQuantity,
			TakeProfit:      position.TakeProfit,
			StopLoss:        position.StopLoss,
			MarginCurrency:  position.MarginCurrency,
			Roe:             position.Roe,
			PositionId:      position.PositionId,
			UnrealizedPnl:   position.UnrealizedPnl,
			RealizedPnl:     position.RealizedPnl,
			PnlCurrency:     position.PnlCurrency,
			FundingCurrency: position.FundingCurrency,
		}
	}

	return result, nil
}

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

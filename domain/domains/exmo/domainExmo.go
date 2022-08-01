package exmo

import (
	"bitbucket.org/shatylos/trader/domain/constant"
	"bitbucket.org/shatylos/trader/domain/domains/exmo/request"
	"bitbucket.org/shatylos/trader/domain/structs"
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

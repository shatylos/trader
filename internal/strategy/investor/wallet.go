package investor

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
)

func (i *Investor) updateWalletInfo() (err error) {
	var wallet domainStructs.DomainWallet
	wallet, err = i.provider.GetWallet()
	if err != nil {
		return
	}
	i.State.Wallet = &wallet
	return
}

func (i *Investor) getMainCurrencyAvailable() (mainCurrencyAvailable float64) {
	for _, coin := range i.State.Wallet.Available {
		if coin.Coin == i.Config.MainCurrency {
			mainCurrencyAvailable = coin.Amount
		}
	}
	return
}

func (i *Investor) getTradeCurrencyAvailable() (tradeCurrencyAvailable float64) {
	for _, coin := range i.State.Wallet.Available {
		if coin.Coin == i.Config.TradeCurrency {
			tradeCurrencyAvailable = coin.Amount
		}
	}
	return
}

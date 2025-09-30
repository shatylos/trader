package investor

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/math"
)

func (i *Investor) updateWalletInfo() (err error) {
	var wallet domainStructs.DomainWallet
	wallet, err = i.provider.GetWallet()
	if err != nil {
		return
	}
	i.Wallet = &wallet
	return
}

func (i *Investor) getMainCurrencyAvailable() (mainCurrencyAvailable float64) {
	for _, coin := range i.Wallet.Available {
		if coin.Coin == i.config.MainCurrency {
			mainCurrencyAvailable = coin.Amount
		}
	}
	return
}

func (i *Investor) getTradeCurrencyAvailable() (tradeCurrencyAvailable float64) {
	for _, coin := range i.Wallet.Available {
		if coin.Coin == i.config.TradeCurrency {
			tradeCurrencyAvailable = coin.Amount
		}
	}
	return
}

func currencyAmountAvailable(wallet *domainStructs.DomainWallet, currency string) (amount float64) {
	for _, coin := range wallet.Available {
		if coin.Coin == currency {
			amount = coin.Amount
		}
	}
	return
}

func currencyAmountReserved(wallet *domainStructs.DomainWallet, currency string) (amount float64) {
	for _, coin := range wallet.Reserved {
		if coin.Coin == currency {
			amount = coin.Amount
		}
	}
	return
}

func currencyAmountTotal(wallet *domainStructs.DomainWallet, currency string) (amount float64) {
	amount = currencyAmountAvailable(wallet, currency) + currencyAmountReserved(wallet, currency)
	return
}

func tradeCurrencyToMain(tradeAmount, price float64) (mainAmount float64) {
	mainAmount = math.Mul(tradeAmount, price)
	return
}

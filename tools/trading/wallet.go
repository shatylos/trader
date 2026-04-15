package trading

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/math"
)

func CurrencyAmountAvailable(wallet *domainStructs.DomainWallet, currency string) (amount float64) {
	for _, coin := range wallet.Available {
		if coin.Coin == currency {
			amount = coin.Amount
		}
	}
	return
}

func CurrencyAmountReserved(wallet *domainStructs.DomainWallet, currency string) (amount float64) {
	for _, coin := range wallet.Reserved {
		if coin.Coin == currency {
			amount = coin.Amount
		}
	}
	return
}

func CurrencyAmountTotal(wallet *domainStructs.DomainWallet, currency string) (amount float64) {
	amount = CurrencyAmountAvailable(wallet, currency) + CurrencyAmountReserved(wallet, currency)
	return
}

func TradeCurrencyToMain(tradeAmount, price float64) (mainAmount float64) {
	mainAmount = math.Mul(tradeAmount, price)
	return
}

func MainCurrencyToTrade(mainAmount, price float64) (tradeAmount float64) {
	tradeAmount = math.Div(mainAmount, price)
	return
}

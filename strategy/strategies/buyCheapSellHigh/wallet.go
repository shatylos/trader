package buyCheapSellHigh

import "bitbucket.org/shatylos/trader/trading/services"

func (s *BuyCheapSellHigh) getBalances() (float64, float64, error) {
	wallet, err := services.LoadWalletInfo(s.DomainCode)
	if err != nil {
		return 0, 0, err
	}
	mainCurrencyBalance := float64(0)
	tradeCurrencyBalance := float64(0)

	for _, coin := range wallet.Available {
		if coin.Coin == s.MainCurrency {
			mainCurrencyBalance = coin.Amount
		}
		if coin.Coin == s.TradeCurrency {
			tradeCurrencyBalance = coin.Amount
		}
	}

	return mainCurrencyBalance, tradeCurrencyBalance, nil
}

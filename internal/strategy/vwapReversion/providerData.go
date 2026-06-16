package vwapReversion

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/apperrors"
)

func (v *VwapReversion) getAvailableBalance() (balance float64, err error) {
	var wallet domainStructs.DomainWallet
	wallet, err = v.provider.GetWallet()
	if err != nil {
		err = apperrors.Wrap(err, "error get wallet")
		return
	}
	for _, coin := range wallet.Available {
		if coin.Coin == v.config.MainCurrency {
			balance = coin.Amount
		}
	}
	return
}

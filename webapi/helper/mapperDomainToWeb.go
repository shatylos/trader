package helper

import (
	"bitbucket.org/shatylos/trader/domain/structs"
	gqlModel "bitbucket.org/shatylos/trader/webapi/graphql/model"
)

func MapWalletDomainToGraphql(domainWallet structs.DomainWalletApiResult) gqlModel.WalletBrokerWallet {

	availableItems := make([]*gqlModel.WalletBrokerWalletCoinItem, 0)
	for _, available := range domainWallet.DomainWallet.Available {
		if available.Amount > 0 {
			availableItems = append(availableItems, &gqlModel.WalletBrokerWalletCoinItem{
				Coin:   available.Coin,
				Amount: available.Amount,
			})
		}
	}

	reservedItems := make([]*gqlModel.WalletBrokerWalletCoinItem, 0)
	for _, reserved := range domainWallet.DomainWallet.Reserved {
		if reserved.Amount > 0 {
			reservedItems = append(reservedItems, &gqlModel.WalletBrokerWalletCoinItem{
				Coin:   reserved.Coin,
				Amount: reserved.Amount,
			})
		}
	}

	result := gqlModel.WalletBrokerWallet{
		DomainCode: domainWallet.DomainCode,
		Available:  availableItems,
		Reserved:   reservedItems,
	}

	return result
}

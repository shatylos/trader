package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"bitbucket.org/shatylos/trader/webapi/graphql/model"
)

// Wallet is the resolver for the wallet field.
func (r *queryResolver) Wallet(ctx context.Context) (*model.Wallet, error) {
	totalAm := 123.3
	bwAm1 := 500000.0
	brokerWallet1 := model.BrokerWallet{
		Amount: &bwAm1,
	}

	bwAm2 := 500000.0
	brokerWallet2 := model.BrokerWallet{
		Amount: &bwAm2,
	}

	brokerWallets := make([]*model.BrokerWallet, 0)
	brokerWallets = append(brokerWallets, &brokerWallet1)
	brokerWallets = append(brokerWallets, &brokerWallet2)

	wallet := model.Wallet{
		TotalAmount:   &totalAm,
		BrokerWallets: brokerWallets,
	}

	return &wallet, nil
}

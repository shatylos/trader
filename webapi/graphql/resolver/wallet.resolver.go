package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bitbucket.org/shatylos/trader/domain"
	"bitbucket.org/shatylos/trader/domain/structs"
	tradeServices "bitbucket.org/shatylos/trader/trading/services"
	"bitbucket.org/shatylos/trader/utils"
	gqlModel "bitbucket.org/shatylos/trader/webapi/graphql/model"
	"bitbucket.org/shatylos/trader/webapi/helper"
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
)

// Wallet is the resolver for the wallet field.
func (r *queryResolver) Wallet(ctx context.Context) (*gqlModel.Wallet, error) {
	domainCodes := domain.GetDomainList()

	walletResultSuccess = make([]gqlModel.WalletBrokerWallet, 0)
	walletResultFailed = make(map[string]error, 0)

	var wg sync.WaitGroup
	bufferRequest := make(chan bool, numberParallelRequests)
	domainWalletResultChan := make(chan structs.DomainWalletApiResult)

	go walletResponseHandle(&wg, domainWalletResultChan)
	for _, domainCode := range domainCodes {
		bufferRequest <- true
		wg.Add(1)
		go tradeServices.LoadWalletInfo(domainCode, domainWalletResultChan, bufferRequest)
	}
	wg.Wait()

	if len(walletResultFailed) > 0 {
		msg := "Service(s) unavailable:"
		for key, errorItem := range walletResultFailed {
			msg += fmt.Sprintf("\n%s: %s", key, errorItem)
		}
		return nil, utils.AppError{
			Message: msg,
		}
	} else {
		totalAvailablesList := make([]*gqlModel.WalletBrokerWalletCoinItem, 0)
		totalAvailablesMap := make(map[string]*gqlModel.WalletBrokerWalletCoinItem, 0)
		totalReservedsList := make([]*gqlModel.WalletBrokerWalletCoinItem, 0)
		totalReservedsMap := make(map[string]*gqlModel.WalletBrokerWalletCoinItem, 0)

		brokerWallets := make([]*gqlModel.WalletBrokerWallet, 0)

		for _, walletItem := range walletResultSuccess {
			for _, walletAvItem := range walletItem.Available {
				if walletAvItem.Amount > 0 {
					if _, ok := totalAvailablesMap[walletAvItem.Coin]; ok {
						totalAvailablesMap[walletAvItem.Coin].Amount += walletAvItem.Amount
					} else {
						totalAvailablesMap[walletAvItem.Coin] = walletAvItem
					}
				}
			}
			for _, walletResItem := range walletItem.Reserved {
				if walletResItem.Amount > 0 {
					if _, ok := totalReservedsMap[walletResItem.Coin]; ok {
						totalReservedsMap[walletResItem.Coin].Amount += walletResItem.Amount
					} else {
						totalReservedsMap[walletResItem.Coin] = walletResItem
					}
				}
			}
			brokerWallets = append(brokerWallets, &walletItem)
		}
		for _, coinItem := range totalAvailablesMap {
			totalAvailablesList = append(totalAvailablesList, coinItem)
		}
		for _, coinItem := range totalReservedsMap {
			totalReservedsList = append(totalReservedsList, coinItem)
		}

		wallet := gqlModel.Wallet{
			TotalAvailableAmount: totalAvailablesList,
			TotalReservedAmount:  totalReservedsList,
			BrokerWallets:        brokerWallets,
		}
		return &wallet, nil
	}
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
const defaultNumberParallelRequests = 5

var (
	numberParallelRequests int64
	walletResultSuccess    = make([]gqlModel.WalletBrokerWallet, 0)
	walletResultFailed     = make(map[string]error, 0)
)

func init() {
	envNumParRuq := os.Getenv("TRADER_NUMBER_PARALLEL_REQUESTS")
	numberParallelRequests = defaultNumberParallelRequests
	if envNumParRuq != "" {
		intNumParRuq, err := strconv.ParseInt(envNumParRuq, 0, 8)
		if err == nil {
			numberParallelRequests = intNumParRuq
		} else {
			println(fmt.Sprintf("Environment variable TRADER_NUMBER_PARALLEL_REQUESTS set not correctly. Use default value %d", defaultNumberParallelRequests))
		}
	}
}

func walletResponseHandle(wg *sync.WaitGroup, domainWalletResult chan structs.DomainWalletApiResult) {
	for domainWalletResultItem := range domainWalletResult {
		if domainWalletResultItem.Error == nil {
			walletResultSuccess = append(walletResultSuccess, helper.MapWalletDomainToGraphql(domainWalletResultItem))
		} else {
			walletResultFailed[domainWalletResultItem.DomainCode] = utils.AppError{
				Message: domainWalletResultItem.Error.Error(),
			}
		}
		wg.Done()
	}
}

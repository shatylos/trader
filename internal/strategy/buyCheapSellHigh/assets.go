package buyCheapSellHigh

import (
	"github.com/shatylos/trader/internal/strategy/buyCheapSellHigh/storage"
	"github.com/shatylos/trader/internal/strategy/buyCheapSellHigh/storage/mongo"
	"github.com/shatylos/trader/internal/strategy/buyCheapSellHigh/storage/structs"
	"time"
)

func (s *BuyCheapSellHigh) AddAssetTransaction(amount float64, dateTime time.Time, transactionType string) error {
	strategyStorage, err := storage.GetStorage(s.Id)
	if err != nil {
		return err
	}

	err = strategyStorage.AddAssetTransaction(structs.AssetTransaction{
		TransactionType: transactionType,
		Amount:          amount,
		CreatedTime:     dateTime.Unix(),
	})
	if err != nil {
		return err
	}

	modifyAmount := amount
	if transactionType == "withdraw" {
		modifyAmount *= -1
	}

	err = modifyOrders(strategyStorage, dateTime, modifyAmount)
	if err != nil {
		return err
	}

	return nil
}

func modifyOrders(strategyStorage *mongo.MongoStorage, from time.Time, amount float64) error {
	orders, err := strategyStorage.GetHistoryOrders(from, time.Now())
	if err != nil {
		return err
	}

	for _, order := range orders {
		order.MainCurrencyAmountBefore += amount
		order.AveragePrice = 0
		order.Revenue = 0
		err = strategyStorage.UpdateOrder(order)
		if err != nil {
			return err
		}
	}

	return nil
}

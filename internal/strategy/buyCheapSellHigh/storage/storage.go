package storage

import (
	"github.com/shatylos/trader/internal/strategy/buyCheapSellHigh/storage/mongo"
	"github.com/shatylos/trader/internal/strategy/buyCheapSellHigh/storage/structs"
	"time"
)

type Storage interface {
	AddDomainOrderOnce(order structs.HistoryOrder) (bool, error)
	GetCalculatedHistoryOrders(from time.Time, to time.Time) ([]structs.HistoryOrder, error)
	GetLastCalculatedOrder() (*structs.HistoryOrder, error)
	GetNotCalculatedHistoryOrders() ([]structs.HistoryOrder, error)
	GetNotFilledHistoryOrders() ([]structs.HistoryOrder, error)
	RemoveOrder(domainOrderId string) error
	ResetHistoryOrderData() error
	UpdateOrder(order structs.HistoryOrder) error
}

var storages = map[string]*mongo.MongoStorage{}

func GetStorage(setupId string) (*mongo.MongoStorage, error) {
	_, ok := storages[setupId]
	if ok {
		return storages[setupId], nil
	}

	var err error
	storages[setupId], err = mongo.New(setupId)
	if err != nil {
		return nil, err
	}

	return storages[setupId], nil
}

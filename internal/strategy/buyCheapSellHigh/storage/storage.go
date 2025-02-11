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

var storages = map[string]Storage{}

func GetStorage(setupId string) (*Storage, error) {
	err := error(nil)

	storage, ok := storages[setupId]
	if !ok {
		//storage, err = sqlite.New("var/buyCheapSellHigh.db", setupId)
		storage, err = mongo.New(setupId)
		if err != nil {
			return nil, err
		}
		storages[setupId] = storage
	}

	return &storage, nil
}

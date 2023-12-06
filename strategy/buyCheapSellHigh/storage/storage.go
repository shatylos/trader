package storage

import (
	"bitbucket.org/shatylos/trader/strategy/buyCheapSellHigh/storage/mongo"
	"bitbucket.org/shatylos/trader/strategy/buyCheapSellHigh/storage/structs"
)

type Storage interface {
	AddDomainOrderOnce(order structs.HistoryOrder) (bool, error)
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

package storage

import (
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/fibonacci/storage/mongo"
)

type AssetTransaction struct {
	Id              *string `bson:"_id,omitempty"`
	TransactionType string  `bson:"type"`
	Amount          float64 `bson:"amount"`
	CreatedTime     int64   `bson:"created_time"`
}

type Storage interface {
	InsertOrder(order structs.DomainOrder) (bool, error)
}

var storages = map[string]*mongo.MongoStorage{}

func GetStorage(setupId string) (storage mongo.MongoStorage, err error) {
	_, ok := storages[setupId]
	if ok {
		storage = *storages[setupId]
		return
	}

	storages[setupId], err = mongo.New(setupId)
	if err != nil {
		return
	}

	storage = *storages[setupId]
	return
}

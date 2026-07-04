package storage

import (
	"github.com/shatylos/trader/internal/strategy/scalper/storage/mongo"
	"github.com/shatylos/trader/tools/apperrors"
)

var storages = map[string]*mongo.MongoStorage{}

func GetStorage(setupId string) (storage mongo.MongoStorage, err error) {
	_, ok := storages[setupId]
	if ok {
		storage = *storages[setupId]
		return
	}

	storages[setupId], err = mongo.New(setupId)
	if err != nil {
		err = apperrors.Wrap(err, "error create new storage: %s", setupId)
		return
	}

	storage = *storages[setupId]
	return
}

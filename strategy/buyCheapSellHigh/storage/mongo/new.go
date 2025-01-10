package mongo

import (
	appStorage "github.com/shatylos/trader/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStorage struct {
	setupCode string
	db        *mongo.Database
}

func New(setupCode string) (*MongoStorage, error) {
	db, err := appStorage.GetDocumentDB()
	if err != nil {
		return nil, err
	}

	return &MongoStorage{
		setupCode: setupCode,
		db:        db,
	}, nil
}

func getOrderCollectionName(setupCode string) string {
	return `bcsh_orders_` + setupCode
}

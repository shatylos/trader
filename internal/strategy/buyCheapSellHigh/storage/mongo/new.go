package mongo

import (
	appStorage "github.com/shatylos/trader/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStorage struct {
	setupCode  string
	db         *mongo.Database
	collection *mongo.Collection
}

func New(setupCode string) (*MongoStorage, error) {
	db, err := appStorage.GetDocumentDB()
	if err != nil {
		return nil, err
	}

	collection := db.Collection(getOrderCollectionName(setupCode))

	return &MongoStorage{
		setupCode:  setupCode,
		db:         db,
		collection: collection,
	}, nil
}

func getOrderCollectionName(setupCode string) string {
	return `bcsh_orders_` + setupCode
}

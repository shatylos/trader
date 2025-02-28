package mongo

import (
	"fmt"
	appStorage "github.com/shatylos/trader/internal/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStorage struct {
	setupCode       string
	db              *mongo.Database
	orderCollection *mongo.Collection
	assetCollection *mongo.Collection
}

func New(setupCode string) (*MongoStorage, error) {
	db, err := appStorage.GetDocumentDB()
	if err != nil {
		return nil, err
	}

	orderCollection := db.Collection(getOrderCollectionName(setupCode))
	assetCollection := db.Collection(getAssetCollectionName(setupCode))

	return &MongoStorage{
		setupCode:       setupCode,
		db:              db,
		orderCollection: orderCollection,
		assetCollection: assetCollection,
	}, nil
}

func getOrderCollectionName(setupCode string) string {
	return fmt.Sprintf("bcsh_orders_%s", setupCode)
}

func getAssetCollectionName(setupCode string) string {
	return fmt.Sprintf("bcsh_assets_%s", setupCode)
}

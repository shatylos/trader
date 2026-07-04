package mongo

import (
	"fmt"
	appStorage "github.com/shatylos/trader/internal/storage"
	"github.com/shatylos/trader/tools/apperrors"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStorage struct {
	setupCode          string
	db                 *mongo.Database
	positionCollection *mongo.Collection
	assetCollection    *mongo.Collection
}

func New(setupCode string) (*MongoStorage, error) {
	db, err := appStorage.GetDocumentDB()
	if err != nil {
		return nil, apperrors.Wrap(err, "error get document DB: %s", setupCode)
	}

	positionCollection := db.Collection(getPositionCollectionName(setupCode))
	assetCollection := db.Collection(getAssetCollectionName(setupCode))

	return &MongoStorage{
		setupCode:          setupCode,
		db:                 db,
		positionCollection: positionCollection,
		assetCollection:    assetCollection,
	}, nil
}

func getPositionCollectionName(setupCode string) string {
	return fmt.Sprintf("scalper_positions_%s", setupCode)
}

func getAssetCollectionName(setupCode string) string {
	return fmt.Sprintf("scalper_assets_%s", setupCode)
}

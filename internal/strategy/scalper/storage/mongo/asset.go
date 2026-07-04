package mongo

import (
	"context"
	"fmt"
	"github.com/shatylos/trader/internal/strategy/scalper/structs"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (s *MongoStorage) AddAssetTransaction(transaction structs.AssetTransaction) error {

	ctx := context.Background()
	_, err := s.assetCollection.InsertOne(ctx, transaction)
	if err != nil {
		logger.Error(fmt.Sprintf("Error inserting asset transaction: %s", err.Error()))
		return apperrors.Wrap(err, "error inserting asset transaction")
	}

	return nil
}

func (s *MongoStorage) GetAssetTransactions(from time.Time, to time.Time) ([]structs.AssetTransaction, error) {
	ctx := context.Background()

	cursor, err := s.assetCollection.Find(ctx, bson.D{
		{"$and", bson.A{
			bson.D{{"CreatedTime", bson.D{{"$gt", from}}}},
			bson.D{{"CreatedTime", bson.D{{"$lt", to}}}},
		}},
	}, options.Find().SetSort(
		bson.D{{"CreatedTime", -1}},
	))
	if err != nil {
		return nil, apperrors.Wrap(err, "error finding assets")
	}
	defer cursor.Close(ctx)

	var assets []structs.AssetTransaction

	if err = cursor.All(ctx, &assets); err != nil {
		return nil, apperrors.Wrap(err, "error mapping assets")
	}

	return assets, nil
}

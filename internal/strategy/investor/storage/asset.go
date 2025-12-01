package storage

import (
	"context"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/apperrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (s *Storage) AddAssetTransaction(ctx context.Context, transaction structs.AssetTransaction) (err error) {
	_, err = s.assetCollection.InsertOne(ctx, transaction)
	if err != nil {
		err = apperrors.Wrap(err, "error inserting asset transaction")
		return
	}

	return
}

func (s *Storage) GetAssetTransactions(ctx context.Context, from time.Time, to time.Time) (assetsResult []*structs.AssetTransaction, err error) {
	var cursor *mongo.Cursor
	cursor, err = s.assetCollection.Find(ctx, bson.D{
		{"$and", bson.A{
			bson.D{{"CreatedTime", bson.D{{"$gt", from}}}},
			bson.D{{"CreatedTime", bson.D{{"$lt", to}}}},
		}},
	}, options.Find().SetSort(
		bson.D{{"CreatedTime", -1}},
	))
	if err != nil {
		err = apperrors.Wrap(err, "error finding assets")
		return
	}
	defer cursor.Close(ctx)

	var assets []structs.AssetTransaction
	if err = cursor.All(ctx, &assets); err != nil {
		err = apperrors.Wrap(err, "error mapping assets")
		return
	}

	for _, asset := range assets {
		assetsResult = append(assetsResult, &asset)
	}

	return
}

package storage

import (
	"context"
	"fmt"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (s *Storage) AddAssetTransaction(ctx context.Context, transaction structs.AssetTransaction) error {
	_, err := s.assetCollection.InsertOne(ctx, transaction)
	if err != nil {
		logger.Error(fmt.Sprintf("Error inserting asset transaction: %s", err.Error()))
		return err
	}

	return nil
}

func (s *Storage) GetAssetTransactions(ctx context.Context, from time.Time, to time.Time) (assetsResult []*structs.AssetTransaction, err error) {
	cursor, err := s.assetCollection.Find(ctx, bson.D{
		{"$and", bson.A{
			bson.D{{"CreatedTime", bson.D{{"$gt", from}}}},
			bson.D{{"CreatedTime", bson.D{{"$lt", to}}}},
		}},
	}, options.Find().SetSort(
		bson.D{{"CreatedTime", -1}},
	))
	if err != nil {
		msg := fmt.Sprintf("Error finding assets. Error: %s", err.Error())
		logger.Error(msg)
		err = tools.AppError{
			Message:     msg,
			ParentError: err,
		}
		return
	}
	defer cursor.Close(ctx)

	var assets []structs.AssetTransaction
	if err = cursor.All(ctx, &assets); err != nil {
		msg := fmt.Sprintf("Error mapping assets. Error: %s", err.Error())
		logger.Error(msg)
		err = tools.AppError{
			Message:     msg,
			ParentError: err,
		}
		return
	}

	for _, asset := range assets {
		assetsResult = append(assetsResult, &asset)
	}

	return
}

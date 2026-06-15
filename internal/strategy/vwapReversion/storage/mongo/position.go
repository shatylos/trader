package mongo

import (
	"context"
	"fmt"
	"github.com/shatylos/trader/internal/strategy/vwapReversion/structs"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (s *MongoStorage) SaveInternalPosition(position structs.Position) (savedPosition structs.Position, err error) {

	ctx := context.Background()
	var primObjectID primitive.ObjectID
	var savedId string
	var ok bool

	if position.Id == nil {
		position.CreatedTime = time.Now()
		position.UpdatedTime = time.Now()
		var inserted *mongo.InsertOneResult
		inserted, err = s.positionCollection.InsertOne(ctx, position)
		if err != nil {
			err = apperrors.Wrap(err, "error inserting new position")
			return
		}
		primObjectID, ok = inserted.InsertedID.(primitive.ObjectID)
		if !ok {
			err = apperrors.New("can not convert InsertedID to ObjectID")
			return
		}
		savedId = primObjectID.Hex()
	} else {
		primObjectID, err = primitive.ObjectIDFromHex(*position.Id)
		if err != nil {
			err = apperrors.Wrap(err, "can not convert Object ID from Hex updating position")
			return
		}
		filter := bson.D{{"_id", primObjectID}}

		update := bson.D{{"$set", bson.D{
			{"Chart", position.Chart},
			{"Trend", position.Trend},
			{"LtTrend", position.LtTrend},
			{"StTrend", position.StTrend},
			{"Order", position.Order},
			{"UpdatedTime", time.Now().Unix()},
			{"ClosedTime", position.ClosedTime},
			{"Status", position.Status},
			{"ProviderPosition", position.ProviderPosition},
			{"BalanceBefore", position.BalanceBefore},
			{"BalanceAfter", position.BalanceAfter},
			{"TotalClosePnl", position.TotalClosePnl},
		}}}

		_, err = s.positionCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			err = apperrors.Wrap(err, "error updating position")
			return
		}
		savedId = *position.Id
	}
	savedPosition, err = s.GetPositionById(savedId)
	if err != nil {
		err = apperrors.Wrap(err, "error get position by id")
		return
	}

	return
}

func (s *MongoStorage) GetLastInternalPosition(skip int64) (position structs.Position, err error) {
	ctx := context.Background()
	err = s.positionCollection.FindOne(ctx,
		bson.D{{}},
		options.FindOne().
			SetSort(bson.D{{"CreatedTime", -1}}).
			SetSkip(skip),
	).Decode(&position)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
		} else {
			err = apperrors.Wrap(err, "error get last internal position")
		}
	}
	return
}

func (s *MongoStorage) GetInternalPositionCreatedBeforeTime(dateTime time.Time) (position structs.Position, err error) {
	ctx := context.Background()
	err = s.positionCollection.FindOne(ctx,
		bson.D{{"CreatedTime", bson.D{{"$lt", dateTime.Unix()}}}},
		options.FindOne().SetSort(bson.D{{"CreatedTime", -1}}),
	).Decode(&position)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Warning(fmt.Sprintf("Not found position created before time %s", dateTime.String()))
			err = nil
		} else {
			err = apperrors.Wrap(err, "error get internal position created before time")
		}
	}
	return
}

func (s *MongoStorage) GetPositionById(idString string) (position structs.Position, err error) {
	var id primitive.ObjectID
	id, err = primitive.ObjectIDFromHex(idString)
	if err != nil {
		err = apperrors.Wrap(err, "can not convert Object ID from Hex")
		return
	}
	ctx := context.Background()
	err = s.positionCollection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&position)
	if err != nil {
		err = apperrors.Wrap(err, "error get position by id")
		return
	}
	return
}

func (s *MongoStorage) GetPositions(from time.Time, to time.Time) (positions []structs.Position, err error) {
	ctx := context.Background()
	var cursor *mongo.Cursor

	cursor, err = s.positionCollection.Find(ctx, bson.D{{
		"$and", bson.A{
			bson.D{{"CreatedTime", bson.D{{"$gt", from.Unix()}}}},
			bson.D{{"CreatedTime", bson.D{{"$lt", to.Unix()}}}},
		},
	}}, options.Find().SetSort(
		bson.D{{"CreatedTime", -1}},
	))
	if err != nil {
		return nil, apperrors.Wrap(err, "error find positions")
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &positions)
	if err != nil {
		err = apperrors.Wrap(err, "error map found positions")
		return
	}

	return
}

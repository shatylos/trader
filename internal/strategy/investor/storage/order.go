package storage

import (
	"context"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	"github.com/shatylos/trader/tools/apperrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (s *Storage) SaveOrder(ctx context.Context, order *entity.Order) (err error) {
	var primObjectID primitive.ObjectID
	var ok bool

	if order.Id == nil {
		order.CreatedTime = time.Now()
		order.UpdatedTime = time.Now()
		var inserted *mongo.InsertOneResult
		inserted, err = s.orderCollection.InsertOne(ctx, order)
		if err != nil {
			err = apperrors.Wrap(err, "error insert one")
		}
		primObjectID, ok = inserted.InsertedID.(primitive.ObjectID)
		if !ok {
			err = apperrors.New("can not convert InsertedID to ObjectID inserting order")
			return
		}
		orderId := primObjectID.Hex()
		order.Id = &orderId
	} else {
		primObjectID, err = primitive.ObjectIDFromHex(*order.Id)
		if err != nil {
			err = apperrors.Wrap(err, "can not convert Object ID from Hex updating order")
			return
		}
		filter := bson.D{{"_id", primObjectID}}
		order.UpdatedTime = time.Now()

		var updateDoc []byte
		updateDoc, err = bson.Marshal(order)
		if err != nil {
			err = apperrors.Wrap(err, "can not marshal order document to update")
			return
		}
		var m bson.M
		err = bson.Unmarshal(updateDoc, &m)
		if err != nil {
			err = apperrors.Wrap(err, "can not unmarshal order update document")
			return
		}
		delete(m, "_id")
		update := bson.D{{"$set", m}}

		_, err = s.orderCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			err = apperrors.Wrap(err, "error updating order")
		}
	}

	return
}

// @TODO: investigate how to add limit
func (s *Storage) GetOrdersByTimeframe(ctx context.Context, timeframe string) (ordersResult []*entity.Order, err error) {
	var cursor *mongo.Cursor
	cursor, err = s.orderCollection.Find(ctx,
		bson.D{{"Timeframe", timeframe}},
		options.Find().SetSort(bson.D{{"CreatedTime", 1}}),
	)
	if err != nil {
		err = apperrors.Wrap(err, "error find by timeframe: %s", timeframe)
		return
	}
	defer cursor.Close(ctx)

	orders := make([]entity.Order, 0)
	err = cursor.All(ctx, &orders)
	if err != nil {
		err = apperrors.Wrap(err, "error map found orders")
		return
	}

	for _, order := range orders {
		ordersResult = append(ordersResult, &order)
	}

	return
}

func (s *Storage) GetOrdersByPeriod(ctx context.Context, from time.Time, to time.Time) (ordersResult []*entity.Order, err error) {
	var cursor *mongo.Cursor
	cursor, err = s.orderCollection.Find(ctx,
		bson.D{{
			"$and", bson.A{
				bson.D{{"CreatedTime", bson.D{{"$gt", from}}}},
				bson.D{{"CreatedTime", bson.D{{"$lt", to}}}},
			},
		}},
		options.Find().SetSort(bson.D{{"CreatedTime", 1}}),
	)
	if err != nil {
		err = apperrors.Wrap(err, "error find orders by period from %s to %s", from, to)
		return
	}
	defer cursor.Close(ctx)

	orders := make([]entity.Order, 0)
	err = cursor.All(ctx, &orders)
	if err != nil {
		err = apperrors.Wrap(err, "error map found orders")
		return
	}

	for _, order := range orders {
		ordersResult = append(ordersResult, &order)
	}

	return
}

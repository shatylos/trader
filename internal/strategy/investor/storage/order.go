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

func (s *Storage) GetOrdersByDealId(ctx context.Context, dealId string) (ordersResult []*entity.Order, err error) {
	var cursor *mongo.Cursor
	cursor, err = s.orderCollection.Find(ctx,
		bson.D{{"DealId", dealId}},
		options.Find().SetSort(bson.D{{"CreatedTime", 1}}),
	)
	if err != nil {
		err = apperrors.Wrap(err, "error find by dealId: %s", dealId)
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

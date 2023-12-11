package mongo

import (
	"bitbucket.org/shatylos/trader/strategy/buyCheapSellHigh/storage/structs"
	"bitbucket.org/shatylos/trader/utils"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddDomainOrderOnce add order to storage.
func (s *MongoStorage) AddDomainOrderOnce(order structs.HistoryOrder) (bool, error) {

	collectionName := getOrderCollectionName(s.setupCode)
	ctx := context.TODO()

	var res structs.HistoryOrder
	err := s.db.Collection(collectionName).FindOne(ctx, bson.D{{"domain_order_id", order.DomainOrderId}}).Decode(&res)

	if errors.Is(err, mongo.ErrNoDocuments) {
		_, err := s.db.Collection(collectionName).InsertOne(ctx, order)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return false, utils.AppError{
		Message: fmt.Sprintf("Order with id %s already exists", order.DomainOrderId),
	}
}

func (s *MongoStorage) GetNotFilledHistoryOrders() ([]structs.HistoryOrder, error) {
	collectionName := getOrderCollectionName(s.setupCode)
	ctx := context.TODO()

	cursor, err := s.db.Collection(collectionName).Find(ctx, bson.D{
		{"$or", bson.A{
			bson.D{{"filled_price", 0}},
			bson.D{{"filled_qty", 0}},
			bson.D{{"side", ""}},
			bson.D{{"updated_time", 0}},
		}},
	})

	if err != nil {
		return nil, err
	}

	var orders []structs.HistoryOrder

	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *MongoStorage) GetNotCalculatedHistoryOrders() ([]structs.HistoryOrder, error) {
	collectionName := getOrderCollectionName(s.setupCode)
	ctx := context.TODO()

	cursor, err := s.db.Collection(collectionName).Find(ctx, bson.D{
		{"$and", bson.A{
			bson.D{{"filled_price", bson.D{{"$gt", 0}}}},
			bson.D{{"filled_qty", bson.D{{"$gt", 0}}}},
			bson.D{{"side", bson.D{{"$ne", ""}}}},
			bson.D{{"updated_time", bson.D{{"$gt", 0}}}},
			bson.D{{"average_price", bson.D{{"$in", bson.A{0, nil}}}}},
		}},
	}, options.Find().SetSort(
		bson.D{{"updated_time", 1}},
	))

	if err != nil {
		return nil, err
	}

	var orders []structs.HistoryOrder

	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *MongoStorage) GetLastCalculatedOrder() (*structs.HistoryOrder, error) {
	collectionName := getOrderCollectionName(s.setupCode)
	ctx := context.TODO()

	var lastCalculatedOrder structs.HistoryOrder
	err := s.db.Collection(collectionName).FindOne(
		ctx,
		bson.D{{
			"average_price", bson.D{{"$gt", 0}},
		}},
		options.FindOne().SetSort(
			bson.D{{"updated_time", -1}},
		),
	).Decode(&lastCalculatedOrder)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &lastCalculatedOrder, nil
}

func (s *MongoStorage) RemoveOrder(domainOrderId string) error {
	collectionName := getOrderCollectionName(s.setupCode)
	ctx := context.TODO()

	_, err := s.db.Collection(collectionName).DeleteOne(ctx, bson.D{{"domain_order_id", domainOrderId}})

	if err != nil {
		return err
	}

	return nil
}

func (s *MongoStorage) UpdateOrder(order structs.HistoryOrder) error {
	collectionName := getOrderCollectionName(s.setupCode)
	ctx := context.TODO()

	if order.Id == nil {
		return utils.AppError{
			Message: "order id is empty",
		}
	}
	id, err := primitive.ObjectIDFromHex(*order.Id)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id", id}}

	update := bson.D{{"$set", order}}

	_, err = s.db.Collection(collectionName).UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	return nil
}

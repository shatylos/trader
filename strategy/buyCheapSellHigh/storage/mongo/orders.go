package mongo

import (
	"bitbucket.org/shatylos/trader/strategy/buyCheapSellHigh/storage/structs"
	"bitbucket.org/shatylos/trader/utils"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

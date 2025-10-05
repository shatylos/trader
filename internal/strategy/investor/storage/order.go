package storage

import (
	"context"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Order struct {
	domainStructs.DomainOrder `bson:",inline"`
	DealId                    string                     `bson:"DealId"`
	Timeframe                 string                     `bson:"Timeframe"`
	CreatedTime               time.Time                  `bson:"CreatedTime"`
	UpdatedTime               time.Time                  `bson:"UpdatedTime"`
	WalletBefore              domainStructs.DomainWallet `bson:"WalletBefore"`
	WalletAfter               domainStructs.DomainWallet `bson:"WalletAfter"`
}

func (s *Storage) SaveOrder(ctx context.Context, order *Order) (err error) {
	var primObjectID primitive.ObjectID
	var ok bool

	if order.Id == nil {
		order.CreatedTime = time.Now()
		order.UpdatedTime = time.Now()
		var inserted *mongo.InsertOneResult
		inserted, err = s.orderCollection.InsertOne(ctx, order)
		if err != nil {
			err = tools.AppError{
				Message:     "Error inserting new order",
				ParentError: err,
			}
		}
		primObjectID, ok = inserted.InsertedID.(primitive.ObjectID)
		if !ok {
			err = tools.AppError{
				Message: "Can not convert InsertedID to ObjectID inserting order",
			}
			return
		}
		orderId := primObjectID.Hex()
		order.Id = &orderId
	} else {
		primObjectID, err = primitive.ObjectIDFromHex(*order.Id)
		if err != nil {
			err = tools.AppError{
				Message:     "Can not convert Object ID from Hex updating order",
				ParentError: err,
			}
			return
		}
		filter := bson.D{{"_id", primObjectID}}
		order.UpdatedTime = time.Now()

		var updateDoc []byte
		updateDoc, err = bson.Marshal(order)
		if err != nil {
			msg := "Can not marshal order document to update"
			logger.Error(msg)
			err = tools.AppError{
				Message:     msg,
				ParentError: err,
			}
			return
		}
		var m bson.M
		err = bson.Unmarshal(updateDoc, &m)
		if err != nil {
			msg := "Can not unmarshal order update document"
			logger.Error(msg)
			err = tools.AppError{
				Message:     msg,
				ParentError: err,
			}
			return
		}
		delete(m, "_id")
		update := bson.D{{"$set", m}}

		_, err = s.orderCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			err = tools.AppError{
				Message:     "Error updating order",
				ParentError: err,
			}
		}
	}

	return
}

func (s *Storage) GetOrdersByDealId(ctx context.Context, dealId string) (ordersResult []*Order, err error) {
	var cursor *mongo.Cursor
	cursor, err = s.orderCollection.Find(ctx,
		bson.D{{"DealId", dealId}},
		options.Find().SetSort(bson.D{{"CreatedTime", 1}}),
	)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	orders := make([]Order, 0)
	err = cursor.All(ctx, &orders)
	if err != nil {
		return
	}

	for _, order := range orders {
		ordersResult = append(ordersResult, &order)
	}

	return
}

func (o *Order) Amount() float64 {
	return math.Mul(o.Qty, o.Price)
}

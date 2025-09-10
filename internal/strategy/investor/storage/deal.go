package storage

import (
	"context"
	"errors"
	"github.com/shatylos/trader/tools"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Deal struct {
	Id          string    `bson:"_id,omitempty"`
	Timeframe   string    `bson:"Timeframe"`
	Status      string    `bson:"Status"`
	CreatedTime time.Time `bson:"CreatedTime"`
	UpdatedTime time.Time `bson:"UpdatedTime"`
	//IsHeap    bool
	//Orders ???
}

const DealStatusNew = "NEW"

func (s *Storage) SaveDeal(deal *Deal) (err error) {

	ctx := context.Background()
	var primObjectID primitive.ObjectID
	var ok bool

	if deal.Id == "" {
		deal.CreatedTime = time.Now()
		deal.UpdatedTime = time.Now()
		var inserted *mongo.InsertOneResult
		inserted, err = s.dealCollection.InsertOne(ctx, deal)
		if err != nil {
			err = tools.AppError{
				Message:     "Error inserting new deal",
				ParentError: err,
			}
		}
		primObjectID, ok = inserted.InsertedID.(primitive.ObjectID)
		if !ok {
			err = tools.AppError{
				Message: "Can not convert InsertedID to ObjectID",
			}
			return
		}
		//s.Id = primObjectID.Hex()
		deal.Id = primObjectID.Hex()
	} else {
		primObjectID, err = primitive.ObjectIDFromHex(deal.Id)
		if err != nil {
			err = tools.AppError{
				Message:     "Can not convert Object ID from Hex",
				ParentError: err,
			}
			return
		}
		filter := bson.D{{"_id", primObjectID}}
		deal.UpdatedTime = time.Now()
		update := bson.D{{"$set", deal}}

		_, err = s.dealCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			err = tools.AppError{
				Message:     "Error updating deal",
				ParentError: err,
			}
		}
	}

	return
}

func (s *Storage) GetLastDealByTimeframe(timeframe string, deal *Deal) (err error) {
	ctx := context.Background()

	err = s.dealCollection.FindOne(ctx,
		bson.D{{"Timeframe", timeframe}},
		options.FindOne().
			SetSort(bson.D{{"CreatedTime", -1}}),
	).Decode(deal)
	if errors.Is(err, mongo.ErrNoDocuments) {
		err = nil
	}
	if err != nil {
		return
	}

	return
}

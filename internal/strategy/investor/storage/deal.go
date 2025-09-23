package storage

import (
	"context"
	"errors"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Deal struct {
	Id          *string   `bson:"_id,omitempty"`
	Timeframe   string    `bson:"Timeframe"`
	Status      string    `bson:"Status"`
	CreatedTime time.Time `bson:"CreatedTime"`
	UpdatedTime time.Time `bson:"UpdatedTime"`
	//IsHeap    bool
	//Orders ???
}

type DealRelation struct {
	Deal   *Deal
	Orders []*Order
}

const DealStatusNew = "NEW"
const DealStatusActive = "ACTIVE"
const DealStatusClosed = "CLOSED"

func (s *Storage) SaveDeal(ctx context.Context, deal *Deal) (err error) {

	var primObjectID primitive.ObjectID
	var ok bool

	if deal.Id == nil {
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
				Message: "Can not convert InsertedID to ObjectID inserting deal",
			}
			return
		}
		dealId := primObjectID.Hex()
		deal.Id = &dealId
	} else {
		primObjectID, err = primitive.ObjectIDFromHex(*deal.Id)
		if err != nil {
			err = tools.AppError{
				Message:     "Can not convert Object ID from Hex updating deal",
				ParentError: err,
			}
			return
		}
		filter := bson.D{{"_id", primObjectID}}
		deal.UpdatedTime = time.Now()

		var updateDoc []byte
		updateDoc, err = bson.Marshal(deal)
		if err != nil {
			msg := "Can not marshal deal document to update"
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
			msg := "Can not unmarshal deal update document"
			logger.Error(msg)
			err = tools.AppError{
				Message:     msg,
				ParentError: err,
			}
			return
		}
		delete(m, "_id")
		update := bson.D{{"$set", m}}

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

func (s *Storage) GetLastDealByTimeframe(ctx context.Context, timeframe string, deal *Deal) (err error) {
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

func (s *Storage) GetDealsByPeriod(ctx context.Context, from time.Time, to time.Time) (dealPointers []*Deal, err error) {
	var cursor *mongo.Cursor

	cursor, err = s.dealCollection.Find(ctx, bson.D{{
		"$and", bson.A{
			bson.D{{"CreatedTime", bson.D{{"$gt", from}}}},
			bson.D{{"CreatedTime", bson.D{{"$lt", to}}}},
		},
	}}, options.Find().SetSort(
		bson.D{{"CreatedTime", -1}},
	))
	if err != nil {
		msg := "Error getting cursor deals by period"
		logger.Error(msg)
		err = tools.AppError{
			Message:     msg,
			ParentError: err,
		}
		return
	}
	defer cursor.Close(ctx)

	var deals []Deal
	err = cursor.All(ctx, &deals)
	if err != nil {
		msg := "Error getting all deals by period"
		logger.Error(msg)
		err = tools.AppError{
			Message:     msg,
			ParentError: err,
		}
		return
	}
	for _, deal := range deals {
		dealPointers = append(dealPointers, &deal)
	}

	return
}

func (s *Storage) GetDealRelation(ctx context.Context, deal *Deal) (dealRelation *DealRelation, err error) {
	if deal.Id == nil {
		msg := "Deal ID is empty. Can not get deal relations"
		logger.Error(msg)
		err = tools.AppError{Message: msg}
		return
	}

	var dealOrders []*Order
	dealOrders, err = s.GetOrdersByDealId(ctx, *deal.Id)
	if err != nil {
		return
	}

	dealRelation = &DealRelation{
		Deal:   deal,
		Orders: dealOrders,
	}
	return
}

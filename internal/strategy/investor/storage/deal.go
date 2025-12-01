package storage

import (
	"context"
	"errors"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/apperrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (s *Storage) SaveDeal(ctx context.Context, deal *entity.Deal) (err error) {

	var primObjectID primitive.ObjectID
	var ok bool

	if deal.Id == nil {
		deal.CreatedTime = time.Now()
		deal.UpdatedTime = time.Now()
		var inserted *mongo.InsertOneResult
		inserted, err = s.dealCollection.InsertOne(ctx, deal)
		if err != nil {
			err = apperrors.Wrap(err, "error inserting new deal")
			return
		}
		primObjectID, ok = inserted.InsertedID.(primitive.ObjectID)
		if !ok {
			err = apperrors.New("can not convert InsertedID to ObjectID inserting deal")
			return
		}
		dealId := primObjectID.Hex()
		deal.Id = &dealId
	} else {
		primObjectID, err = primitive.ObjectIDFromHex(*deal.Id)
		if err != nil {
			err = apperrors.Wrap(err, "can not convert Object ID from Hex updating deal")
			return
		}
		filter := bson.D{{"_id", primObjectID}}
		deal.UpdatedTime = time.Now()

		var updateDoc []byte
		updateDoc, err = bson.Marshal(deal)
		if err != nil {
			err = apperrors.Wrap(err, "can not marshal deal document to update")
			return
		}
		var m bson.M
		err = bson.Unmarshal(updateDoc, &m)
		if err != nil {
			err = apperrors.Wrap(err, "can not unmarshal deal update document")
			return
		}
		delete(m, "_id")
		update := bson.D{{"$set", m}}

		_, err = s.dealCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			err = apperrors.Wrap(err, "error updating deal")
			return
		}
	}

	return
}

func (s *Storage) GetActiveDealByTimeframe(ctx context.Context, timeFrame _struct.Timeframe) (deal *entity.Deal, err error) {
	err = s.dealCollection.FindOne(ctx, bson.D{{
		"$and", bson.A{
			bson.D{{"Timeframe", timeFrame.Resolution()}},
			bson.D{{"Status", bson.D{{"$in", bson.A{entity.DealStatusNew, entity.DealStatusActive}}}}},
		},
	}}, options.FindOne().SetSort(
		bson.D{{"CreatedTime", -1}},
	)).Decode(&deal)
	if errors.Is(err, mongo.ErrNoDocuments) {
		err = nil
	}
	if err != nil {
		err = apperrors.Wrap(err, "error find one, resolution: %s", timeFrame.Resolution())
		return
	}

	if deal == nil {
		deal = &entity.Deal{
			Timeframe:             timeFrame.Resolution(),
			Status:                entity.DealStatusNew,
			MinPercentRangeToSell: timeFrame.MinPercentRangeToSell(),
		}
		if timeFrame.IsHeap() {
			deal.IsHeap = true
		}
		err = s.SaveDeal(ctx, deal)
		if err != nil {
			err = apperrors.Wrap(err, "error save deal")
			return
		}
	}

	// @TODO: Remove it after fill the values in DB
	if deal.MinPercentRangeToSell == 0 {
		deal.MinPercentRangeToSell = timeFrame.MinPercentRangeToSell()
	}
	// @TODO: Remove it after fill the values in DB
	if timeFrame.IsHeap() {
		deal.IsHeap = true
	}

	return
}

func (s *Storage) GetDealsByPeriod(ctx context.Context, from time.Time, to time.Time) (dealPointers []*entity.Deal, err error) {
	var cursor *mongo.Cursor

	cursor, err = s.dealCollection.Find(ctx, bson.D{{
		"$and", bson.A{
			bson.D{{"ClosedTime", bson.D{{"$gt", from}}}},
			bson.D{{"ClosedTime", bson.D{{"$lt", to}}}},
		},
	}}, options.Find().SetSort(
		bson.D{{"ClosedTime", -1}},
	))
	if err != nil {
		err = apperrors.Wrap(err, "error getting cursor deals by period")
		return
	}
	defer cursor.Close(ctx)

	var deals []entity.Deal
	err = cursor.All(ctx, &deals)
	if err != nil {
		err = apperrors.Wrap(err, "error getting all deals by period")
		return
	}
	for _, deal := range deals {
		dealPointers = append(dealPointers, &deal)
	}

	return
}

func (s *Storage) GetDealsOnHeap(ctx context.Context) (dealPointers []*entity.Deal, err error) {
	var cursor *mongo.Cursor

	cursor, err = s.dealCollection.Find(ctx, bson.D{{
		"$and", bson.A{
			bson.D{{"ClosedTime", bson.D{{"$gt", time.Time{}}}}},
			bson.D{{"IsHeap", true}},
		},
	}}, options.Find().SetSort(
		bson.D{{"ClosedTime", -1}},
	))
	if err != nil {
		err = apperrors.Wrap(err, "error getting cursor deals for heap")
		return
	}
	defer cursor.Close(ctx)

	var deals []entity.Deal
	err = cursor.All(ctx, &deals)
	if err != nil {
		err = apperrors.Wrap(err, "error getting all deals for heap")
		return
	}
	for _, deal := range deals {
		dealPointers = append(dealPointers, &deal)
	}

	return
}

func (s *Storage) GetActiveDeals(ctx context.Context) (dealPointers []*entity.Deal, err error) {
	var cursor *mongo.Cursor

	cursor, err = s.dealCollection.Find(ctx, bson.D{{
		"$and", bson.A{
			bson.D{{"Status", entity.DealStatusActive}},
		},
	}}, options.Find().SetSort(
		bson.D{{"UpdatedTime", -1}},
	))
	if err != nil {
		err = apperrors.Wrap(err, "error getting cursor active deals")
		return
	}
	defer cursor.Close(ctx)

	var deals []entity.Deal
	err = cursor.All(ctx, &deals)
	if err != nil {
		err = apperrors.Wrap(err, "error getting all active deals")
		return
	}
	for _, deal := range deals {
		dealPointers = append(dealPointers, &deal)
	}

	return
}

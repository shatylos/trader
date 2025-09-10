package storage

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Order struct {
	Id          *string   `bson:"_id,omitempty"`
	DealId      string    `bson:"DealId"`
	Timeframe   string    `bson:"Timeframe"`
	Price       float64   `bson:"Price"`
	Qty         float64   `bson:"Qty"`
	Side        string    `bson:"Side"`
	CreatedTime time.Time `bson:"CreatedTime"`
}

func (o *Order) Save() (err error) {
	return
}

func (s *Storage) GetOrdersByDealId(dealId string, orders []Order) (err error) {
	ctx := context.Background()

	var cursor *mongo.Cursor
	cursor, err = s.dealCollection.Find(ctx,
		bson.D{{"DealId", dealId}},
		options.Find().SetSort(bson.D{{"CreatedTime", -1}}),
	)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, orders)
	if err != nil {
		return
	}

	return
}

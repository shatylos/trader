package mongo

import (
	"context"
	"github.com/shatylos/trader/internal/strategy/fibonacci/structs"
	"github.com/shatylos/trader/tools"
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
		position.CreatedTime = time.Now().Unix()
		position.UpdatedTime = time.Now().Unix()
		var inserted *mongo.InsertOneResult
		inserted, err = s.positionCollection.InsertOne(ctx, position)
		if err != nil {
			err = tools.AppError{
				Message:     "Error inserting new position",
				ParentError: err,
			}
		}
		primObjectID, ok = inserted.InsertedID.(primitive.ObjectID)
		if !ok {
			err = tools.AppError{Message: "Can not convert InsertedID to ObjectID"}
			return
		}
		savedId = primObjectID.Hex()
	} else {
		primObjectID, err = primitive.ObjectIDFromHex(*position.Id)
		if err != nil {
			return
		}
		filter := bson.D{{"_id", primObjectID}}

		update := bson.D{{"$set", bson.D{
			{"FibonacciChart", position.FibonacciChart},
			{"Trend", position.Trend},
			{"Orders", position.Orders},
			{"UpdatedTime", time.Now().Unix()},
			{"Status", position.Status},
		}}}

		_, err = s.positionCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			err = tools.AppError{
				Message:     "Error updating position",
				ParentError: err,
			}
		}
		savedId = *position.Id
	}
	savedPosition, err = s.GetPositionById(savedId)

	return
}

func (s *MongoStorage) GetLastInternalPosition() (position structs.Position, err error) {
	ctx := context.Background()
	err = s.positionCollection.FindOne(ctx,
		bson.D{{}},
		options.FindOne().SetSort(
			bson.D{{"CreatedTime", -1}},
		)).Decode(&position)
	return
}

func (s *MongoStorage) GetPositionById(idString string) (position structs.Position, err error) {
	var id primitive.ObjectID
	id, err = primitive.ObjectIDFromHex(idString)
	if err != nil {
		return
	}
	ctx := context.Background()
	err = s.positionCollection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&position)
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
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &positions)

	return
}

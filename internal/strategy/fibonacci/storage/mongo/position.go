package mongo

import (
	"context"
	"fmt"
	"github.com/shatylos/trader/internal/strategy/fibonacci/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (s *MongoStorage) SaveInternalPosition(position structs.Position) (err error) {

	ctx := context.Background()

	if position.Id == nil {
		position.CreatedTime = time.Now().Unix()
		position.UpdatedTime = time.Now().Unix()
		_, err = s.positionCollection.InsertOne(ctx, position)
	} else {
		var id primitive.ObjectID
		id, err = primitive.ObjectIDFromHex(*position.Id)
		if err != nil {
			return err
		}
		filter := bson.D{{"_id", id}}

		update := bson.D{{"$set", bson.D{
			{"FibonacciChart", position.FibonacciChart},
			{"Trend", position.Trend},
			{"Orders", position.Orders},
			{"FullQty", position.FullQty},
			{"UpdatedTime", time.Now().Unix()},
			{"Status", position.Status},
		}}}

		var updateResult *mongo.UpdateResult
		updateResult, err = s.positionCollection.UpdateOne(ctx, filter, update)
		fmt.Println(updateResult)
	}

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

package storage

import (
	appStorage "github.com/shatylos/trader/internal/storage"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	"github.com/shatylos/trader/tools/apperrors"
	"go.mongodb.org/mongo-driver/mongo"
)

type Storage struct {
	Id                  string
	orderCollection     *mongo.Collection
	dealCollection      *mongo.Collection
	assetCollection     *mongo.Collection
	activeDealRelations map[string]*entity.DealRelation
}

func (s *Storage) InitStorage() (err error) {
	var db *mongo.Database
	db, err = appStorage.GetDocumentDB()
	if err != nil {
		err = apperrors.Wrap(err, "error get document DB: %s", s.Id)
		return
	}
	s.orderCollection = db.Collection("investor_" + s.Id + "_orders")
	s.dealCollection = db.Collection("investor_" + s.Id + "_deals")
	s.assetCollection = db.Collection("investor_" + s.Id + "_assets")
	s.activeDealRelations = make(map[string]*entity.DealRelation)
	return
}

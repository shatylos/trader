package storage

import (
	"bitbucket.org/shatylos/trader/utils"
	"context"
	"github.com/FerretDB/FerretDB/ferretdb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strings"
)

var documentStorage *mongo.Client

func GetDocumentDB() (*mongo.Database, error) {
	if documentStorage == nil {
		err := InitStorage()
		if err != nil {
			return nil, err
		}
	}
	return documentStorage.Database("trader"), nil
}

func InitStorage() error {
	var err error
	documentStorage, err = getFerretdbStorage()
	if err != nil {
		return err
	}
	return nil
}

func getFerretdbStorage() (*mongo.Client, error) {

	ferretdbListenerTcp := utils.AppConfig("FERRETDB_LISTENER_TCP")
	ferretdbListenerHandler := utils.AppConfig("FERRETDB_HANDLER")
	ferretdbSqliteUrl := utils.AppConfig("FERRETDB_SQLITE_URL")

	dirPath := strings.TrimPrefix(ferretdbSqliteUrl, "file:")
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	fdb, err := ferretdb.New(&ferretdb.Config{
		Listener: ferretdb.ListenerConfig{
			TCP: ferretdbListenerTcp,
		},
		Handler:   ferretdbListenerHandler,
		SQLiteURL: ferretdbSqliteUrl,
	})
	if err != nil {
		return nil, err
	}

	//ctx, cancel := context.WithCancel(context.Background())
	ctx := context.Background()

	done := make(chan struct{})

	go func() {
		err := fdb.Run(ctx)
		if err != nil {
			utils.LogError(err.Error())
		}
		close(done)
	}()

	uri := fdb.MongoDBURI()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		utils.LogError(err.Error())
		//cancel()
		<-done
		return nil, err
	}

	return mongoClient, nil
}

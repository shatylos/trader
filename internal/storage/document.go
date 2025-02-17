package storage

import (
	"context"
	"github.com/FerretDB/FerretDB/ferretdb"
	"github.com/shatylos/trader/internal/config"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strings"
)

var documentStorage *mongo.Client

func GetDocumentDB() (*mongo.Database, error) {
	if documentStorage == nil {
		var err error

		appConfig, err := config.GetConfig()
		if err != nil {
			return nil, err
		}
		mongodbUri := appConfig.App["mongodb_uri"]

		ctx := context.Background()
		documentStorage, err = mongo.Connect(ctx, options.Client().ApplyURI(mongodbUri))
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}

	}
	return documentStorage.Database("trader"), nil
}

func InitStorage() error {

	var err error

	appConfig, err := config.GetConfig()
	if err != nil {
		return err
	}
	ferretdbListenerTcp := appConfig.App["ferretdb_listener_tcp"]
	ferretdbListenerHandler := appConfig.App["ferretdb_handler"]
	ferretdbSqliteUrl := appConfig.App["ferretdb_sqlite_url"]

	if ferretdbListenerTcp == "" || ferretdbListenerHandler == "" || ferretdbSqliteUrl == "" {
		logger.Error("ferretdb config is not set")
		return tools.AppError{Message: "ferretdb config is not set"}
	}

	dirPath := strings.TrimPrefix(ferretdbSqliteUrl, "file:")
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	fdb, err := ferretdb.New(&ferretdb.Config{
		Listener: ferretdb.ListenerConfig{
			TCP: ferretdbListenerTcp,
		},
		Handler:   ferretdbListenerHandler,
		SQLiteURL: ferretdbSqliteUrl,
	})
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = fdb.Run(ctx)
	if err != nil {
		logger.Error(err.Error())
	}

	return nil
}

package storage

import (
	"context"
	"github.com/FerretDB/FerretDB/ferretdb"
	"github.com/shatylos/trader/internal/config"
	"github.com/shatylos/trader/tools/apperrors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strings"
)

var documentStorage *mongo.Client

func GetDocumentDB() (*mongo.Database, error) {
	if documentStorage == nil {
		var err error
		var appConfig *config.Config
		appConfig, err = config.GetConfig()
		if err != nil {
			err = apperrors.Wrap(err, "error get config")
			return nil, err
		}
		mongodbUri := appConfig.App["mongodb_uri"]

		ctx := context.Background()
		documentStorage, err = mongo.Connect(ctx, options.Client().ApplyURI(mongodbUri))
		if err != nil {
			err = apperrors.Wrap(err, "error connect to mongo, mongodbUri: %s", mongodbUri)
			return nil, err
		}

	}
	return documentStorage.Database("trader"), nil
}

func InitStorage() error {

	var err error
	var appConfig *config.Config

	appConfig, err = config.GetConfig()
	if err != nil {
		err = apperrors.Wrap(err, "error get config")
		return err
	}
	ferretdbListenerTcp := appConfig.App["ferretdb_listener_tcp"]
	ferretdbListenerHandler := appConfig.App["ferretdb_handler"]
	ferretdbSqliteUrl := appConfig.App["ferretdb_sqlite_url"]

	if ferretdbListenerTcp == "" || ferretdbListenerHandler == "" || ferretdbSqliteUrl == "" {
		err = apperrors.New("ferretdb config is not set")
		return err
	}

	dirPath := strings.TrimPrefix(ferretdbSqliteUrl, "file:")
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		err = apperrors.Wrap(err, "error make dir, dirPath: %s", dirPath)
		return err
	}

	var fdb *ferretdb.FerretDB
	fdb, err = ferretdb.New(&ferretdb.Config{
		Listener: ferretdb.ListenerConfig{
			TCP: ferretdbListenerTcp,
		},
		Handler:   ferretdbListenerHandler,
		SQLiteURL: ferretdbSqliteUrl,
	})
	if err != nil {
		err = apperrors.Wrap(err, "error get ferretdb new")
		return err
	}

	ctx := context.Background()
	err = fdb.Run(ctx)
	if err != nil {
		err = apperrors.Wrap(err, "error run ferretdb")
		return err
	}

	return nil
}

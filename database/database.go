package database

import (
	"context"
	"log"

	"github.com/manuelladantas/go-crawler/database/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx    = context.TODO()
	cfg    = config.NewConfig()
	client *mongo.Client
)

func init() {
	var err error
	db := cfg.Database

	opts := options.Client().ApplyURI(db.Uri)
	opts.SetAppName("goCrawler")
	// opts.SetAuth(options.Credential{Username: db.User, Password: db.Pass})

	if client, err = mongo.Connect(ctx, opts); err != nil {
		log.Fatal(err, "deu ruim")
	}

}

func Client() *mongo.Database {
	return client.Database(cfg.Database.Name)
}

func Ping() {
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("[ERROR] ", err.Error())
	}
}

func Disconnect() {
	if client != nil {
		client.Disconnect(ctx)
	}
}

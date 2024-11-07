package database

import (
	"context"
	"log"
	"sync"
	"time"
	"ws-chat/configs"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var (
	once       sync.Once
	DbInstance *mongo.Client
)

func DbConn(pctx context.Context, cfg *configs.Config) *mongo.Client {
	once.Do(func() {
		ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
		defer cancel()

		var err error
		DbInstance, err = mongo.Connect(options.Client().ApplyURI(cfg.Db.URI))

		if err != nil {
			log.Fatalf("Error: Connect to database error: %s", err.Error())
		}

		// Ping to ensure that the connection is established
		if err := DbInstance.Ping(ctx, readpref.Primary()); err != nil {
			log.Fatalf("Error: Pinging to database error: %s", err.Error())
		}

		log.Println("Successfully connected to the database.")

		log.Println("Indexes created successfully.")
	})

	return DbInstance
}

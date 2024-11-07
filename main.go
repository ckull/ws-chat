package main

import (
	"context"
	"log"
	"os"
	"ws-chat/configs"
	server "ws-chat/pkg/Server"
	"ws-chat/pkg/database"
)

func main() {
	ctx := context.Background()

	cfg := configs.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path is required")
		}
		return os.Args[1]
	}())

	db := database.DbConn(ctx, cfg)

	server.Start(ctx, cfg, db)

}

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load configurations %v", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("failed to make a connection %v", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	fmt.Println("config is ", len(config.TokenSymmetricKey))
	if err != nil {
		log.Fatal("cannot create server :", err)
	}
	serverErr := server.StartServer(config.ServerAddress)
	if serverErr != nil {
		log.Fatalf("Failed to start server %v", serverErr)
	}
}

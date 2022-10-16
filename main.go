package main

import (
	"database/sql"
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
	server := api.NewServer(store)
	serverErr := server.StartServer(config.ServerAddress)
	if serverErr != nil {
		log.Fatalf("Failed to start server %v", serverErr)
	}
}

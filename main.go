package main

import (
	"database/sql"
	"log"

	"github.com/kaayce/zen-bank/api"
	db "github.com/kaayce/zen-bank/db/sqlc"
	"github.com/kaayce/zen-bank/utils"
	_ "github.com/lib/pq"
)

func main() {
	var err error

	// load env
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	conn, err := sql.Open("postgres", config.DBUrl)
	if err != nil {
		log.Fatalf("cannot connect to the database: %v", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create server: %v", err)
	}
	err = server.Start(config.ServerPort)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}

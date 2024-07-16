package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/kaayce/zen-bank/api"
	db "github.com/kaayce/zen-bank/db/sqlc"
	"github.com/kaayce/zen-bank/utils"
	_ "github.com/lib/pq"
)

func main() {
	var err error

	err = utils.LoadEnv()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	dbUrl := os.Getenv("DB_URL")
	port := os.Getenv("PORT")

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("cannot connect to the database: %v", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(port)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}

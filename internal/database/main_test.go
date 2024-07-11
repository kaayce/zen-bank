package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	var err error

	// err = utils.LoadEnv()
	// if err != nil {
	// 	log.Fatalf("error loading .env file: %v", err)
	// }

	// dbUrl := os.Getenv("DB_URL")
	dbUrl := "postgres://root:secret@localhost:5434/zen_bank?sslmode=disable"

	testDB, err = sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Cannot connect to the database: %v", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}

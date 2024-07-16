package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/kaayce/zen-bank/utils"
	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	var err error

	err = utils.LoadEnv()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	dbUrl := os.Getenv("DB_URL")

	testDB, err = sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Cannot connect to the database: %v", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}

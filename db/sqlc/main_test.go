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

	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	testDB, err = sql.Open("postgres", config.DBSource)
	if err != nil {
		log.Fatalf("Cannot connect to the database: %v", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}

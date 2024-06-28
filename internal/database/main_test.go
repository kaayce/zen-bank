package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	dbUrl := os.Getenv("DB_URL")
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Cannot connect to the database: %v", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}

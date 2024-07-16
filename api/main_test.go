package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	db "github.com/kaayce/zen-bank/db/sqlc"
)

func newTestServer(store db.Store) *Server {
	return NewServer(store)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

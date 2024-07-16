package api

import (
	"github.com/gin-gonic/gin"

	db "github.com/kaayce/zen-bank/internal/database"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// creates new server instance and setup api routes for our service
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// routes
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

// run server at specific address
func (server *Server) Start(address string) error {
	return server.router.Run(":" + address)
}

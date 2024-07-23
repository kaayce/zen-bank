package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	db "github.com/kaayce/zen-bank/db/sqlc"
	"github.com/kaayce/zen-bank/token"
	"github.com/kaayce/zen-bank/utils"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store      db.Store
	tokenMaker token.Maker
	config     *utils.Config
	router     *gin.Engine
}

// creates new server instance and setup api routes for our service
func NewServer(config *utils.Config, store db.Store) (*Server, error) {
	//  create token maker
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// routes
	router.POST("/users", server.createUser)

	router.POST("/accounts", server.createAccount)
	router.POST("/accounts/:id", server.updateAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server, nil
}

// run server at specific address
func (server *Server) Start(address string) error {
	return server.router.Run(":" + address)
}

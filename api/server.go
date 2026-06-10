// Package api provides the HTTP API for the banking service.
package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/santiagot714/SimpleBank/db/sqlc"
	"github.com/santiagot714/SimpleBank/token"
	"github.com/santiagot714/SimpleBank/util"
)

// Server serves HTTP requests for the banking service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{config: config, store: store, tokenMaker: tokenMaker}

	// Custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("currency", validCurrency); err != nil {
			return nil, fmt.Errorf("cannot register currency validator: %w", err)
		}
	}
	// Setup router
	server.setupRouter()

	return server, nil
}

// setupRouter sets up the router for the server.
func (server *Server) setupRouter() {
	router := gin.Default()
	// Routes that do not require authentication
	// User routes
	router.POST("/users", server.createUser)
	router.GET("/users/:username", server.getUser)
	router.POST("/users/login", server.loginUser)

	// Routes that require authentication
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

// Start starts the server on the given address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse returns a JSON response with the error message.
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

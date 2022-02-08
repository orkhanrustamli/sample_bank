package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/orkhanrustamli/simplebank/db/sqlc"
	"github.com/orkhanrustamli/simplebank/token"
	"github.com/orkhanrustamli/simplebank/util"
)

type Server struct {
	store        db.Store
	router       *gin.Engine
	tokenManager token.Manager
	config       util.Config
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenManager, err := token.NewPasetoManager(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token manager: %v", err)
	}

	server := &Server{
		store:        store,
		config:       config,
		tokenManager: tokenManager,
	}

	// Register custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRoutes()

	return server, nil
}

func (server *Server) setupRoutes() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/login", server.login)

	// Routes that require Authentication
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenManager))

	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

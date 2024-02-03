package api

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/dpomian/gobind/db/sqlc"
	"github.com/dpomian/gobind/token"
	"github.com/dpomian/gobind/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type Server struct {
	config     utils.Config
	storage    db.Storage
	tokenMaker token.TokenMaker
	router     *gin.Engine
}

func NewServer(config utils.Config, storage db.Storage) (*Server, error) {
	tokenMaker, err := token.NewPasetoTokenMaker(config.TokenSymmetricKey) // TODO: add symmetric key as an ENV var
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %w", err)
	}

	var server = &Server{
		config:     config,
		storage:    storage,
		tokenMaker: tokenMaker,
	}

	server.configureRoutes()
	return server, nil
}

func newTestServer(t *testing.T, storage db.Storage) *Server {
	config := utils.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, storage)
	require.NoError(t, err)

	return server
}

func (server *Server) configureRoutes() {
	router := gin.Default()
	router.Use(corsMiddleware())

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	notebookApiHandler := NewNotebooksHandler(server.storage, context.Background())
	authRoutes.GET("/api/v1/notebooks", notebookApiHandler.ListNotebooksHandler)
	authRoutes.GET("/api/v1/notebooks/:id", notebookApiHandler.ListNotebookByIdHandler)
	authRoutes.POST("/api/v1/notebooks", notebookApiHandler.AddNewNotebookHandler)
	authRoutes.PUT("/api/v1/notebooks/:id", notebookApiHandler.UpdateNotebookHandler)
	authRoutes.GET("/api/v1/notebooks/search", notebookApiHandler.SearchNotebookHandler)
	authRoutes.DELETE("/api/v1/notebooks/:id", notebookApiHandler.DeleteNotebookHandler)
	authRoutes.GET("/api/v1/notebooks/topics", notebookApiHandler.ListTopicsHandler)

	userApiHandler := NewUserHander(server.config, server.tokenMaker, server.storage, context.Background())
	router.POST("/api/v1/users", userApiHandler.AddNewUserHandler)
	router.POST("/api/v1/users/login", userApiHandler.LoginUser)

	tokenApiHandler := NewTokenHandler(server.config, server.tokenMaker, server.storage, context.Background())
	router.POST("/api/v1/tokens/renew_access", tokenApiHandler.RenewAccessToken)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

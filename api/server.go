package api

import (
	"context"
	"fmt"

	db "github.com/dpomian/gobind/db/sqlc"
	"github.com/dpomian/gobind/token"
	"github.com/dpomian/gobind/utils"
	"github.com/gin-gonic/gin"
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

func (server *Server) configureRoutes() {
	router := gin.Default()

	notebookApiHandler := NewNotebooksHandler(server.storage, context.Background())

	router.GET("/api/v1/notebooks", notebookApiHandler.ListNotebooksHandler)
	router.GET("/api/v1/notebooks/:id", notebookApiHandler.ListNotebookByIdHandler)
	router.POST("/api/v1/notebooks", notebookApiHandler.AddNewNotebookHandler)
	router.PUT("/api/v1/notebooks/:id", notebookApiHandler.UpdateNotebookHandler)
	router.GET("/api/v1/notebooks/search", notebookApiHandler.SearchNotebookHandler)
	router.DELETE("/api/v1/notebooks/:id", notebookApiHandler.DeleteNotebookHandler)

	userApiHandler := NewUserHander(server.config, server.tokenMaker, server.storage, context.Background())

	router.POST("/api/v1/users", userApiHandler.AddNewUserHandler)
	router.POST("/api/v1/users/login", userApiHandler.LoginUser)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

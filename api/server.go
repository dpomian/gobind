package api

import (
	"context"

	db "github.com/dpomian/gobind/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	storage db.Storage
	router  *gin.Engine
}

func NewServer(storage db.Storage) (*Server, error) {
	var server = &Server{
		storage: storage,
	}

	server.configureRoutes()
	return server, nil
}

func (server *Server) configureRoutes() {
	router := gin.Default()

	apiHandler := NewNotebooksHandler(server.storage, context.Background())

	router.GET("/api/v1/notebooks", apiHandler.ListNotebooksHandler)
	router.GET("/api/v1/notebooks/:id", apiHandler.ListNotebookByIdHandler)
	router.POST("/api/v1/notebooks", apiHandler.AddNewNotebookHandler)
	router.PUT("/api/v1/notebooks/:id", apiHandler.UpdateNotebookHandler)
	router.GET("/api/v1/notebooks/search", apiHandler.SearchNotebookHandler)
	router.DELETE("/api/v1/notebooks/:id", apiHandler.DeleteNotebookHandler)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

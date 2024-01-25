package main

import (
	"github.com/dpomian/gobind/handlers"
	"github.com/gin-gonic/gin"
)

var apiHandler *handlers.NotebooksHandler

func init() {
	apiHandler = handlers.NewNotebooksHandler()
}

func main() {
	router := gin.Default()

	router.GET("/api/v1/notebooks", apiHandler.ListNotebooksHandler)
	router.GET("/api/v1/notebooks/:id", apiHandler.ListNotebookByIdHandler)
	router.POST("/api/v1/notebooks", apiHandler.AddNewNotebookHandler)
	router.PUT("/api/v1/notebooks/:id", apiHandler.UpdateNotebookHandler)
	router.GET("/api/v1/notebooks/search", apiHandler.SearchNotebookHandler)

	router.Run(":5050")
}

package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/dpomian/gobind/handlers"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var apiHandler *handlers.NotebooksHandler

func init() {
	var dbSource = os.Getenv("BINDER_DB_SOURCE")
	var dbDriver = os.Getenv("BINDER_DB_DRIVER")
	db, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("cannot connect do db:", err)
	}

	apiHandler = handlers.NewNotebooksHandler(db, context.Background())
}

func main() {
	router := gin.Default()

	router.GET("/api/v1/notebooks", apiHandler.ListNotebooksHandler)
	router.GET("/api/v1/notebooks/:id", apiHandler.ListNotebookByIdHandler)
	router.POST("/api/v1/notebooks", apiHandler.AddNewNotebookHandler)
	router.PUT("/api/v1/notebooks/:id", apiHandler.UpdateNotebookHandler)
	router.GET("/api/v1/notebooks/search", apiHandler.SearchNotebookHandler)
	router.DELETE("/api/v1/notebooks/:id", apiHandler.DeleteNotebookHandler)

	router.Run(":5050")
}

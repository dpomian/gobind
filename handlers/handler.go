package handlers

import (
	"fmt"
	"net/http"

	"github.com/dpomian/gobind/models"
	"github.com/gin-gonic/gin"
)

type NotebooksHandler struct {
	// ctx context.Context
}

func NewNotebooksHandler() *NotebooksHandler {
	return &NotebooksHandler{}
}

func (handler *NotebooksHandler) ListNotebooksHandler(c *gin.Context) {
	notebookList := []string{}
	c.JSON(http.StatusOK, notebookList)
}

func (handler *NotebooksHandler) ListNotebookByIdHandler(c *gin.Context) {
	notebookId := c.Param("id")
	fmt.Println("notebookId:", notebookId)

	c.JSON(http.StatusNotFound, gin.H{
		"message": "notebook not found",
	})
}

func (handler *NotebooksHandler) AddNewNotebookHandler(c *gin.Context) {
	var reqNotebook models.Notebook
	if err := c.ShouldBindJSON(&reqNotebook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Println("req notebook:", reqNotebook)

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "not implemented",
	})
}

func (handler *NotebooksHandler) UpdateNotebookHandler(c *gin.Context) {
	reqNotebookId := c.Param("id")
	var reqNotebook models.Notebook
	if err := c.ShouldBindJSON(&reqNotebook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	fmt.Println("reqNotebookId:", reqNotebookId)
	fmt.Println("reqNotebook:", reqNotebook)

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "not implemented",
	})
}

func (handler *NotebooksHandler) SearchNotebookHandler(c *gin.Context) {
	searchBy := c.Query("text")
	fmt.Println("searchBy:", searchBy)

	notebooks := []models.Notebook{}

	c.JSON(http.StatusOK, notebooks)
}

package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/dpomian/gobind/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	MiscTopic string = "Misc"
)

var (
	InternalError    = gin.H{"error": "internal error"}
	InvalidId        = gin.H{"error": "invalid notebook id"}
	NotebookNotFound = gin.H{"error": "notebook not found"}
	NotImplemented   = gin.H{"error": "not implemented"}
	NotebookDeleted  = gin.H{"message": "notebook deleted"}
)

type NotebooksHandler struct {
	db  *sql.DB
	ctx context.Context
}

func NewNotebooksHandler(db *sql.DB, ctx context.Context) *NotebooksHandler {
	return &NotebooksHandler{
		db:  db,
		ctx: ctx,
	}
}

func (handler *NotebooksHandler) ListNotebooksHandler(c *gin.Context) {
	limit := 50
	offset := 0

	query := db.New(handler.db)

	arg := db.ListNotebooksParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	notebooks, err := query.ListNotebooks(handler.ctx, arg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, InternalError)
	}

	if notebooks == nil {
		notebooks = []db.Notebook{}
	}

	c.JSON(http.StatusOK, notebooks)
}

type listNotebookByIdRq struct {
	Id string `uri:"id" binding:"required"`
}

func (handler *NotebooksHandler) ListNotebookByIdHandler(c *gin.Context) {
	var rq listNotebookByIdRq
	if err := c.ShouldBindUri(&rq); err != nil {
		c.JSON(http.StatusBadRequest, InvalidId)
		return
	}

	rqNotebookId, err := uuid.Parse(rq.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, InvalidId)
		return
	}

	query := db.New(handler.db)
	dbNotebook, err := query.GetNotebook(handler.ctx, rqNotebookId)

	fmt.Println(dbNotebook)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, struct{}{})
		return
	} else if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	c.JSON(http.StatusOK, dbNotebook)
}

type addNewNotebookRq struct {
	Title   string `json:"title" binding:"required"`
	Topic   string `json:"topic" binding:"-"`
	Content string `json:"content" binding:"-"`
}

func (handler *NotebooksHandler) AddNewNotebookHandler(c *gin.Context) {
	var newNotebook addNewNotebookRq
	if err := c.ShouldBindJSON(&newNotebook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(newNotebook.Topic) == 0 {
		newNotebook.Topic = MiscTopic
	}

	query := db.New(handler.db)
	arg := db.CreateNotebookParams{
		ID:        uuid.New(),
		Title:     newNotebook.Title,
		Topic:     newNotebook.Topic,
		Content:   newNotebook.Content,
		CreatedAt: time.Now(),
	}

	dbNotebook, err := query.CreateNotebook(handler.ctx, arg)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	c.JSON(http.StatusOK, dbNotebook)
}

type updateNotebookRq struct {
	Title   string `json:"title" binding:"required"`
	Topic   string `json:"topic" binding:"-"`
	Content string `json:"content" binding:"-"`
}

func (handler *NotebooksHandler) UpdateNotebookHandler(c *gin.Context) {
	rqNotebookId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, InvalidId)
		return
	}

	var rqNotebook updateNotebookRq
	if err := c.ShouldBindJSON(&rqNotebook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if rqNotebook.Topic == "" {
		rqNotebook.Topic = MiscTopic
	}

	arg := db.UpdateNotebookParams{
		ID:           rqNotebookId,
		Title:        rqNotebook.Title,
		Content:      rqNotebook.Content,
		Topic:        rqNotebook.Topic,
		LastModified: time.Now(),
	}
	query := db.New(handler.db)
	dbNotebook, err := query.UpdateNotebook(handler.ctx, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	c.JSON(http.StatusOK, dbNotebook)
}

func (handler *NotebooksHandler) SearchNotebookHandler(c *gin.Context) {
	searchBy := "%" + c.Query("text") + "%"

	fmt.Println(searchBy)
	query := db.New(handler.db)
	notebooks, err := query.SearchNotebooks(handler.ctx, searchBy)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	c.JSON(http.StatusOK, notebooks)
}

func (handler *NotebooksHandler) DeleteNotebookHandler(c *gin.Context) {
	rqNotebookId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, InvalidId)
		return
	}

	args := db.DeleteNotebookParams{
		ID:           rqNotebookId,
		LastModified: time.Now(),
	}

	query := db.New(handler.db)
	err = query.DeleteNotebook(handler.ctx, args)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	c.JSON(http.StatusOK, NotebookDeleted)
}

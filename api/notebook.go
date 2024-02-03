package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/dpomian/gobind/db/sqlc"
	"github.com/dpomian/gobind/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	Forbidden        = gin.H{"message": "forbidden"}
)

type NotebooksHandler struct {
	storage db.Storage
	ctx     context.Context
}

func NewNotebooksHandler(storage db.Storage, ctx context.Context) *NotebooksHandler {
	return &NotebooksHandler{
		storage: storage,
		ctx:     ctx,
	}
}

func (handler *NotebooksHandler) ListNotebooksHandler(c *gin.Context) {
	limit := 100
	offset := 0

	authPayload := c.MustGet(authPayloadKey).(*token.Payload)

	if topicQuery := c.Query("topic"); len(topicQuery) > 0 {
		handler.GetNotebookTitlesByTopic(c)
		return
	}

	arg := db.ListNotebooksParams{
		UserID: authPayload.UserId,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	notebooks, err := handler.storage.ListNotebooks(handler.ctx, arg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, InternalError)
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

	authPayload := c.MustGet(authPayloadKey).(*token.Payload)

	arg := db.GetNotebookParams{
		ID:     rqNotebookId,
		UserID: authPayload.UserId,
	}
	dbNotebook, err := handler.storage.GetNotebook(handler.ctx, arg)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, NotebookNotFound)
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

	authPayload := c.MustGet(authPayloadKey).(*token.Payload)

	arg := db.CreateNotebookParams{
		ID:        uuid.New(),
		UserID:    authPayload.UserId,
		Title:     newNotebook.Title,
		Topic:     newNotebook.Topic,
		Content:   newNotebook.Content,
		CreatedAt: time.Now(),
	}

	dbNotebook, err := handler.storage.CreateNotebook(handler.ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				c.JSON(http.StatusForbidden, Forbidden)
			}
		} else {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, InternalError)

		}
		return
	}

	c.JSON(http.StatusCreated, dbNotebook)
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

	authPayload := c.MustGet(authPayloadKey).(*token.Payload)

	arg := db.UpdateNotebookParams{
		ID:           rqNotebookId,
		UserID:       authPayload.UserId,
		Title:        rqNotebook.Title,
		Content:      rqNotebook.Content,
		Topic:        rqNotebook.Topic,
		LastModified: time.Now(),
	}
	dbNotebook, err := handler.storage.UpdateNotebook(handler.ctx, arg)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, NotebookNotFound)
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	c.JSON(http.StatusOK, dbNotebook)
}

func (handler *NotebooksHandler) SearchNotebookHandler(c *gin.Context) {
	searchText := c.Query("text")

	if len(searchText) == 0 {
		handler.ListNotebooksHandler(c)
		return
	}

	searchBy := "%" + searchText + "%"
	authPayload := c.MustGet(authPayloadKey).(*token.Payload)

	arg := db.SearchNotebooksParams{
		UserID: authPayload.UserId,
		Title:  searchBy,
	}
	notebooks, err := handler.storage.SearchNotebooks(handler.ctx, arg)
	if err != nil {
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

	authPayload := c.MustGet(authPayloadKey).(*token.Payload)

	args := db.DeleteNotebookParams{
		ID:           rqNotebookId,
		UserID:       authPayload.UserId,
		LastModified: time.Now(),
	}

	_, err = handler.storage.DeleteNotebook(handler.ctx, args)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, NotebookNotFound)
		return
	} else if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	c.JSON(http.StatusOK, NotebookDeleted)
}

type rsTopicList struct {
	Topics []string `json:"topics"`
}

func (handler *NotebooksHandler) ListTopicsHandler(c *gin.Context) {
	authPayload := c.MustGet(authPayloadKey).(*token.Payload)

	arg := db.ListTopicsParams{
		UserID: authPayload.UserId,
		Topic:  "%" + c.Query("topic") + "%",
	}
	topics, err := handler.storage.ListTopics(handler.ctx, arg)

	if err != nil && err != sql.ErrNoRows {
		fmt.Println("error fetching topics:", err)
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	c.JSON(http.StatusOK, rsTopicList{Topics: topics})
}

func (handler *NotebooksHandler) GetNotebookTitlesByTopic(c *gin.Context) {
	authPayload := c.MustGet(authPayloadKey).(*token.Payload)
	topicQuery := c.Query("topic")

	arg := db.GetNotebookTitlesByTopicParams{
		UserID: authPayload.UserId,
		Topic:  topicQuery,
	}
	notebookTitles, err := handler.storage.GetNotebookTitlesByTopic(handler.ctx, arg)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("error fetching notebooks:", err)
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	fmt.Println("api data:", notebookTitles)

	c.JSON(http.StatusOK, notebookTitles)
}

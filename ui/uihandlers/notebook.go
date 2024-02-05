package uihandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/dpomian/gobind/ui/httputils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type rsNotebook struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Topic        string    `json:"topic"`
	Content      string    `json:"content"`
	Deleted      bool      `json:"deleted"`
	LastModified time.Time `json:"last_modified"`
	CreatedAt    time.Time `json:"created_at"`
	UserID       uuid.UUID `json:"user_id"`
}

type rsNotebookLite struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type notebooksByTitle []rsNotebookLite

type rsNotebookLitePerTopic struct {
	Topic     string           `json:"topic"`
	Notebooks []rsNotebookLite `json:"notebooks"`
}

func (nb notebooksByTitle) Len() int           { return len(nb) }
func (nb notebooksByTitle) Less(i, j int) bool { return nb[i].Title < nb[j].Title }
func (nb notebooksByTitle) Swap(i, j int)      { nb[i], nb[j] = nb[j], nb[i] }

func (handler *RqHandler) HandleNotebookLite(c *gin.Context) {
	const htmlTemplate = "notebooks_accordion.html"
	accessToken := c.GetString(CtxAccessTokenKey)

	rsNotebookPerTopic := []rsNotebookLitePerTopic{}
	notebookPerTopicMap := make(map[string][]rsNotebookLite)
	notebooks := []rsNotebook{}

	url := handler.config.BinderApiBaseUrl + "/api/v1/notebooks/search"
	headers := httputils.NewHeaders().WithAuthorizationHeader("Bearer " + accessToken)
	queries := httputils.NewQueries().Add("text", c.Query("text"))
	response, statusCode, err := httputils.SendRequest(httputils.RequestTypeGet, url, headers, queries, nil)
	if err != nil {
		handler.errorHandler.HandlerHtmlErrorRs(statusCode, htmlTemplate, c)
		return
	}

	if statusCode != http.StatusOK {
		handler.errorHandler.HandlerHtmlErrorRs(statusCode, htmlTemplate, c)
		return
	}

	if err = json.Unmarshal(response, &notebooks); err != nil {
		handler.errorHandler.HandlerHtmlErrorRs(http.StatusInternalServerError, htmlTemplate, c)
		return
	}

	for _, notebook := range notebooks {
		nbLite := rsNotebookLite{ID: notebook.ID, Title: notebook.Title}
		notebookPerTopicMap[notebook.Topic] = append(notebookPerTopicMap[notebook.Topic], nbLite)
	}

	nbKeys := make([]string, 0, len(notebookPerTopicMap))
	for key := range notebookPerTopicMap {
		nbKeys = append(nbKeys, key)
	}
	sort.Strings(nbKeys)

	for _, nbk := range nbKeys {
		sort.Sort(notebooksByTitle(notebookPerTopicMap[nbk]))
		nbLitePerTopic := rsNotebookLitePerTopic{
			Topic:     nbk,
			Notebooks: notebookPerTopicMap[nbk],
		}
		rsNotebookPerTopic = append(rsNotebookPerTopic, nbLitePerTopic)
	}

	c.HTML(statusCode, "notebooks_accordion.html", gin.H{
		"NotebooksPerTopic": rsNotebookPerTopic,
	})
}

type notebookDetailsRs struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Topic        string    `json:"topic"`
	Content      string    `json:"content"`
	Deleted      bool      `json:"deleted"`
	LastModified time.Time `json:"last_modified"`
	CreatedAt    time.Time `json:"created_at"`
	UserID       uuid.UUID `json:"user_id"`
}

type NotebookBriefDetails struct {
	Id      string `json:"notebook_id"`
	Title   string `json:"notebook_title"`
	Topic   string `json:"notebook_topic"`
	Content string `json:"notebook_content"`
}

func (handler *RqHandler) HandleGetNotebookDetails(c *gin.Context) {
	accessToken := c.GetString(CtxAccessTokenKey)

	notebookBriefDetails := NotebookBriefDetails{}
	notebookId := c.Param("id")
	if len(notebookId) > 0 {
		url := handler.config.BinderApiBaseUrl + "/api/v1/notebooks/" + notebookId
		headers := httputils.NewHeaders().
			WithJsonContentTypeHeader().
			WithAuthorizationHeader("Bearer " + accessToken)
		responseData, statusCode, err := httputils.SendGETRequest(url, nil, headers)
		if err != nil {
			fmt.Println("handleGetNotebookDetails err: ", err)
			handler.errorHandler.HandleJsonErrorRs(statusCode, nil, c)
			return
		}

		if statusCode == http.StatusUnauthorized {
			fmt.Println("unauthorized response")
			handler.errorHandler.HandleJsonErrorRs(statusCode, nil, c)
			return
		}

		var rsData notebookDetailsRs
		err = json.Unmarshal(responseData, &rsData)
		if err != nil {
			fmt.Println("error unmarshalling notebookDetailsRs data")
			handler.errorHandler.HandleJsonErrorRs(http.StatusInternalServerError, nil, c)
			return
		}

		notebookBriefDetails = NotebookBriefDetails{
			Id:      rsData.ID.String(),
			Title:   rsData.Title,
			Topic:   rsData.Topic,
			Content: rsData.Content,
		}
	} else {
		fmt.Println("notebook id is empty")
		handler.errorHandler.HandleJsonErrorRs(http.StatusInternalServerError, nil, c)
		return
	}

	c.JSON(http.StatusOK, notebookBriefDetails)
}

type saveNotebookRq struct {
	Id      string `json:"notebook_id"`
	Title   string `json:"notebook_title"`
	Topic   string `json:"notebook_topic"`
	Content string `json:"notebook_content"`
}

type updateNotebookRq struct {
	Title   string `json:"title" binding:"required"`
	Topic   string `json:"topic" binding:"-"`
	Content string `json:"content" binding:"-"`
}

func (handler *RqHandler) HandleSaveNotebook(c *gin.Context) {
	accessToken := c.GetString(CtxAccessTokenKey)
	notebookId := c.Param("id")
	notebookContent := c.PostForm("active_notebook")

	if len(notebookContent) == 0 {
		var rqData saveNotebookRq
		err := c.ShouldBindJSON(&rqData)
		if err != nil {
			fmt.Println("error binding data to json: ", err)
			c.Status(http.StatusBadRequest)
			return
		} else {
			rawData, err := json.Marshal(rqData)
			if err != nil {
				fmt.Println("error marshaling data")
				c.Status(http.StatusBadRequest)
				return
			} else {
				notebookContent = string(rawData)
			}
		}
	}

	var saveRqData saveNotebookRq
	err := json.Unmarshal([]byte(notebookContent), &saveRqData)
	if err != nil {
		fmt.Println("error unmarshalling data: ", err)
		c.Status(http.StatusBadRequest)
		return
	}

	updateNotebook := updateNotebookRq{
		Title:   saveRqData.Title,
		Topic:   saveRqData.Topic,
		Content: saveRqData.Content,
	}

	updateNotebookData, err := json.Marshal(updateNotebook)
	if err != nil {
		fmt.Println("error marshalling data: ", err)
		c.Status(http.StatusBadRequest)
		return
	}

	url := handler.config.BinderApiBaseUrl + "/api/v1/notebooks/" + notebookId
	headers := httputils.NewHeaders().
		WithJsonContentTypeHeader().
		WithAuthorizationHeader("Bearer " + accessToken)
	_, statusCode, err := httputils.SendRequest(httputils.RequestTypePut, url, headers, nil, updateNotebookData)

	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if statusCode == http.StatusUnauthorized {
		fmt.Println("unauthorized request")
		c.Status(http.StatusUnauthorized)
		return
	}

	c.Status(http.StatusOK)
}

func (handler *RqHandler) HandleNotebooksModal(c *gin.Context) {
	c.HTML(http.StatusFound, "new_notebook_form.html", nil)
}

type rqAddNotebook struct {
	Title string `json:"title"`
	Topic string `json:"topic"`
}

func (handler *RqHandler) HandleAddNotebook(c *gin.Context) {
	accessToken := c.GetString(CtxAccessTokenKey)

	rq := rqAddNotebook{
		Title: c.Request.FormValue("title"),
		Topic: c.Request.FormValue("topic"),
	}

	if rq.Title == "" || rq.Topic == "" {
		fmt.Println("error binding input params:", rq)
		return
	}

	jsonData, _ := json.Marshal(rq)

	url := handler.config.BinderApiBaseUrl + "/api/v1/notebooks"
	headers := httputils.NewHeaders().
		WithJsonContentTypeHeader().
		WithAuthorizationHeader("Bearer " + accessToken)

	_, statusCode, err := httputils.SendPOSTRequest(url, headers, jsonData)
	if err != nil {
		fmt.Println("error in http request:", err)
	}

	c.Status(statusCode)
}

type rsTopicList struct {
	Topics []string `json:"topics"`
}

func (handler *RqHandler) HandleGetNotebookTopics(c *gin.Context) {
	topicSubstring := c.Query("topic")
	var topicOptions rsTopicList

	url := handler.config.BinderApiBaseUrl + "/api/v1/notebooks/topics"
	queries := httputils.NewQueries().Add("topic", topicSubstring)
	headers := httputils.NewHeaders().
		WithJsonContentTypeHeader().
		WithAuthorizationHeader("Bearer " + c.GetString(CtxAccessTokenKey))

	responseData, statusCode, err := httputils.SendRequest(httputils.RequestTypeGet, url, headers, queries, nil)
	if err != nil {
		fmt.Println("error retrieving topics: ", err)
		return
	}

	err = json.Unmarshal(responseData, &topicOptions)
	if err != nil {
		fmt.Println("error unmarshalling response: ", err)
		return
	}

	defer func(statusCode int, topicOptions rsTopicList) {
		c.HTML(statusCode, "topic_options.html", gin.H{"TopicOptions": topicOptions.Topics})
		return
	}(statusCode, topicOptions)
}

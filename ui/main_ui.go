package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dpomian/gobind/ui/httputils"
	"github.com/dpomian/gobind/utils"
	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
)

type RqHandler struct {
	ctx         context.Context
	redisClient *redis.Client
}

func NewRqHandler(redisClient *redis.Client) *RqHandler {
	return &RqHandler{
		ctx:         context.Background(),
		redisClient: redisClient,
	}
}

func redirectHome(userLoggedIn bool, statusCode int, c *gin.Context) {
	if statusCode >= http.StatusBadRequest {
		c.HTML(http.StatusOK, "login.html", nil)
		return
	}

	if userLoggedIn {
		c.HTML(statusCode, "app.html", gin.H{
			"userLoggedIn": userLoggedIn,
		})
	} else {
		c.HTML(statusCode, "login.html", nil)
	}
}

func (handler *RqHandler) renderMainApp(c *gin.Context) {
	c.HTML(http.StatusOK, "app_main.html", nil)
}

func (handler *RqHandler) renderLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (handler *RqHandler) renderHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "app.html", nil)
}

func (handler *RqHandler) handleBinderContent(c *gin.Context) {
	accessToken := handler.getAccessToken(c)
	if accessToken == nil {
		fmt.Println("nil access token => renderLogin()")
		handler.renderLogin(c)
		return
	}
	fmt.Println("handleBinderContent: access token:", *accessToken)

	handler.renderMainApp(c)
}

func (handler *RqHandler) handleBinderHomePage(c *gin.Context) {
	handler.renderHomePage(c)
}

type loginRq struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type loginDetails struct {
	SessionId             string    `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func (handler *RqHandler) handlePOSTUserLogin(c *gin.Context) {
	session := sessions.Default(c)
	sessionId := session.Get("session_id")

	fmt.Println("session id:", sessionId)

	if sessionId == nil {
		sessionId = uuid.NewString()
		var loginRequestData loginRq
		if err := c.ShouldBind(&loginRequestData); err != nil {
			fmt.Println(loginRequestData)
			fmt.Println("error binding json", err)
			handler.handleBinderHomePage(c)
			return
		}

		postData, err := json.Marshal(loginRequestData)
		if err != nil {
			fmt.Println("error marshalling data:", err)
			return
		}
		url := "http://localhost:5050/api/v1/users/login"
		headers := httputils.NewHeaders().WithJsonContentTypeHeader()
		responseData, statusCode, err := httputils.SendPOSTRequest(url, headers, postData)

		if err != nil || statusCode >= http.StatusBadRequest {
			fmt.Println("error in Post rq:", err)
			fmt.Println("status:", statusCode)
			handler.handleBinderHomePage(c)
			return
		}

		session.Set("session_id", sessionId)
		session.Save()

		// add sessionId to redis
		handler.redisClient.Set(handler.ctx, sessionId.(string), responseData, time.Duration(24*time.Hour))

		fmt.Println(string(responseData))
	} else {
		_, err := handler.redisClient.Get(handler.ctx, sessionId.(string)).Result()
		if err != nil {
			fmt.Println("sessionId:", sessionId, "does not exist in cache, invalidate session")
			utils.InvalidateSession(session, c)
			handler.handleBinderHomePage(c)
		}
	}

	handler.handleBinderHomePage(c)
}

type rqRenewAccessToken struct {
	RefreshToken string `json:"refresh_token"`
}

type rsRenewAccesstoken struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (handler *RqHandler) getAccessToken(c *gin.Context) *string {
	session := sessions.Default(c)
	sessionId := session.Get("session_id")

	if sessionId == nil {
		fmt.Println("session_id not in request")
		return nil
	}

	fmt.Println("sessionId:", sessionId)

	authData, err := handler.redisClient.Get(handler.ctx, sessionId.(string)).Result()
	if err == redis.Nil {
		fmt.Println("session_id: ", sessionId, "not found in redis: ", err)
		utils.InvalidateSession(session, c)
		return nil
	}

	// fmt.Println("authData:", authData)
	var lgDetails loginDetails
	if len(authData) > 0 {
		if err := json.Unmarshal([]byte(authData), &lgDetails); err == nil {
			// fmt.Println(lgDetails)
			// check if access token is expired
			if time.Now().After(lgDetails.AccessTokenExpiresAt) {
				fmt.Println("access token expired try to renew it")

				// check if refresh token is expired
				if time.Now().After(lgDetails.RefreshTokenExpiresAt) {
					fmt.Println("refresh token expired")
					return nil
				} else {
					postData, err := json.Marshal(rqRenewAccessToken{RefreshToken: lgDetails.RefreshToken})
					url := "http://localhost:5050/api/v1/tokens/renew_access"
					headers := httputils.NewHeaders().WithJsonContentTypeHeader()
					responseData, statusCode, err := httputils.SendPOSTRequest(url, headers, postData)

					if statusCode == http.StatusUnauthorized {
						handler.handleUnauthorized(c)
						return nil
					}

					var renewAccesTokenResponse rsRenewAccesstoken
					err = json.Unmarshal(responseData, &renewAccesTokenResponse)
					if err != nil {
						fmt.Println("error unmarshalling response data", err)
						return nil
					}

					lgDetails.AccessToken = renewAccesTokenResponse.AccessToken
					lgDetails.AccessTokenExpiresAt = renewAccesTokenResponse.AccessTokenExpiresAt

					marshaledTokens, _ := json.Marshal(lgDetails)

					handler.redisClient.Set(handler.ctx, sessionId.(string), marshaledTokens, time.Duration(24*time.Hour))
				}
			} else {
				return &lgDetails.AccessToken
			}
		} else {
			fmt.Println("error unmarshalling authData")
			return nil
		}
	}

	return &lgDetails.AccessToken
}

func (handler *RqHandler) handleUserLogout(c *gin.Context) {
	session := sessions.Default(c)
	sessionId := session.Get("session_id")
	fmt.Println("session id:", sessionId)

	utils.InvalidateSession(session, c)

	if sessionId != nil {
		handler.redisClient.Del(handler.ctx, sessionId.(string))
	}

	handler.renderLogin(c)
}

type NotebookBriefDetails struct {
	Id      string `json:"notebook_id"`
	Title   string `json:"notebook_title"`
	Topic   string `json:"notebook_topic"`
	Content string `json:"notebook_content"`
}

type NotebookIdAndtitle struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func (handler *RqHandler) handleGetNotebookTitles(c *gin.Context) {
	accessToken := handler.getAccessToken(c)
	if accessToken == nil {
		handler.handleBinderHomePage(c)
		return
	}

	notebookBriefs := []NotebookIdAndtitle{}
	rqTopic := c.Query("topic")

	headers := httputils.NewHeaders().Add("Authorization", "Bearer "+*accessToken)
	queries := httputils.NewQueries().Add("topic", rqTopic)

	url := "http://localhost:5050/api/v1/notebooks"
	response, statusCode, err := httputils.SendGETRequest(url, queries, headers)

	if statusCode == http.StatusUnauthorized {
		handler.handleUnauthorized(c)
		return
	}

	if err != nil {
		fmt.Println(err)
		c.HTML(http.StatusOK, "list.html", gin.H{
			"NotebookIdAndtitleList": notebookBriefs,
		})
		return
	}

	if err = json.Unmarshal(response, &notebookBriefs); err != nil {
		fmt.Println("error unmarshalling response body:", err)
		fmt.Println("response:", string(response))
		c.HTML(http.StatusOK, "list.html", gin.H{
			"NotebookIdAndtitleList": notebookBriefs,
		})
		return
	}

	c.HTML(http.StatusOK, "list.html", gin.H{
		"NotebookIdAndtitleList": notebookBriefs,
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

func (handler *RqHandler) handleGetNotebookDetails(c *gin.Context) {
	accessToken := handler.getAccessToken(c)
	if accessToken == nil {
		handler.renderLogin(c)
		return
	}

	notebookBriefDetails := NotebookBriefDetails{}
	notebookId := c.Param("id")
	if len(notebookId) > 0 {
		url := "http://localhost:5050/api/v1/notebooks/" + notebookId
		headers := httputils.NewHeaders().
			WithJsonContentTypeHeader().
			WithAuthorizationHeader("Bearer " + *accessToken)
		responseData, statusCode, err := httputils.SendGETRequest(url, nil, headers)
		if err != nil {
			fmt.Println("handleGetNotebookDetails err: ", err)
			c.JSON(http.StatusBadRequest, NotebookBriefDetails{})
			return
		}

		if statusCode == http.StatusUnauthorized {
			fmt.Println("unauthorized response")
			c.JSON(http.StatusBadRequest, NotebookBriefDetails{})
			return
		}

		var rsData notebookDetailsRs
		err = json.Unmarshal(responseData, &rsData)
		if err != nil {
			fmt.Println("error unmarshalling notebookDetailsRs data")
			c.JSON(http.StatusInternalServerError, NotebookBriefDetails{})
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
		c.JSON(http.StatusInternalServerError, NotebookBriefDetails{})
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

func (handler *RqHandler) handleSaveNotebook(c *gin.Context) {
	accessToken := handler.getAccessToken(c)
	if accessToken == nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	notebookId := c.Param("id")
	notebookContent := c.PostForm("active_notebook")

	if len(notebookContent) == 0 {
		var rqData saveNotebookRq
		err := c.ShouldBindJSON(&rqData)
		if err != nil {
			fmt.Println("error binding data to json: ", err)
		} else {
			rawData, err := json.Marshal(rqData)
			if err != nil {
				fmt.Println("error marshaling data")
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

	url := "http://localhost:5050/api/v1/notebooks/" + notebookId
	headers := httputils.NewHeaders().
		WithJsonContentTypeHeader().
		WithAuthorizationHeader("Bearer " + *accessToken)
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

func (handler *RqHandler) handleNotebooksModal(c *gin.Context) {
	session := sessions.Default(c)
	sessionId := session.Get("session_id")

	if sessionId == nil {
		c.HTML(http.StatusOK, "app.html", gin.H{
			"userLoggedIn": false,
		})
		return
	}

	c.HTML(http.StatusFound, "new_notebook_form.html", nil)
}

type rqAddNotebook struct {
	Title string `json:"title"`
	Topic string `json:"topic"`
}

func (handler *RqHandler) handleAddNotebook(c *gin.Context) {
	accessToken := handler.getAccessToken(c)
	if accessToken == nil {
		handler.handleBinderHomePage(c)
		return
	}

	rq := rqAddNotebook{
		Title: c.Request.FormValue("title"),
		Topic: c.Request.FormValue("topic"),
	}

	if rq.Title == "" || rq.Topic == "" {
		fmt.Println("error binding input params:", rq)
		handler.renderMainApp(c)
		return
	}

	jsonData, _ := json.Marshal(rq)

	url := "http://localhost:5050/api/v1/notebooks"
	headers := httputils.NewHeaders().
		WithJsonContentTypeHeader().
		WithAuthorizationHeader("Bearer " + *accessToken)

	_, statusCode, err := httputils.SendPOSTRequest(url, headers, jsonData)
	if err != nil {
		fmt.Println("error in http request:", err)
	}
	if statusCode == http.StatusUnauthorized {
		fmt.Println("unauthorized request:", err)
		handler.handleBinderHomePage(c)
		return
	}

	handler.renderMainApp(c)
}

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

type rsNotebookLitePerTopic struct {
	Topic     string           `json:"topic"`
	Notebooks []rsNotebookLite `json:"notebooks"`
}

func (handler *RqHandler) handleNotebookLite(c *gin.Context) {
	accessToken := handler.getAccessToken(c)
	rsNotebookPerTopic := []rsNotebookLitePerTopic{}
	notebookPerTopicMap := make(map[string][]rsNotebookLite)
	notebooks := []rsNotebook{}

	defer func() {
		c.HTML(http.StatusOK, "notebooks_accordion.html", gin.H{
			"NotebooksPerTopic": rsNotebookPerTopic,
		})
	}()

	if accessToken != nil {
		url := "http://localhost:5050/api/v1/notebooks/search"
		headers := httputils.NewHeaders().WithAuthorizationHeader("Bearer " + *accessToken)
		queries := httputils.NewQueries().Add("text", c.Query("text"))
		response, statusCode, err := httputils.SendRequest(httputils.RequestTypeGet, url, headers, queries, nil)
		if err != nil {
			fmt.Println("error sending GET request: %w", err)
			return
		}

		if statusCode == http.StatusUnauthorized {
			handler.handleBinderHomePage(c)
			return
		}

		if err = json.Unmarshal(response, &notebooks); err != nil {
			fmt.Println("error unmarshalling response body:", err)
			return
		}

		for _, notebook := range notebooks {
			nbLite := rsNotebookLite{ID: notebook.ID, Title: notebook.Title}
			notebookPerTopicMap[notebook.Topic] = append(notebookPerTopicMap[notebook.Topic], nbLite)
		}

		for nk, nv := range notebookPerTopicMap {
			nbLitePerTopic := rsNotebookLitePerTopic{
				Topic:     nk,
				Notebooks: nv,
			}
			rsNotebookPerTopic = append(rsNotebookPerTopic, nbLitePerTopic)
		}
	}
}

func (handler RqHandler) handleUnauthorized(c *gin.Context) {
	session := sessions.Default(c)
	sessionId := session.Get("session_id")

	fmt.Println(sessionId)
	if sessionId != nil {
		result := handler.redisClient.Del(handler.ctx, sessionId.(string))
		fmt.Println("removed ", result, "keys")
	}

	utils.InvalidateSession(session, c)

	handler.handleBinderHomePage(c)
}

type UIConfig struct {
	REDIS_URI    string `mapstructure:"REDIS_URI"`
	REDIS_SECRET string `mapstructure:"REDIS_SECRET"`
}

func main() {
	uiConfig := readUiConfig()
	store, err := redisStore.NewStore(10, "tcp", uiConfig.REDIS_URI, "", []byte(uiConfig.REDIS_SECRET))
	if err != nil {
		log.Fatal(err)
	}
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24, // 24 hours
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   false, // Set to true in production with HTTPS
	})

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URI"),
		Password: "", // no password set
		DB:       0,  // use default DB
		Protocol: 3,  // specify 2 for RESP 2 or 3 for RESP 3
	})
	status := redisClient.Ping(context.Background())
	if status.String() != "ping: PONG" {
		log.Fatal("cannot connect to redis:", status.String())
	}
	rqHandler := NewRqHandler(redisClient)

	router := gin.Default()
	router.Use(sessions.Sessions("binder-com", store))
	router.LoadHTMLGlob("ui/templates/*")
	router.Static("/static", "ui/static")

	router.GET("/binder", rqHandler.handleBinderHomePage)
	router.GET("/binder/content", rqHandler.handleBinderContent)
	router.POST("/login", rqHandler.handlePOSTUserLogin)
	router.GET("/logout", rqHandler.handleUserLogout)
	router.GET("/notebooks", rqHandler.handleGetNotebookTitles)
	router.GET("/notebooks/:id", rqHandler.handleGetNotebookDetails)
	router.PUT("/notebooks/:id", rqHandler.handleSaveNotebook)
	router.GET("/notebooks-modal", rqHandler.handleNotebooksModal)
	router.POST("/notebooks", rqHandler.handleAddNotebook)
	router.GET("/notebooks/search", rqHandler.handleNotebookLite)
	router.GET("/notebooks/lite", rqHandler.handleNotebookLite)

	router.Run(":5051")
}

func readUiConfig() UIConfig {
	return UIConfig{
		REDIS_URI:    os.Getenv("REDIS_URI"),
		REDIS_SECRET: "secret",
	}
}

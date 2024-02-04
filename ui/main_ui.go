package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dpomian/gobind/ui/uihandlers"
	"github.com/dpomian/gobind/utils"
	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
)

var (
	UnauthorizedRs = gin.H{"error": "unauthorized"}
)

func main() {
	uiConfig, _ := utils.LoadUiConfig("path/to/config")
	fmt.Println("uiconfig: ", uiConfig)
	fmt.Println("redisUri: ", uiConfig.RedisUri)

	store, err := redisStore.NewStore(10, "tcp", uiConfig.RedisUri, "", []byte(uiConfig.RedisSecret))
	if err != nil {
		log.Fatal("error creating the redisStore, err:", err)
	}
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 60 * 60 * 24, // 24 hours
		// SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})

	redisClient := redis.NewClient(&redis.Options{
		Addr:     uiConfig.RedisUri,
		Password: "", // no password set
		DB:       0,  // use default DB
		Protocol: 3,  // specify 2 for RESP 2 or 3 for RESP 3
	})
	status := redisClient.Ping(context.Background())
	if status.String() != "ping: PONG" {
		log.Fatal("cannot connect to redis:", status.String())
	}

	router := gin.Default()
	withLoginRouter := router.Group("/").
		Use(sessions.Sessions("binder-com", store)).
		Use(uihandlers.UiMiddleware(redisClient, uiConfig))

	errorHandler := uihandlers.NewUiErrorResponseHandler(redisClient)
	rqHandler := uihandlers.NewRqHandler(*errorHandler, redisClient, uiConfig)

	router.LoadHTMLGlob("ui/templates/*")
	router.Static("/static", "ui/static")

	// login
	withLoginRouter.POST("/login", rqHandler.HandleLoginPost)
	withLoginRouter.GET("/login", rqHandler.HandleLoginGet)
	withLoginRouter.GET("/logout", rqHandler.HandleUserLogout)
	withLoginRouter.GET("/register", rqHandler.HandleUserRegistrationGet)
	withLoginRouter.POST("/register", rqHandler.HandleUserRegistrationPost)

	// notebooks
	withLoginRouter.GET("/binder", rqHandler.HandleBinderHomePage)
	withLoginRouter.GET("/notebooks/:id", rqHandler.HandleGetNotebookDetails)
	withLoginRouter.PUT("/notebooks/:id", rqHandler.HandleSaveNotebook)
	withLoginRouter.GET("/notebooks-modal", rqHandler.HandleNotebooksModal)
	withLoginRouter.POST("/notebooks", rqHandler.HandleAddNotebook)
	withLoginRouter.GET("/notebooks/search", rqHandler.HandleNotebookLite)
	withLoginRouter.GET("/notebooks/lite", rqHandler.HandleNotebookLite)
	withLoginRouter.GET("/notebooks/topics", rqHandler.HandleGetNotebookTopics)

	router.Run(uiConfig.ServerAddress)
}

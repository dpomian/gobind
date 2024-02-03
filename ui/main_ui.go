package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/dpomian/gobind/ui/uihandlers"
	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
)

var (
	UnauthorizedRs = gin.H{"error": "unauthorized"}
)

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

	router := gin.Default()
	withLoginRouter := router.Group("/").
		Use(sessions.Sessions("binder-com", store)).
		Use(uihandlers.UiMiddleware(redisClient))

	errorHandler := uihandlers.NewUiErrorResponseHandler(redisClient)
	rqHandler := uihandlers.NewRqHandler(*errorHandler, redisClient)

	router.LoadHTMLGlob("ui/templates/*")
	router.Static("/static", "ui/static")

	// login
	withLoginRouter.POST("/login", rqHandler.HandleLoginPost)
	withLoginRouter.GET("/login", rqHandler.HandleLoginGet)
	withLoginRouter.GET("/logout", rqHandler.HandleUserLogout)

	// notebooks
	withLoginRouter.GET("/binder", rqHandler.HandleBinderHomePage)
	withLoginRouter.GET("/notebooks/:id", rqHandler.HandleGetNotebookDetails)
	withLoginRouter.PUT("/notebooks/:id", rqHandler.HandleSaveNotebook)
	withLoginRouter.GET("/notebooks-modal", rqHandler.HandleNotebooksModal)
	withLoginRouter.POST("/notebooks", rqHandler.HandleAddNotebook)
	withLoginRouter.GET("/notebooks/search", rqHandler.HandleNotebookLite)
	withLoginRouter.GET("/notebooks/lite", rqHandler.HandleNotebookLite)
	withLoginRouter.GET("/notebooks/topics", rqHandler.HandleGetNotebookTopics)

	router.Run(":5051")
}

func readUiConfig() UIConfig {
	return UIConfig{
		REDIS_URI:    os.Getenv("REDIS_URI"),
		REDIS_SECRET: "secret",
	}
}

package uihandlers

import (
	"context"
	"net/http"

	"github.com/dpomian/gobind/utils"
	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
)

type RqHandler struct {
	ctx          context.Context
	redisClient  *redis.Client
	errorHandler UiErrorResponseHandler
	config       utils.UIConfig
}

func NewRqHandler(errorHandler UiErrorResponseHandler, redisClient *redis.Client, uiConfig utils.UIConfig) *RqHandler {
	return &RqHandler{
		ctx:          context.Background(),
		redisClient:  redisClient,
		errorHandler: errorHandler,
		config:       uiConfig,
	}
}

func (handler *RqHandler) HandleBinderHomePage(c *gin.Context) {
	accessToken := c.GetString("access_token")
	if len(accessToken) == 0 {
		c.HTML(http.StatusOK, "login.html", nil)
	} else {
		c.HTML(http.StatusOK, "app.html", nil)
	}
}

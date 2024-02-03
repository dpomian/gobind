package uihandlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type IErrorResponseHandler interface {
	HandlerHtmlErrorRs(statusCode int, htmlTemplate string, c *gin.Context)
	HandleJsonErrorRs(statusCode int, jsonData gin.H, c *gin.Context)
}

type UiErrorResponseHandler struct {
	ctx         context.Context
	redisClient *redis.Client
}

func NewUiErrorResponseHandler(redisClient *redis.Client) *UiErrorResponseHandler {
	return &UiErrorResponseHandler{
		ctx:         context.Background(),
		redisClient: redisClient,
	}
}

func (uiErrHandler *UiErrorResponseHandler) HandlerHtmlErrorRs(statusCode int, htmlTemplate string, c *gin.Context) {
	c.HTML(statusCode, htmlTemplate, nil)
}

func (uiErrHandler *UiErrorResponseHandler) HandleJsonErrorRs(statusCode int, jsonData gin.H, c *gin.Context) {
	c.JSON(statusCode, jsonData)
}

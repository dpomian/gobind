package uihandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dpomian/gobind/ui/httputils"
	"github.com/dpomian/gobind/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type loginRq struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

func (handler *RqHandler) HandleLoginPost(c *gin.Context) {
	session := sessions.Default(c)

	utils.InvalidateSessionAndCache(handler.redisClient, handler.ctx, session)

	sessionId := uuid.NewString()
	var loginRequestData loginRq
	if err := c.ShouldBind(&loginRequestData); err != nil {
		fmt.Println(loginRequestData)
		fmt.Println("error binding json", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
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

		c.JSON(statusCode, gin.H{
			"error": "unauthorized",
		})
		return
	}

	session.Set("session_id", sessionId)
	session.Save()

	// add sessionId to redis
	handler.redisClient.Set(handler.ctx, sessionId, responseData, time.Duration(24*time.Hour))

	fmt.Println(string(responseData))

	c.JSON(http.StatusOK, nil)
}

func (handler *RqHandler) HandleLoginGet(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (handler *RqHandler) HandleUserLogout(c *gin.Context) {
	utils.InvalidateSessionAndCache(handler.redisClient, handler.ctx, sessions.Default(c))
	c.Redirect(http.StatusSeeOther, "/login")
}

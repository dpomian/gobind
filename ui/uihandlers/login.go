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

var (
	PasswordsDontMatchError = gin.H{"message": "passwords don't match", "success": false}
	EmailAlreadyExistsError = gin.H{"message": "user with this email already exists", "success": false}
	InternalServerError     = gin.H{"message": "internal server error", "success": false}
	UserSuccessfullyCreated = gin.H{"message": "user successfully created, proceed to login", "success": true}
	BadRequest              = gin.H{"message": "something wrong with the input data", "success": false}
	PasswordTooShort        = gin.H{"message": "password is too short. Min length: 6 characters", "success": false}
)

type loginRqData struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

func (handler *RqHandler) HandleLoginPost(c *gin.Context) {
	session := sessions.Default(c)

	utils.InvalidateSessionAndCache(handler.redisClient, handler.ctx, session)

	sessionId := uuid.NewString()
	var loginRequestData loginRqData
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

func (handler *RqHandler) HandleUserRegistrationGet(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

type registrationRqData struct {
	Username        string `json:"username" form:"username" binding:"required,alphanum"`
	Password        string `json:"password" form:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required,min=6"`
	Email           string `json:"email" form:"email" binding:"required,email"`
}

func (handler *RqHandler) HandleUserRegistrationPost(c *gin.Context) {
	templateFile := "register.html"
	var registrationRequestData registrationRqData
	if err := c.ShouldBind(&registrationRequestData); err != nil {
		fmt.Println(registrationRequestData)
		fmt.Println("error binding json", err)
		c.HTML(http.StatusBadRequest, templateFile, BadRequest)
		return
	}

	if registrationRequestData.Password != registrationRequestData.ConfirmPassword {
		fmt.Println("passwords do not match")
		c.HTML(http.StatusBadRequest, templateFile, PasswordsDontMatchError)
		return
	}

	postData, err := json.Marshal(registrationRequestData)
	if err != nil {
		fmt.Println("error marshalling data:", err)
		c.HTML(http.StatusInternalServerError, templateFile, InternalServerError)
		return
	}

	fmt.Println("postData:", postData)
	url := "http://localhost:5050/api/v1/users"
	headers := httputils.NewHeaders().WithJsonContentTypeHeader()
	_, statusCode, err := httputils.SendPOSTRequest(url, headers, postData)

	switch statusCode {
	case http.StatusOK:
		c.HTML(statusCode, templateFile, UserSuccessfullyCreated)
		return
	case http.StatusForbidden:
		c.HTML(statusCode, templateFile, EmailAlreadyExistsError)
		return
	case http.StatusBadRequest:
		c.HTML(statusCode, templateFile, PasswordTooShort)
	default:
		fmt.Println("err:", err)
		c.HTML(statusCode, templateFile, InternalServerError)
	}
}

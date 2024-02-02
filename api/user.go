package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/dpomian/gobind/db/sqlc"
	"github.com/dpomian/gobind/token"
	"github.com/dpomian/gobind/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserHandler struct {
	config     utils.Config
	tokenMaker token.TokenMaker
	storage    db.Storage
	ctx        context.Context
}

var (
	UserWithEmailAlreadyExists = gin.H{"error": "there is a username already registered with this email"}
	UserWithEmailDoesNotExist  = gin.H{"error": "no user is registered with that email"}
	UserNotAuthorized          = gin.H{"error": "user not authorized"}
)

func NewUserHander(config utils.Config, tokenMaker token.TokenMaker, storage db.Storage, ctx context.Context) *UserHandler {
	return &UserHandler{
		config:     config,
		tokenMaker: tokenMaker,
		storage:    storage,
		ctx:        ctx,
	}
}

type addNewUserRq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

func (handler *UserHandler) AddNewUserHandler(c *gin.Context) {
	var newUserRq addNewUserRq
	if err := c.ShouldBindJSON(&newUserRq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	password, err := utils.HashPassword(newUserRq.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	arg := db.CreateUserParams{
		ID:       uuid.New(),
		Username: newUserRq.Username,
		Email:    newUserRq.Email,
		Password: password,
	}

	_, err = handler.storage.CreateUser(handler.ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			fmt.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "unique_violation":
				c.JSON(http.StatusForbidden, UserWithEmailAlreadyExists)
				return
			}
		} else {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, InternalError)
			return
		}
	}

	c.JSON(http.StatusOK, struct{}{})
}

type loginUserRq struct {
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type loginUserRs struct {
	SessionId             uuid.UUID `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	UserId                string    `json:"user_id"`
}

func (handler *UserHandler) LoginUser(c *gin.Context) {
	var rq loginUserRq
	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user, err := handler.storage.GetUserByEmail(handler.ctx, rq.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, UserWithEmailDoesNotExist)
			return
		}
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	err = utils.CheckPassword(user.Password, rq.Password)
	if err != nil {
		fmt.Println("password check failed. rqpwd:", rq.Password, err)
		c.JSON(http.StatusUnauthorized, UserNotAuthorized)
		return
	}

	accessToken, accessPayload, err := handler.tokenMaker.CreateToken(user.ID.String(), handler.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	refreshToken, refreshPayload, err := handler.tokenMaker.CreateToken(user.ID.String(), handler.config.RefreshTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	session, err := handler.storage.CreateSession(handler.ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	rs := loginUserRs{
		SessionId:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		UserId:                user.ID.String(),
	}
	c.JSON(http.StatusOK, rs)
}

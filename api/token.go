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
)

type TokenHandler struct {
	config     utils.Config
	tokenMaker token.TokenMaker
	storage    db.Storage
	ctx        context.Context
}

func NewTokenHandler(
	config utils.Config,
	tokenMaker token.TokenMaker,
	storage db.Storage,
	ctx context.Context) *TokenHandler {
	return &TokenHandler{
		config:     config,
		tokenMaker: tokenMaker,
		storage:    storage,
		ctx:        ctx,
	}
}

type renewAccessTokenRq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenRs struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (handler *TokenHandler) RenewAccessToken(c *gin.Context) {
	var rq renewAccessTokenRq
	if err := c.ShouldBindJSON(&rq); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	refreshPayload, err := handler.tokenMaker.VerifyToken(rq.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, UserNotAuthorized)
		return
	}

	session, err := handler.storage.GetSession(handler.ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, UserWithEmailDoesNotExist)
			return
		}
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	if session.UserID != refreshPayload.UserId {
		err := fmt.Errorf("incorrect session user")
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	if session.RefreshToken != rq.RefreshToken {
		err := fmt.Errorf("invalid refresh token")
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	accessToken, accessPayload, err := handler.tokenMaker.CreateToken(
		refreshPayload.UserId.String(),
		handler.config.AccessTokenDuration,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, InternalError)
		return
	}

	rs := renewAccessTokenRs{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}
	c.JSON(http.StatusOK, rs)
}

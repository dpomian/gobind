package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dpomian/gobind/token"
	"github.com/gin-gonic/gin"
)

const (
	authHeaderKey  = "authorization"
	authTypeBearer = "bearer"
	authPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.TokenMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authHeaderKey)
		if len(authHeader) == 0 {
			err := errors.New("authorization header is not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != authTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authType)
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		fmt.Println("err:", err)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}
		c.Set(authPayloadKey, payload)
		c.Next()
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("origin")
		fmt.Println("origin:", origin)

		if origin == "http://localhost" || origin == "http://localhost:5051" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		} else {
			c.Next()
		}
	}
}

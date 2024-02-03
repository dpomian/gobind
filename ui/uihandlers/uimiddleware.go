package uihandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dpomian/gobind/ui/httputils"
	"github.com/dpomian/gobind/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
)

const (
	CtxAccessTokenKey = "access_token"
)

func UiMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := getAccessToken(redisClient, context.Background(), c)

		if c.Request.URL.Path == "/login" && accessToken == nil {
			fmt.Println("url path is /login and accessToken is nil")
			c.Next()
			return
		}

		if accessToken == nil {
			// c.HTML(http.StatusOK, "login.html", nil)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		c.Set(CtxAccessTokenKey, *accessToken)
		c.Next()
	}
}

type rsRenewAccesstoken struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

type rqRenewAccessToken struct {
	RefreshToken string `json:"refresh_token"`
}

type loginDetails struct {
	SessionId             string    `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func getAccessToken(redisClient *redis.Client, ctx context.Context, c *gin.Context) *string {
	session := sessions.Default(c)
	sessionId := session.Get("session_id")

	if sessionId == nil {
		fmt.Println("session_id not in request")
		return nil
	}

	fmt.Println("sessionId:", sessionId)

	authData, err := redisClient.Get(ctx, sessionId.(string)).Result()
	if err == redis.Nil {
		fmt.Println("session_id: ", sessionId, "not found in redis: ", err)
		utils.InvalidateSession(session)
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
						fmt.Println("unauthorized trying to refresh access token")
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

					redisClient.Set(ctx, sessionId.(string), marshaledTokens, time.Duration(24*time.Hour))
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

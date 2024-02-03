package utils

import (
	"context"

	"github.com/gin-contrib/sessions"
	"github.com/redis/go-redis/v9"
)

const (
	SessionIdKey = "session_id"
)

func InvalidateSession(session sessions.Session) {
	session.Clear()
	session.Save()
}

func InvalidateSessionAndCache(redisClient *redis.Client, ctx context.Context, session sessions.Session) {
	sessionId := session.Get(SessionIdKey)
	if sessionId != nil {
		redisClient.Del(ctx, sessionId.(string))
	}
	InvalidateSession(session)
}

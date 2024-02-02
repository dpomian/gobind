package utils

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func InvalidateSession(session sessions.Session, c *gin.Context) {
	session.Clear()
	session.Save()
}

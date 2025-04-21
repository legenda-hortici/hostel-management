package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, status int, message string) {
	session := sessions.Default(c)
	session.AddFlash(message)
	session.Save()
	c.Redirect(status, c.Request.Referer())
}

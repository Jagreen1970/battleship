package endpoints

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (c *Controller) Logout(context *gin.Context) {
	session := sessions.Default(context)
	session.Options(sessions.Options{
		MaxAge: -1,
	})
	session.Clear()
	err := session.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
	}
	context.JSON(http.StatusOK, gin.H{"playerName": ""})
}

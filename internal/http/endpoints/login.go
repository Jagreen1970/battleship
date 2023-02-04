package endpoints

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (c *Controller) Login(context *gin.Context) {
	type login struct {
		Username string `json:"username"`
	}
	var l login
	err := context.ShouldBindJSON(&l)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := sessions.Default(context)
	session.Options(sessions.Options{
		MaxAge: 15 * 60,
	})
	var playerName string
	v := session.Get(sessionKeyPlayerName)
	if v == nil {
		playerName = l.Username
	} else {
		playerName = v.(string)
		context.JSON(http.StatusUnauthorized, fmt.Errorf("you are already logged in as player %q", playerName))
		return
	}
	session.Set(sessionKeyPlayerName, playerName)
	err = session.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"playerName": playerName})
}

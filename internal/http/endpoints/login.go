package endpoints

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/Jagreen1970/battleship/internal/battleship"
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
		context.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("you are already logged in as player %q", playerName)})
		return
	}
	session.Set(sessionKeyPlayerName, playerName)
	err = session.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	player, err := c.gameAPI.GetPlayer(playerName)
	if player == nil && errors.Is(err, battleship.ErrorNotFound) {
		player, err = c.gameAPI.NewPlayer(playerName)
	}

	context.JSON(http.StatusOK, player)
}

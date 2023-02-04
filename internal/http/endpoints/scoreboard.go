package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *Controller) Scoreboard(context *gin.Context) {
	playerName := playerFromSession(context)

	scoreboard, err := c.gameAPI.ScoreBoard(playerName)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}
	context.JSON(http.StatusOK, scoreboard)
}

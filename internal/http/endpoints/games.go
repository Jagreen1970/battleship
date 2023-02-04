package endpoints

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const DefaultGamesPerPage = 10

func (c *Controller) Games(context *gin.Context) {
	// fetch pagination from query params
	page, count := paginationParams(context, DefaultGamesPerPage)

	games, err := c.gameAPI.Games(page, count)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}
	context.JSON(http.StatusOK, games)
}

func (c *Controller) GetGame(context *gin.Context) {
	gameID := context.Param("id")
	if gameID == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
		return
	}
	g, err := c.gameAPI.GetGame(gameID)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}
	context.JSON(http.StatusOK, g)
}

func (c *Controller) CreateGame(context *gin.Context) {
	player := playerFromSession(context)
	if player == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid player"})
		return
	}

	game, err := c.gameAPI.NewGame(player)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}

	context.JSON(http.StatusCreated, game)
}

type position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (c *Controller) PlacePin(context *gin.Context) {
	gameID := context.Param("id")
	if gameID == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
		return
	}

	var pin position
	err := context.ShouldBindJSON(pin)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}

	playerName := playerFromSession(context)
	if playerName == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid player"})
		return
	}

	game, err := c.gameAPI.GetGame(gameID)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}

	err = game.PlacePin(playerName, pin.X, pin.Y)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}

	err = c.gameAPI.UpdateGame(game)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}

	context.JSON(http.StatusCreated, game)
}

func (c *Controller) RecoverPin(context *gin.Context) {

}

func (c *Controller) Shoot(context *gin.Context) {

}

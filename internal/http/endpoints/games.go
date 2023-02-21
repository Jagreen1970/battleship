package endpoints

import (
	"fmt"
	"github.com/Jagreen1970/battleship/internal/battleship"
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

	game, err := c.gameAPI.GetGame(gameID)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}

	player := playerFromSession(context)
	if player == "" {
		context.JSON(http.StatusOK, mapToViewerPerspective(game))
		return
	}

	context.JSON(http.StatusOK, mapToPlayerPerspective(player, game))
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

	context.JSON(http.StatusCreated, mapToPlayerPerspective(player, game))
}

func (c *Controller) JoinGame(context *gin.Context) {
	playerName := playerFromSession(context)
	if playerName == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid player - you must be logged in"})
		return
	}

	gameID := context.Param("id")
	if gameID == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid game"})
	}

	player, err := c.gameAPI.GetPlayer(playerName)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "could not load player object"})
		return
	}

	game, err := c.gameAPI.GetGame(gameID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "could not load game object"})
		return
	}

	err = game.Join(player)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}

	game, err = c.gameAPI.UpdateGame(game)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}

	context.JSON(http.StatusAccepted, mapToPlayerPerspective(playerName, game))
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

	pinID := context.Param("pin")
	var pin position
	_, err := fmt.Sscanf(pinID, "%d-%d", pin.X, pin.Y)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid pin id"})
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

	err = game.ValidSetup()
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}

	game, err = c.gameAPI.UpdateGame(game)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}

	context.JSON(http.StatusCreated, mapToPlayerPerspective(playerName, game))
}

func (c *Controller) RecoverPin(context *gin.Context) {
	gameID := context.Param("id")
	if gameID == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
		return
	}

	pinID := context.Param("pin")
	var pin position
	_, err := fmt.Sscanf(pinID, "%d-%d", pin.X, pin.Y)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid pin id"})
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

	err = game.RecoverPin(playerName, pin.X, pin.Y)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}

	game, err = c.gameAPI.UpdateGame(game)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}

	context.JSON(http.StatusCreated, mapToPlayerPerspective(playerName, game))
}

func (c *Controller) StartGame(context *gin.Context) {
	gameID := context.Param("id")
	if gameID == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
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

	err = game.Start(playerName)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}

	game, err = c.gameAPI.UpdateGame(game)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}

	context.JSON(http.StatusOK, mapToPlayerPerspective(playerName, game))
}

func (c *Controller) Target(context *gin.Context) {
	gameID := context.Param("id")
	if gameID == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
		return
	}

	playerName := playerFromSession(context)
	if playerName == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid player"})
		return
	}

	var move battleship.Move
	err := context.ShouldBindJSON(&move)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	g, err := c.gameAPI.GetGame(gameID)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}

	move.Player = playerName
	err = g.MakeMove(move)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}

	g, err = c.gameAPI.UpdateGame(g)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}

	context.JSON(http.StatusOK, mapToPlayerPerspective(playerName, g))
}

type playerGame struct {
	ID      string            `json:"_id,omitempty"`
	User    string            `json:"user"`
	Board   *battleship.Board `json:"board"`
	History []battleship.Move `json:"history"`
	Status  battleship.Status `json:"status"`

	Player1      *battleship.Player `json:"player_1"`
	Player2      *battleship.Player `json:"player_2"`
	PlayerToMove string             `json:"player_to_move"`
}

func mapToPlayerPerspective(name string, g *battleship.Game) playerGame {
	return playerGame{
		ID:      g.ID,
		User:    name,
		Board:   g.Boards[name],
		History: g.History,
		Status:  g.Status,

		Player1:      g.Player1,
		Player2:      g.Player2,
		PlayerToMove: g.PlayerToMove,
	}
}

type viewerGame struct {
	ID      string            `json:"_id,omitempty"`
	User    string            `json:"user"`
	Board   *battleship.Board `json:"board"`
	History []battleship.Move `json:"history"`
	Status  battleship.Status `json:"status"`

	Player1      *battleship.Player `json:"player_1"`
	Player2      *battleship.Player `json:"player_2"`
	PlayerToMove string             `json:"player_to_move"`
}

func mapToViewerPerspective(game *battleship.Game) viewerGame {
	return viewerGame{
		ID:           game.ID,
		User:         "guest",
		Board:        makeViewerBoard(game),
		History:      game.History,
		Status:       game.Status,
		Player1:      game.Player1,
		Player2:      game.Player2,
		PlayerToMove: game.PlayerToMove,
	}
}

func makeViewerBoard(game *battleship.Game) *battleship.Board {
	if game == nil {
		return nil
	}

	var board *battleship.Board
	board.ShipsMap = game.Boards[game.Player1.Name].ShipsMap
	board.ShotsMap = game.Boards[game.Player2.Name].ShipsMap

	return board
}

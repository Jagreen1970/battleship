package endpoints

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Jagreen1970/battleship/internal/game"
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

	user := playerFromSession(context)
	if user == "" {
		user = "guest"
	}

	response := struct {
		Games []*game.Game `json:"games"`
		User  string       `json:"user"`
	}{
		Games: games,
		User:  user,
	}

	context.JSON(http.StatusOK, response)
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
		context.JSON(http.StatusOK, viewerPerspective(game))
		return
	}

	context.JSON(http.StatusOK, playerPerspective(player, game))
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

	context.JSON(http.StatusCreated, playerPerspective(player, game))
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

	context.JSON(http.StatusAccepted, playerPerspective(playerName, game))
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

	context.JSON(http.StatusCreated, playerPerspective(playerName, game))
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

	game, err = c.gameAPI.UpdateGame(game)
	if err != nil {
		context.JSON(mapErrorToStatusErr(err))
		return
	}

	context.JSON(http.StatusCreated, playerPerspective(playerName, game))
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

	context.JSON(http.StatusOK, playerPerspective(playerName, game))
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

	var move game.Move
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

	context.JSON(http.StatusOK, playerPerspective(playerName, g))
}

type gameView struct {
	ID      string      `json:"_id,omitempty"`
	User    string      `json:"user"`
	Board   *game.Board `json:"board"`
	History []game.Move `json:"history"`
	Status  game.Status `json:"status"`

	Player1      *game.Player `json:"player_1"`
	Player2      *game.Player `json:"player_2"`
	PlayerToMove string       `json:"player_to_move"`
}

func playerPerspective(name string, g *game.Game) gameView {
	return gameView{
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

func viewerPerspective(game *game.Game) gameView {
	return gameView{
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

func makeViewerBoard(g *game.Game) *game.Board {
	if g == nil {
		return nil
	}

	board := &game.Board{
		PinsAvailable: 0,
		Maps: [2]*game.BoardMap{
			{
				Title: "shots",
				Map:   [10]game.FieldRow{},
			},
			{
				Title: "ships",
				Map:   [10]game.FieldRow{},
			},
		},
		Fleet: nil,
	}
	return board
}

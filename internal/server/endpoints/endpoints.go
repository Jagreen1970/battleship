package endpoints

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/Jagreen1970/battleship/internal/game"
)

type GameAPI interface {
	ScoreBoard(playerName string) (*game.ScoreBoard, error)
	Games(page int, count int) ([]*game.Game, error)
	GetGame(id string) (*game.Game, error)
	NewGame(player string) (*game.Game, error)
	UpdateGame(g *game.Game) (*game.Game, error)
	GetPlayer(playerName string) (*game.Player, error)
	NewPlayer(playerName string) (*game.Player, error)
}

func playerFromSession(context *gin.Context) string {
	session := sessions.Default(context)
	v := session.Get(sessionKeyPlayerName)
	if v == nil {
		return ""
	}
	playerName := v.(string)

	return playerName
}

func paginationParams(context *gin.Context, defaultItems int) (int, int) {
	// fetch page and items from context
	page := context.Query("page")
	items := context.Query("items")
	// return page and items
	return pageAsInt(page), itemsAsInt(items, defaultItems)
}

func itemsAsInt(items string, defaultItems int) int {
	i, err := strconv.Atoi(items)
	if err != nil {
		return defaultItems
	}
	return i
}

func pageAsInt(page string) int {
	i, err := strconv.Atoi(page)
	if err != nil {
		return 1
	}
	return i
}

// mapErrorToStatusError maps an error to a http-status and an error message
//
//	ErrorNotFound  = errors.New("not found")
//	ErrorIllegal   = errors.New("illegal action")
//	ErrorNotReady  = errors.New("not ready")
//	ErrorInvalid   = errors.New("invalid")
//	ErrorAmbiguous = errors.New("duplicate")
func mapErrorToStatusErr(err error) (int, any) {
	if errors.Is(err, game.ErrorNotFound) {
		return http.StatusNotFound, gin.H{"error": err.Error()}
	}
	if errors.Is(err, game.ErrorIllegal) {
		return http.StatusForbidden, gin.H{"error": err.Error()}
	}
	if errors.Is(err, game.ErrorInvalid) {
		return http.StatusBadRequest, gin.H{"error": err.Error()}
	}
	if errors.Is(err, game.ErrorNotReady) {
		return http.StatusConflict, gin.H{"error": err.Error()}
	}
	return http.StatusInternalServerError, gin.H{"error": err.Error()}
}

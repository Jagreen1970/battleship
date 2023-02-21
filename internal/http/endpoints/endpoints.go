package endpoints

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/Jagreen1970/battleship/internal/battleship"
)

type GameAPI interface {
	ScoreBoard(playerName string) (*battleship.ScoreBoard, error)
	Games(page int, count int) ([]*battleship.Game, error)
	GetGame(id string) (*battleship.Game, error)
	NewGame(player string) (*battleship.Game, error)
	UpdateGame(g *battleship.Game) (*battleship.Game, error)
	GetPlayer(playerName string) (*battleship.Player, error)
	NewPlayer(playerName string) (*battleship.Player, error)
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
	p, err := strconv.Atoi(page)
	if err != nil {
		return 1
	}
	return p
}

// mapErrorToStatusError maps an error to a http-status and an error message
//	ErrorNotFound  = errors.New("not found")
//	ErrorIllegal   = errors.New("illegal action")
//	ErrorNotReady  = errors.New("not ready")
//	ErrorInvalid   = errors.New("invalid")
//	ErrorAmbiguous = errors.New("duplicate")
func mapErrorToStatusErr(err error) (int, any) {
	switch {
	case errors.Is(err, battleship.ErrorNotFound):
		return http.StatusNotFound, gin.H{"error": err.Error()}
	case errors.Is(err, battleship.ErrorIllegal):
		return http.StatusForbidden, gin.H{"error": err.Error()}
	case errors.Is(err, battleship.ErrorNotReady):
		return http.StatusBadRequest, gin.H{"error": err.Error()}
	case errors.Is(err, battleship.ErrorInvalid):
		return http.StatusBadRequest, gin.H{"error": err.Error()}
	case errors.Is(err, battleship.ErrorAmbiguous):
		return http.StatusConflict, gin.H{"error": err.Error()}
	}
	return http.StatusInternalServerError, gin.H{"error": err.Error()}
}

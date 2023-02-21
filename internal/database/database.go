package database

import (
	"github.com/Jagreen1970/battleship/internal/app"
	"github.com/Jagreen1970/battleship/internal/battleship"
	"github.com/Jagreen1970/battleship/internal/database/mongodb"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	UnknownDatabaseError Error = "UnknownDatabaseError"
)

type Database interface {
	Connect() error
	Disconnect() error
	Ping() error
	Close() error

	CreatePlayer(playerName string) (*battleship.Player, error)
	FindPlayerByName(username string) (*battleship.Player, error)

	QueryGames(page int, count int) ([]*battleship.Game, error)
	CreateGame(game *battleship.Game) (*battleship.Game, error)
	FindGameByID(id string) (*battleship.Game, error)
	UpdateGame(game *battleship.Game) (*battleship.Game, error)
}

func New() (Database, error) {
	switch app.DatabaseDriver() {
	case "mongo":
		return mongodb.NewMongoDB()
	}

	return nil, UnknownDatabaseError
}

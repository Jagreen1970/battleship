package database

import (
	"github.com/Jagreen1970/battleship/internal/app"
	"github.com/Jagreen1970/battleship/internal/game"
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

	FindPlayerByName(username string) (game.Player, error)
}

func New() (Database, error) {
	switch app.DatabaseDriver() {
	case "mongo":
		return newMongoDB()
	}

	return nil, UnknownDatabaseError
}

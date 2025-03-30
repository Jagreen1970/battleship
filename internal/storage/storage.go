package storage

import (
	"fmt"

	"github.com/Jagreen1970/battleship/internal/app"
	"github.com/Jagreen1970/battleship/internal/game"
	"github.com/Jagreen1970/battleship/internal/storage/mongodb"
)

// Storage defines the interface for storage operations
type Storage interface {
	Connect() error
	Disconnect() error
	Ping() error
	Close() error

	CreatePlayer(playerName string) (*game.Player, error)
	FindPlayerByName(username string) (*game.Player, error)

	QueryGames(page int, count int) ([]*game.Game, error)
	CreateGame(game *game.Game) (*game.Game, error)
	FindGameByID(id string) (*game.Game, error)
	UpdateGame(game *game.Game) (*game.Game, error)
}

// New creates a new storage instance based on the configuration
func New(cfg app.DatabaseConfig) (Storage, error) {
	switch cfg.Driver {
	case "mongo":
		return mongodb.NewMongoDB(cfg)
	default:
		return nil, fmt.Errorf("unsupported storage driver: %s", cfg.Driver)
	}
}

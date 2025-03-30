package mongodb

import (
	"context"
	"fmt"

	"github.com/Jagreen1970/battleship/internal/game"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Player struct {
	ID     primitive.ObjectID `bson:"_id"`
	Player *game.Player       `bson:"player"`
}

// CreatePlayer creates a new player in the database
func (m *MongoDB) CreatePlayer(playerName string) (*game.Player, error) {
	player := Player{
		Player: &game.Player{
			Name: playerName,
		},
	}

	collection := m.client.Database(m.cfg.Name).Collection("players")

	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeout)
	defer cancel()

	result, err := collection.InsertOne(ctx, &player)
	if err != nil {
		return nil, fmt.Errorf("error creating player: %w", err)
	}

	player.Player.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return player.Player, nil
}

// FindPlayerByName retrieves a player by their username
func (m *MongoDB) FindPlayerByName(username string) (*game.Player, error) {
	collection := m.client.Database(m.cfg.Name).Collection("players")

	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeout)
	defer cancel()

	result := collection.FindOne(ctx, bson.D{primitive.E{Key: "player.name", Value: username}})

	var player Player
	if err := result.Decode(&player); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, game.ErrorNotFound
		}
		return nil, fmt.Errorf("error fetching player: %w", err)
	}

	player.Player.ID = player.ID.Hex()
	return player.Player, nil
}

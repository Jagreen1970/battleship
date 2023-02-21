package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jagreen1970/battleship/internal/app"
	"github.com/Jagreen1970/battleship/internal/battleship"
)

type Player struct {
	ID     primitive.ObjectID `bson:"_id"`
	Player *battleship.Player `bson:"player"`
}

func (m *MongoDB) CreatePlayer(playerName string) (*battleship.Player, error) {
	player := Player{
		Player: &battleship.Player{
			Name: playerName,
		},
	}

	collection := m.db.Database(app.DatabaseName()).Collection("players")

	ctx, cancelContext := context.WithTimeout(context.Background(), app.DatabaseTimeout())
	defer cancelContext()

	result, err := collection.InsertOne(ctx, &player)
	if err != nil {
		return nil, err
	}

	player.Player.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return player.Player, nil
}

func (m *MongoDB) FindPlayerByName(username string) (*battleship.Player, error) {
	collection := m.db.Database(app.DatabaseName()).Collection("players")

	ctx, cancelContext := context.WithTimeout(context.Background(), app.DatabaseTimeout())
	defer cancelContext()

	result := collection.FindOne(ctx, bson.D{primitive.E{Key: "player.name", Value: username}})

	var player Player
	if err := result.Decode(&player); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, battleship.ErrorNotFound
		}
		return nil, fmt.Errorf("error fetching player: %w", err)
	}

	player.Player.ID = player.ID.Hex()
	return player.Player, nil
}

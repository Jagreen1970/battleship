package mongodb

import (
	"context"
	"fmt"

	"github.com/Jagreen1970/battleship/internal/game"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Game struct {
	ID   primitive.ObjectID `bson:"_id"`
	Game *game.Game         `bson:"game"`
}

// QueryGames retrieves a list of games with pagination
func (m *MongoDB) QueryGames(page int, count int) ([]*game.Game, error) {
	collection := m.client.Database(m.cfg.Name).Collection("games")

	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeout)
	defer cancel()

	skip := int64((page - 1) * count)
	opts := options.Find().SetSkip(skip).SetLimit(int64(count))
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, fmt.Errorf("error querying games: %w", err)
	}
	defer cursor.Close(ctx)

	var games []Game
	if err := cursor.All(ctx, &games); err != nil {
		return nil, fmt.Errorf("error decoding games: %w", err)
	}

	ret := make([]*game.Game, len(games))
	for i, g := range games {
		ret[i] = g.Game
	}
	return ret, nil
}

// CreateGame creates a new game in the database
func (m *MongoDB) CreateGame(game *game.Game) (*game.Game, error) {
	g := Game{
		Game: game,
	}

	collection := m.client.Database(m.cfg.Name).Collection("games")

	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeout)
	defer cancel()

	result, err := collection.InsertOne(ctx, &g)
	if err != nil {
		return nil, fmt.Errorf("error creating game: %w", err)
	}

	g.Game.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return g.Game, nil
}

// FindGameByID retrieves a game by its ID
func (m *MongoDB) FindGameByID(id string) (*game.Game, error) {
	gameID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid game ID: %w", err)
	}

	collection := m.client.Database(m.cfg.Name).Collection("games")

	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeout)
	defer cancel()

	var g Game
	err = collection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: gameID}}).Decode(&g)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, game.ErrorNotFound
		}
		return nil, fmt.Errorf("error finding game: %w", err)
	}

	return g.Game, nil
}

// UpdateGame updates an existing game in the database
func (m *MongoDB) UpdateGame(g *game.Game) (*game.Game, error) {
	gameID, err := primitive.ObjectIDFromHex(g.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid game ID: %w", err)
	}

	collection := m.client.Database(m.cfg.Name).Collection("games")

	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeout)
	defer cancel()

	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "game", Value: g},
		}},
	}

	result, err := collection.UpdateOne(ctx, bson.D{primitive.E{Key: "_id", Value: gameID}}, update)
	if err != nil {
		return nil, fmt.Errorf("error updating game: %w", err)
	}

	if result.ModifiedCount == 0 {
		return nil, game.ErrorNotFound
	}

	return g, nil
}

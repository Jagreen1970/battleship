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

	// Fix: Make sure skip is never negative - page is zero-indexed
	skip := int64(page * count)
	if skip < 0 {
		skip = 0
	}
	
	opts := options.Find().SetSkip(skip).SetLimit(int64(count))
	// Sort by _id in descending order to get newest games first
	opts.SetSort(bson.D{primitive.E{Key: "_id", Value: -1}})
	
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
		// Ensure the game ID is set correctly
		ret[i] = g.Game
		ret[i].ID = g.ID.Hex()
	}
	return ret, nil
}

// CreateGame creates a new game in the database
func (m *MongoDB) CreateGame(game *game.Game) (*game.Game, error) {
	// Generate a new ObjectID
	objID := primitive.NewObjectID()
	
	g := Game{
		ID:   objID,  // Set the generated ID
		Game: game,
	}

	collection := m.client.Database(m.cfg.Name).Collection("games")

	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeout)
	defer cancel()

	_, err := collection.InsertOne(ctx, &g)
	if err != nil {
		return nil, fmt.Errorf("error creating game: %w", err)
	}

	// Set the ID from the generated ObjectID
	g.Game.ID = objID.Hex()
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

	// Make sure to set the ID in the game object
	g.Game.ID = id
	return g.Game, nil
}

// FindGameByName retrieves a game by its name
func (m *MongoDB) FindGameByName(name string) (*game.Game, error) {
	if name == "" {
		return nil, fmt.Errorf("game name cannot be empty: %w", game.ErrorInvalidInput)
	}

	collection := m.client.Database(m.cfg.Name).Collection("games")

	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeout)
	defer cancel()

	var g Game
	err := collection.FindOne(ctx, bson.D{primitive.E{Key: "game.name", Value: name}}).Decode(&g)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, game.ErrorNotFound
		}
		return nil, fmt.Errorf("error finding game by name: %w", err)
	}

	// Make sure to set the ID in the game object
	g.Game.ID = g.ID.Hex()
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

// DeleteGame deletes a specific game by ID
func (m *MongoDB) DeleteGame(id string) error {
	gameID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid game ID: %w", err)
	}

	collection := m.client.Database(m.cfg.Name).Collection("games")

	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeout)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.D{primitive.E{Key: "_id", Value: gameID}})
	if err != nil {
		return fmt.Errorf("error deleting game: %w", err)
	}

	if result.DeletedCount == 0 {
		return game.ErrorNotFound
	}

	return nil
}

// DeleteAllGames deletes all games from the database
func (m *MongoDB) DeleteAllGames() (int, error) {
	collection := m.client.Database(m.cfg.Name).Collection("games")

	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeout)
	defer cancel()

	result, err := collection.DeleteMany(ctx, bson.D{})
	if err != nil {
		return 0, fmt.Errorf("error deleting all games: %w", err)
	}

	return int(result.DeletedCount), nil
}

package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Jagreen1970/battleship/internal/app"
	"github.com/Jagreen1970/battleship/internal/battleship"
)

type Game struct {
	ID   primitive.ObjectID `bson:"_id"`
	Game *battleship.Game   `bson:"battleship"`
}

func (m *MongoDB) QueryGames(page int, count int) ([]*battleship.Game, error) {
	collection := m.db.Database(app.DatabaseName()).Collection("games")

	ctx, cancelContext := context.WithTimeout(context.Background(), app.DatabaseTimeout())
	defer cancelContext()

	skip := int64((page - 1) * count)
	limit := int64(count)
	cursor, err := collection.Find(ctx, bson.M{}, &options.FindOptions{Limit: &limit, Skip: &skip})
	if err != nil {
		return nil, err
	}

	var games []Game
	err = cursor.All(ctx, &games)
	if err != nil {
		return nil, err
	}

	ret := make([]*battleship.Game, len(games))
	for i, g := range games {
		nextGame := g.Game
		nextGame.ID = g.ID.Hex()
		ret[i] = nextGame
	}
	return ret, nil
}

func (m *MongoDB) CreateGame(game *battleship.Game) (*battleship.Game, error) {
	g := Game{
		Game: game,
	}

	collection := m.db.Database(app.DatabaseName()).Collection("games")

	ctx, cancelContext := context.WithTimeout(context.Background(), app.DatabaseTimeout())
	defer cancelContext()

	result, err := collection.InsertOne(ctx, &g)
	if err != nil {
		return nil, err
	}

	game.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return game, nil
}

func (m *MongoDB) FindGameByID(id string) (*battleship.Game, error) {
	gameID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	collection := m.db.Database(app.DatabaseName()).Collection("games")

	ctx, cancelContext := context.WithTimeout(context.Background(), app.DatabaseTimeout())
	defer cancelContext()

	result := collection.FindOne(ctx, bson.M{"_id": gameID})

	var g Game
	if err := result.Decode(&g); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, battleship.ErrorNotFound
		}
		return nil, fmt.Errorf("error decoding game: %w", err)
	}

	g.Game.ID = gameID.Hex()
	return g.Game, nil
}

func (m *MongoDB) UpdateGame(g *battleship.Game) (*battleship.Game, error) {
	gameId, err := primitive.ObjectIDFromHex(g.ID)
	if err != nil {
		return nil, err
	}

	collection := m.db.Database(app.DatabaseName()).Collection("games")
	ctx, cancelContext := context.WithTimeout(context.Background(), app.DatabaseTimeout())
	defer cancelContext()

	result, err := collection.UpdateByID(ctx, gameId, bson.M{"$set": bson.M{"battleship": g}})
	if err != nil {
		return nil, fmt.Errorf("could not update game: %w", err)
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("could not update game with id %s: %w", g.ID, battleship.ErrorNotFound)
	}

	return g, nil
}

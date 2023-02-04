package database

import (
	"context"
	"github.com/Jagreen1970/battleship/internal/app"
	"github.com/Jagreen1970/battleship/internal/game"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	db *mongo.Client
}

func (m *MongoDB) FindPlayerByName(username string) (game.Player, error) {
	panic("implement me")
}

func newMongoDB() (*MongoDB, error) {
	clientOptions := mongoClientOptions()
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	return &MongoDB{
		db: client,
	}, nil
}

func mongoClientOptions() *options.ClientOptions {
	mongoOptions := options.Client().SetHosts(
		[]string{app.DatabaseURL()},
	).SetAuth(options.Credential{
		Username: app.DatabaseUser(),
		Password: app.DatabasePassword(),
	})
	return mongoOptions
}

func (m *MongoDB) Connect() error {
	ctx, cancelContext := context.WithTimeout(context.Background(), app.DatabaseTimeout())
	defer cancelContext()

	return m.db.Connect(ctx)
}

func (m *MongoDB) Disconnect() error {
	ctx, cancelContext := context.WithTimeout(context.Background(), app.DatabaseTimeout())
	defer cancelContext()

	return m.db.Disconnect(ctx)
}

func (m *MongoDB) Ping() error {
	ctx, cancelContext := context.WithTimeout(context.Background(), app.DatabaseTimeout())
	defer cancelContext()

	return m.db.Ping(ctx, readpref.Primary())
}

func (m *MongoDB) Close() error {
	return m.Disconnect()
}

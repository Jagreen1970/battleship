package mongodb

import (
	"context"
	"fmt"

	"github.com/Jagreen1970/battleship/internal/app"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
	cfg    app.DatabaseConfig
}

func NewMongoDB(cfg app.DatabaseConfig) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(cfg.URL)
	if cfg.User != "" && cfg.Password != "" {
		clientOptions.SetAuth(options.Credential{
			Username: cfg.User,
			Password: cfg.Password,
		})
	}

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %v", err)
	}

	return &MongoDB{
		client: client,
		cfg:    cfg,
	}, nil
}

func (m *MongoDB) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeout)
	defer cancel()

	if err := m.client.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	return nil
}

func (m *MongoDB) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeout)
	defer cancel()

	if err := m.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %v", err)
	}

	return nil
}

func (m *MongoDB) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeout)
	defer cancel()

	if err := m.client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	return nil
}

func (m *MongoDB) Close() error {
	return m.Disconnect()
}

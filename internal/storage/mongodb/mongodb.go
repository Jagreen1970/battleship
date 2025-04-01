package mongodb

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Jagreen1970/battleship/internal/app"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
	dbName string
	cfg    app.DatabaseConfig
}

func NewMongoDB(cfg app.DatabaseConfig) (*MongoDB, error) {
	// Create the connection URI
	uri := cfg.URL
	if cfg.User != "" && cfg.Password != "" {
		// Parse the existing URL to extract host and path components
		parsedURL, err := url.Parse(cfg.URL)
		if err != nil {
			return nil, fmt.Errorf("invalid MongoDB URL: %w", err)
		}
		
		host := parsedURL.Host
		
		// Create the connection URI with authentication credentials
		uri = fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=admin",
			url.QueryEscape(cfg.User),
			url.QueryEscape(cfg.Password),
			host,
			cfg.Name,
		)
	}

	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetTimeout(cfg.Timeout)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	return &MongoDB{
		client: client,
		dbName: cfg.Name,
		cfg:    cfg,
	}, nil
}

func (m *MongoDB) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.cfg.Timeout)
	defer cancel()

	if err := m.client.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
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

// NOTE: Other implementations like FindPlayerByName are in player.go and game.go
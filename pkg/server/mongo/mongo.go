package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/gookit/slog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/4chain-ag/go-overlay-services/pkg/config"
)

// Client represents the MongoDB client and database connection
type Client struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// Connect establishes a connection to the MongoDB server using the provided configuration.
func Connect(cfg *config.Config) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.Mongo.URI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(cfg.Mongo.Database)
	slog.Infof("MongoDB connected to %s, using DB: %s", cfg.Mongo.URI, cfg.Mongo.Database)

	return &Client{
		Client:   client,
		Database: db,
	}, nil
}

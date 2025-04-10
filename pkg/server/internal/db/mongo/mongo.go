package mongo

import (
	"context"
	"fmt"
	"strings"

	"github.com/gookit/slog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	AuthDB   string `mapstructure:"auth_db"`
}

func (c *Config) HasCredentials() bool {
	return strings.TrimSpace(c.Username) != "" && strings.TrimSpace(c.Password) != ""
}

var DefaultConfig = Config{
	URI:      "mongodb://localhost:27017",
	Database: "overlay",
	Username: "",
	Password: "",
	AuthDB:   "admin",
}

func New(ctx context.Context, cfg *Config) (*mongo.Client, *mongo.Database, error) {
	opts := options.Client().ApplyURI(cfg.URI)
	if cfg.HasCredentials() {
		opts.SetAuth(options.Credential{
			Username:   cfg.Username,
			Password:   cfg.Password,
			AuthSource: cfg.AuthDB,
		})
	}

	cli, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	if err := cli.Ping(ctx, nil); err != nil {
		return nil, nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := cli.Database(cfg.Database)
	slog.Infof("MongoDB connected to %s, using DB: %s", cfg.URI, cfg.Database)
	return cli, db, nil
}

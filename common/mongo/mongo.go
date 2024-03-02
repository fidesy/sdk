package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

func Connect(ctx context.Context, conf Config) (*mongo.Client, error) {
	if conf.Port != "" {
		conf.Host += ":" + conf.Port
	}

	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s",
		conf.Username,
		conf.Password,
		conf.Host,
	))

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}

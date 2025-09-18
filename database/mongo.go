package database

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	ClusterURL string
	Host       string
	Port       string
	User       string
	Password   string
	Database   string
}

type MongoConnection struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoConnection(config MongoConfig) (*MongoConnection, error) {
	url := buildMongioURL(config)

	clientOptions := options.Client().
		ApplyURI(url)
	clientOptions.SetConnectTimeout(10 * time.Second)
	clientOptions.SetServerSelectionTimeout(10 * time.Second).
		SetMaxPoolSize(100).
		SetMinPoolSize(5)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	if err = client.Ping(ctx, nil); err != nil {
		return nil, logrus.Errorf("could not connect to MongoDB: %v", err)
	}
}

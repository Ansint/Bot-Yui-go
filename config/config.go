package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	Discord     DiscordConfig  `json:"discord"`
	DB          DatabaseConfig `json:"db"`
	MongoClient *mongo.Client
}

type DiscordConfig struct {
	Token        string `json:"token" env:"DISCORD_TOKEN"`
	ClientID     string `json:"client_id" env:"DISCORD_CLIENT_ID"`
	LogGuildID   string `json:"log_guild_id" env:"DISCORD_LOG_GUILD_ID"`
	LogChannelID string `json:"log_channel_id" env:"DISCORD_LOG_CHANNEL_ID"`
}

type DatabaseConfig struct {
	ClusterURL string `json:"cluster_url" env:"DB_CLUSTER_URL"`
	Port       string `json:"port" env:"DB_PORT"`
	User       string `json:"user" env:"DB_USER"`
	Password   string `json:"pwd" env:"DB_PWD"`
	Database   string `json:"database" env:"DB_DATABASE"`
	Collection string `json:"collection" env:"DB_COLLECTION"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logrus.Fatal("Error loading .env file")
	}
	config := &Config{
		Discord: DiscordConfig{
			Token:        os.Getenv("DISCORD_TOKEN"),
			ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
			LogGuildID:   os.Getenv("DISCORD_LOG_GUILD_ID"),
			LogChannelID: os.Getenv("DISCORD_LOG_CHANNEL_ID"),
		},
		DB: DatabaseConfig{
			ClusterURL: os.Getenv("DB_CLUSTER_URL"),
			Port:       os.Getenv("DB_PORT"),
			User:       os.Getenv("DB_USER"),
			Password:   os.Getenv("DB_PWD"),
			Database:   os.Getenv("DB_DATABASE"),
		},
	}

	return config, nil
}

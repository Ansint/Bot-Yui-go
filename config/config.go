package config

import (
	"os"

	"github.com/Ansint/yui-lootbot-for-sirus/database"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Discord DiscordConfig  `json:"discord"`
	DB      DatabaseConfig `json:"db"`
}

type DiscordConfig struct {
	Token    string `json:"token" env:"DISCORD_TOKEN"`
	ClientID string `json:"client_id" env:"DISCORD_CLIENT_ID"`
}

type DatabaseConfig struct {
	ClusterURL string `json:"cluster_url" env:"DB_CLUSTER_URL"`
	Host       string `json:"host" env:"DB_HOST"`
	Port       string `json:"port" env:"DB_PORT"`
	User       string `json:"user" env:"DB_USER"`
	Password   string `json:"pwd" env:"DB_PWD"`
	Database   string `json:"database" env:"DB_DATABASE"`
	Collection string `json:"collection" env:"DB_COLLECTION"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Error loading .env file, using system environment variables")
	}

	config := &Config{
		Discord: DiscordConfig{
			Token:    os.Getenv("DISCORD_TOKEN"),
			ClientID: os.Getenv("DISCORD_CLIENT_ID"),
		},
		DB: DatabaseConfig{
			ClusterURL: os.Getenv("DB_CLUSTER_URL"),
			Host:       os.Getenv("DB_HOST"),
			Port:       os.Getenv("DB_PORT"),
			User:       os.Getenv("DB_USER"),
			Password:   os.Getenv("DB_PWD"),
			Database:   os.Getenv("DB_DATABASE"),
			Collection: os.Getenv("DB_COLLECTION"),
		},
	}

	return config, nil
}

func (c *Config) ToMongoConfig() database.MongoConfig {
	return database.MongoConfig{
		ClusterURL: c.DB.ClusterURL,
		Host:       c.DB.Host,
		Port:       c.DB.Port,
		User:       c.DB.User,
		Password:   c.DB.Password,
		Database:   c.DB.Database,
	}
}

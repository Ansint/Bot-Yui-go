package database

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Определяем структуры данных
type Settings struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty"`
	GuildID     string              `bson:"guild_id"`
	ChangelogID *primitive.ObjectID `bson:"changelog_id,omitempty"`
	LootID      *primitive.ObjectID `bson:"loot_id,omitempty"`
}

type Loot struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ChannelID    string             `bson:"channel_id"`
	RealmID      string             `bson:"realm_id"`
	GuildSirusID string             `bson:"guild_sirus_id"`
}

type MongoDB struct {
	Client    *mongo.Client
	db        *mongo.Database
	settings  *mongo.Collection
	changelog *mongo.Collection
	loot      *mongo.Collection
	records   *mongo.Collection
}

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

// NewMongoDB создает экземпляр MongoDB с инициализированными коллекциями
func NewMongoDB(connection *MongoConnection) *MongoDB {
	return &MongoDB{
		Client:    connection.Client,
		db:        connection.Database,
		settings:  connection.Database.Collection("settings"),
		changelog: connection.Database.Collection("changelog"),
		loot:      connection.Database.Collection("loot"),
		records:   connection.Database.Collection("records"),
	}
}

func (m *MongoDB) SetLootChanel(ctx context.Context, guildID string, channelID string, realmID string, guildSirusID string) error {
	var existingSettings Settings
	err := m.settings.FindOne(ctx, bson.M{"guild_id": guildID}).Decode(&existingSettings)

	if err == mongo.ErrNoDocuments {
		// Создаем новую запись loot
		lootResult, err := m.loot.InsertOne(ctx, Loot{
			ChannelID:    channelID,
			RealmID:      realmID,
			GuildSirusID: guildSirusID,
		})
		if err != nil {
			return err
		}

		lootID := lootResult.InsertedID.(primitive.ObjectID)

		// Создаем новые настройки
		_, err = m.settings.InsertOne(ctx, Settings{
			GuildID:     guildID,
			ChangelogID: nil,
			LootID:      &lootID,
		})
		return err

	} else if err != nil {
		return err
	} else {
		// Обновляем существующую запись loot
		if existingSettings.LootID == nil {
			// Если LootID не существует, создаем новую запись loot
			lootResult, err := m.loot.InsertOne(ctx, Loot{
				ChannelID:    channelID,
				RealmID:      realmID,
				GuildSirusID: guildSirusID,
			})
			if err != nil {
				return err
			}

			lootID := lootResult.InsertedID.(primitive.ObjectID)
			_, err = m.settings.UpdateOne(
				ctx,
				bson.M{"guild_id": guildID},
				bson.M{"$set": bson.M{"loot_id": lootID}},
			)
			return err
		}

		// Обновляем существующую запись loot
		_, err := m.loot.UpdateOne(
			ctx,
			bson.M{"_id": existingSettings.LootID},
			bson.M{"$set": bson.M{
				"channel_id":     channelID,
				"realm_id":       realmID,
				"guild_sirus_id": guildSirusID,
			}},
		)
		return err
	}
}

func NewMongoConnection(config MongoConfig) (*MongoConnection, error) {
	url := buildMongoURL(config)

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
		return nil, err
	}
	db := client.Database(config.Database)
	logrus.Info("Connected to MongoDB")

	return &MongoConnection{
		Client:   client,
		Database: db,
	}, nil
}

func (mc *MongoConnection) GetCollection(collection string) *mongo.Collection {
	return mc.Database.Collection(collection)
}

func (mc *MongoConnection) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mc.Client.Disconnect(ctx)
}

func (mc *MongoConnection) IsConnected() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mc.Client.Ping(ctx, nil) == nil
}

func buildMongoURL(config MongoConfig) string {
	if config.ClusterURL != "" {
		return config.ClusterURL
	}
	if config.User != "" && config.Password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%s", config.User, config.Password, config.Host, config.Port)
	}
	return fmt.Sprintf("mongodb://%s:%s", config.Host, config.Port)
}

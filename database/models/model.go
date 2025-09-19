package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Settings struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	GuildID     string              `bson:"guild_id" json:"guild_id"`
	ChangelogID *primitive.ObjectID `bson:"changelog_id,omitempty" json:"changelog_id,omitempty"`
	LootID      *primitive.ObjectID `bson:"loot_id,omitempty" json:"loot_id,omitempty"`
}

type Changelog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ChannelID string             `bson:"channel_id" json:"channel_id"`
}

type Loot struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ChannelID    string             `bson:"channel_id" json:"channel_id"`
	RealmID      string             `bson:"realm_id" json:"realm_id"`
	GuildSirusID string             `bson:"guild_sirus_id" json:"guild_sirus_id"`
}

type Records struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GuildID string             `bson:"guild_id" json:"guild_id"`
	Records []string           `bson:"records" json:"records"`
}

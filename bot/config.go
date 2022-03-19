package bot

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"
)

type Guild struct {
	Active   bool      `firestore:"active"`
	JoinedAt time.Time `firestore:"joined-at"`

	// The person who owns the Discord
	OwnerID string `firestore:"owner-id"`

	// GuildID is the ID of the guild
	GuildID   string `firestore:"guild-id"`
	GuildName string `firestore:"guild-name"`

	// Premint settings
	PremintAPIKey string `firestore:"premint-api-key"`
	PremintRole   string `firestore:"premint-role"`
	PremintRoleID string `firestore:"premint-role-id"`
}

type ConfigParams struct {
	config *Guild
	doc    *firestore.DocumentSnapshot
}

func getConfig(
	ctx context.Context,
	logger *zap.SugaredLogger,
	database *firestore.Client,
	guildID string,
) *ConfigParams {
	// Fetch the config for the guild
	docSnap, err := database.Collection("guilds").Doc(guildID).Get(ctx)
	if err != nil {
		logger.Errorw("Failed to get config", "guild", guildID, "error", err)
		return nil
	}

	// Get the config
	config := &Guild{}
	err = docSnap.DataTo(config)
	if err != nil {
		return nil
	}

	return &ConfigParams{
		config: config,
		doc:    docSnap,
	}
}

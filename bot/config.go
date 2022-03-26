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

	// GuildOwnerID is the person who owns the Discord guild
	GuildOwnerID string `firestore:"owner-id"`
	// GuildID is the ID of the guild
	GuildID   string `firestore:"guild-id"`
	GuildName string `firestore:"guild-name"`
	// GuildAdminRoleID is the Premintbot role ID
	GuildAdminRoleID string `firestore:"guild-admin-role-id"`
	// GuildAdmins is a list of users who were in the audit log when the bot joined the guild
	GuildAdmins []string `firestore:"guild-admins"`

	// PREMINT settings
	// PremintAPIKey is the API key for the user's Premint project
	PremintAPIKey string `firestore:"premint-api-key"`
	// PremintRoleName is the ID of the role that is given to users who have registered for Premint
	PremintRoleID string `firestore:"premint-role-id"`
	// PremintRoleName is the name of the role that is given to users who have registered for Premint
	PremintRoleName string `firestore:"premint-role-name"`
}

type ConfigParams struct {
	Config *Guild
	doc    *firestore.DocumentSnapshot
}

func GetConfig(
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
		Config: config,
		doc:    docSnap,
	}
}

func GetConfigWithAPIKey(
	ctx context.Context,
	logger *zap.SugaredLogger,
	database *firestore.Client,
	apiKey string,
) *ConfigParams {
	// Fetch the config for the guild based on the premint-api-key
	docSnap, err := database.Collection("guilds").Where("premint-api-key", "==", apiKey).Limit(1).Documents(ctx).Next()
	if err != nil {
		logger.Errorw("Failed to get config", "api_key", apiKey, "error", err)
		return nil
	}

	// Get the config
	config := &Guild{}
	err = docSnap.DataTo(config)
	if err != nil {
		return nil
	}

	return &ConfigParams{
		Config: config,
		doc:    docSnap,
	}
}

package bot

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"github.com/mager/premintbot/premint"
	"go.uber.org/zap"
)

func premintCommand(
	ctx context.Context,
	logger *zap.SugaredLogger,
	database *firestore.Client,
	premintClient *premint.PremintClient,
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) {
	if m.Content != "!premint" {
		return
	}

	// Get the config
	p := getConfig(ctx, logger, database, m.GuildID)

	snowflake := m.Author.ID

	// Check Premint status
	status, err := premint.CheckPremintStatus(p.config.PremintAPIKey, snowflake)
	if err != nil {
		logger.Errorw("Failed to check premint status", "guild", m.GuildID, "error", err)
		return
	}

	var message string
	if status {
		message = "You are registered for Premint!"
	} else {
		message = "You are not registered for Premint"
	}

	// Send the message
	s.ChannelMessageSend(m.ChannelID, message)
}

package bot

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	// "github.com/kyokomi/emoji/v2"
)

func setupCommand(
	ctx context.Context,
	logger *zap.SugaredLogger,
	database *firestore.Client,
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) {
	if m.Content != "!premint-setup" {
		return
	}

	p := getConfig(ctx, logger, database, m.GuildID)

	// TODO: Make sure they have the admin role

	completed := false
	msg := "```\n"
	msg += "Here are the bot settings:\n\n"

	// Check if the guild is already setup
	if p.config.PremintAPIKey != "" {
		completed = true
		msg += "✅ Connected to project API Key: `" + p.config.PremintAPIKey + "`\n"
	} else {
		msg += "❌ Your Premint project is not connected, run !premint-set-api-key <API Key> to update it.\n"
	}

	if completed {
		msg += "\n✅ Your guild is setup!\n"
	} else {
		msg += "\nComplete the steps above.\n"
	}

	msg += "```"

	s.ChannelMessageSend(m.ChannelID, msg)
}

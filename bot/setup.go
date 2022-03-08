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
	if m.Content != "!setup" {
		return
	}

	p := getConfig(ctx, logger, database, m.GuildID)

	completed := false
	msg := "```\n"
	msg += "Here are the bot settings:\n\n"

	// Check if the guild is already setup
	if p.config.PremintAPIKey != "" {
		completed = true
		msg += "✅ Your Premint is connected: " + p.config.PremintAPIKey + "\n"
	} else {
		msg += "❌ Your Premint is not connected, run !set-premint <API Key> to update it.\n"
	}

	if completed {
		msg += "\n✅ Your guild is setup!\n"
	} else {
		msg += "\nComplete the steps above.\n"
	}

	msg += "```"

	s.ChannelMessageSend(m.ChannelID, msg)
}

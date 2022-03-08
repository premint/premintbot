package bot

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func helpCommand(
	ctx context.Context,
	logger *zap.SugaredLogger,
	database *firestore.Client,
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) {
	if m.Content != "!help" {
		return
	}

	p := getConfig(ctx, logger, database, m.GuildID)

	msg := "```\n" +
		"!help - Show this message\n" +
		"!premint - Check if you are a registered on Premint\n"

	if p.config.OwnerID == m.Author.ID {

		msg += "\nAdmin commands:\n\n" +
			"!setup - Show the bot settings\n" +
			"!set-premint <API Key> - Set the Premint project API Key\n" +
			"!nuke - Delete all channels and set the guild to inactive\n"
	}

	msg += "```"

	s.ChannelMessageSend(m.ChannelID, msg)
}

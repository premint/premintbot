package bot

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func premintNukeCommand(
	ctx context.Context,
	logger *zap.SugaredLogger,
	database *firestore.Client,
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) {
	if m.Content != "!premint-nuke" {
		return
	}

	logger.Infow("!premint-nuke called", zap.String("guild", m.GuildID), zap.String("user", m.Author.ID))

	cfg := GetConfig(ctx, logger, database, m.GuildID)
	g := getGuildFromMessage(s, m)

	// Make sure owner is sending the message
	if m.Author.ID != g.OwnerID {
		s.ChannelMessageSend(m.ChannelID, "You must be the guild owner to use this command.")
		return
	}

	// Delete channels
	for _, channel := range g.Channels {
		if channel.Name == premintConfigChannelName || channel.Name == premintCategoryName {
			s.ChannelDelete(channel.ID)
			logger.Infow("Deleted channel", zap.String("channel", channel.Name))
		}
	}

	// Delete roles
	for _, role := range g.Roles {
		if role.Name == "Premintbot" || role.Name == "Premint Bot" {
			s.GuildRoleDelete(g.ID, role.ID)
			logger.Infow("Deleted role", zap.String("role", role.Name))
		}
	}

	// Delete the record
	cfg.doc.Ref.Delete(ctx)
	logger.Infow("Deleted config record", zap.String("guild", m.GuildID))
}

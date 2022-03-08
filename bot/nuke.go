package bot

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func nukeCommand(
	ctx context.Context,
	logger *zap.SugaredLogger,
	database *firestore.Client,
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) {
	if m.Content != "!nuke" {
		return
	}

	p := getConfig(ctx, logger, database, m.GuildID)
	g := getGuildFromMessage(s, m)

	// Delete channels
	for _, channel := range g.Channels {
		if channel.Name == "portal-config" || channel.Name == "portal" {
			s.ChannelDelete(channel.ID)
			logger.Infow("Deleted channel", zap.String("channel", channel.Name))
		}
	}

	// Delete roles
	for _, role := range g.Roles {
		if role.Name == "premintbot" {
			s.GuildRoleDelete(g.ID, role.ID)
			logger.Infow("Deleted role", zap.String("role", role.Name))
		}
	}

	// Set to inactive
	p.doc.Ref.Update(ctx, []firestore.Update{{Path: "active", Value: false}})
}

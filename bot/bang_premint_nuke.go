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

	p := GetConfig(ctx, logger, database, m.GuildID)
	g := getGuildFromMessage(s, m)

	// Make sure owner is sending the message
	if m.Author.ID != g.OwnerID {
		s.ChannelMessageSend(m.ChannelID, "You must be the guild owner to use this command.")
		return
	}

	// Delete channels
	for _, channel := range g.Channels {
		if channel.Name == "premint-config" || channel.Name == "premint" {
			s.ChannelDelete(channel.ID)
			logger.Infow("Deleted channel", zap.String("channel", channel.Name))
		}
	}

	// Delete roles
	for _, role := range g.Roles {
		if role.Name == "Premintbot" {
			s.GuildRoleDelete(g.ID, role.ID)
			logger.Infow("Deleted role", zap.String("role", role.Name))
		}
	}

	// Set to inactive
	p.doc.Ref.Update(ctx, []firestore.Update{{Path: "active", Value: false}})
}

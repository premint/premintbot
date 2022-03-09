package bot

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// guildCreate is a function that is called when the bot joins a guild.
func guildCreate(ctx context.Context, logger *zap.SugaredLogger, database *firestore.Client) func(s *discordgo.Session, g *discordgo.GuildCreate) {
	return func(s *discordgo.Session, g *discordgo.GuildCreate) {
		// Return early if the portal channel category exists
		for _, category := range g.Guild.Channels {
			if category.Name == "premint" {
				return
			}
		}

		ownerID := g.Guild.OwnerID

		// Create premintbot role
		role, err := s.GuildRoleCreate(g.Guild.ID)
		if err != nil {
			logger.Errorw("Failed to create role", "guild", g.Guild.ID, "error", err)
			return
		}

		// Edit role name
		role, err = s.GuildRoleEdit(g.Guild.ID, role.ID, "Premintbot", 000000, false, 380910054518, false)
		if err != nil {
			logger.Errorw("Failed to edit role", "guild", g.Guild.ID, "error", err)
			return
		}

		// Add owner to role
		err = s.GuildMemberRoleAdd(g.Guild.ID, ownerID, role.ID)
		if err != nil {
			logger.Errorw("Failed to add owner to role", "guild", g.Guild.ID, "error", err)
			return
		}

		// Create Portal group
		permissionOverwrites := []*discordgo.PermissionOverwrite{
			// Allow for role
			{
				ID:   g.ID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: 0x0000000400,
			},
			// Hide for everyone else
			{
				ID:    role.ID,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: 0x0000000400,
			},
		}
		group, err := s.GuildChannelCreateComplex(
			g.Guild.ID,
			discordgo.GuildChannelCreateData{
				Type:                 discordgo.ChannelTypeGuildCategory,
				Name:                 "premint",
				PermissionOverwrites: permissionOverwrites,
			},
		)
		if err != nil {
			logger.Errorf("Failed to create channel: %v", err)
		}

		// Create #portal-config channel
		s.GuildChannelCreateComplex(
			g.Guild.ID,
			discordgo.GuildChannelCreateData{
				Type:                 discordgo.ChannelTypeGuildText,
				Name:                 "premint-config",
				ParentID:             group.ID,
				PermissionOverwrites: permissionOverwrites,
			},
		)

		// Add or update config in database
		docsnap, err := database.Collection("guilds").Doc(g.Guild.ID).Get(ctx)
		if err != nil {
			logger.Errorf("Failed to get guild: %v", err)
		}

		var guild Guild
		if docsnap.Exists() {
			err = docsnap.DataTo(&guild)
			if err != nil {
				logger.Errorf("Failed to decode guild: %v", err)
			}
			logger.Info("Guild exists in database")
			if guild.Active {
				logger.Info("Guild is active")
				return
			}
			guild.Active = true
		} else {
			logger.Info("Guild does not exist in database, creating")
			// Create the guild
			guild = Guild{
				Active:    true,
				GuildID:   g.Guild.ID,
				GuildName: g.Guild.Name,
				OwnerID:   ownerID,
				JoinedAt:  time.Now(),
			}

		}
		_, err = database.Collection("guilds").Doc(g.Guild.ID).Set(ctx, guild)
		if err != nil {
			logger.Errorf("Failed to create guild: %v", err)
		}
		logger.Info("Guild updated in database")

		// Announce in #general
		// for _, channel := range g.Guild.Channels {
		// 	if channel.Name == "general" {
		// 		s.ChannelMessageSendEmbed(channel.ID, createGeneralEmbed())
		// 	}
		// }
	}
}

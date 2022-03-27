package bot

import (
	"context"
	"time"

	bq "github.com/mager/premintbot/bigquery"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var (
	premintCategoryName      = "premint"
	premintConfigChannelName = "premint-config"
)

// guildCreate is a function that is called when the bot joins a guild.
func guildCreate(ctx context.Context, logger *zap.SugaredLogger, database *firestore.Client, bqClient *bigquery.Client) func(s *discordgo.Session, g *discordgo.GuildCreate) {
	return func(s *discordgo.Session, g *discordgo.GuildCreate) {
		var (
			ownerID   = g.Guild.OwnerID
			guildID   = g.Guild.ID
			guildName = g.Guild.Name
		)

		// Return early if the premint channel category exists
		for _, category := range g.Guild.Channels {
			if category.Name == premintCategoryName {
				logger.Infow("Premint category already exists", zap.String("guild", guildID), zap.String("guild", guildName))
				return
			}
		}

		docsnap, err := database.Collection("guilds").Doc(guildID).Get(ctx)
		if err != nil {
			logger.Errorf("Failed to get guild: %v", err)
		}

		var guild Guild
		if docsnap.Exists() {
			err = docsnap.DataTo(&guild)
			if err != nil {
				logger.Errorf("Failed to decode guild: %v", err)
			}
			logger.Infow("Guild exists in database", zap.String("guild", guildID), zap.String("guild", guildName))
			if !guild.Active {
				logger.Infow("Guild is inactive", zap.String("guild", guildID), zap.String("name", guildName))
			}
			return
		}

		// Create Premintbot role
		logger.Infow("Guild does not exist in database, creating", zap.String("guild", guildID), zap.String("name", guildName))
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

		// Fetch list of non-bot users in the audit log and add them as admins
		auditLog, err := s.GuildAuditLog(guildID, "", "", 28, 0)
		if err != nil {
			logger.Errorw("Failed to get audit log", "guild", g.Guild.ID, "error", err)
			return
		}
		guildAdmins := make([]string, 0)
		for _, entry := range auditLog.Users {
			if !entry.Bot {
				guildAdmins = append(guildAdmins, entry.ID)
			}
		}
		for _, userID := range guildAdmins {
			err = s.GuildMemberRoleAdd(g.Guild.ID, userID, role.ID)
			if err != nil {
				logger.Errorw("Failed to add user to role", "guild", g.Guild.ID, "user", userID, "error", err)
			}
		}

		// Add role to owner, just in case
		err = s.GuildMemberRoleAdd(g.Guild.ID, ownerID, role.ID)
		if err != nil {
			logger.Errorw("Failed to add owner to role", "guild", g.Guild.ID, "ownerID", ownerID, "error", err)
		}

		// Create Premint group
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
				Name:                 premintCategoryName,
				PermissionOverwrites: permissionOverwrites,
			},
		)
		if err != nil {
			logger.Errorf("Failed to create channel: %v", err)
		}

		// Create #premint-config channel
		c, err := s.GuildChannelCreateComplex(
			g.Guild.ID,
			discordgo.GuildChannelCreateData{
				Type:                 discordgo.ChannelTypeGuildText,
				Name:                 premintConfigChannelName,
				ParentID:             group.ID,
				PermissionOverwrites: permissionOverwrites,
			},
		)
		if err != nil {
			logger.Errorf("Failed to create channel: %v", err)
		}

		// Create the guild
		joinedAt := time.Now()
		guild = Guild{
			Active:           true,
			GuildID:          g.Guild.ID,
			GuildName:        g.Guild.Name,
			GuildAdminRoleID: role.ID,
			GuildOwnerID:     ownerID,
			JoinedAt:         joinedAt,
			GuildAdmins:      guildAdmins,
		}

		_, err = database.Collection("guilds").Doc(g.Guild.ID).Set(ctx, guild)
		if err != nil {
			logger.Errorf("Failed to create guild in Firestore: %v", err)
		}
		logger.Infow("Guild updated in database", zap.String("guild", guildID), zap.String("guild", guildName))

		// Send an event to BigQuery
		evt := &bq.BQGuildsCreate{
			GuildID:          g.Guild.ID,
			Timestamp:        joinedAt,
			GuildName:        g.Guild.Name,
			GuildAdminRoleID: role.ID,
			OwnerID:          ownerID,
			GuildAdmins:      guildAdmins,
		}
		bq.RecordGuildsCreate(bqClient, evt)

		s.ChannelMessageSendEmbed(c.ID, createGeneralEmbed())
	}
}

func createGeneralEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Premint Bot",
		Description: "Hello! My name is Premint Bot.\n\n" + "To get started, you need to set up the bot.\n\n" + "1. Set your Premint API Key with `!premint-set-api-key <API Key>`. You can find your API key on the Premint website: https://www.premint.xyz/dashboard/. Click on a project, then click Edit Settings, then API.\n" + "2. Set your Premint Role with `!premint-set-role <Role ID>`\n" + "3. Run `!premint-setup` to check everything is set up correctly.\n\n_Run `!premint` for a list of commands._",
		Color:       0x00ff00,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://cdn.discordapp.com/avatars/420864490981227266/b7f9f9a9c7b6e5e6f7e8f8c1b7f1f2d6.png?size=2048",
		},
	}
}

package bot

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	// "github.com/kyokomi/emoji/v2"
)

func premintSetupCommand(
	ctx context.Context,
	logger *zap.SugaredLogger,
	database *firestore.Client,
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) {
	if m.Content != "!premint-setup" {
		return
	}

	logger.Infow("!premint-setup called", zap.String("guild", m.GuildID), zap.String("user", m.Author.ID))

	cfg := GetConfig(ctx, logger, database, m.GuildID)

	// TODO: Make sure they have the admin role

	apiKeySet := cfg.Config.PremintAPIKey != ""
	roleSet := cfg.Config.PremintRoleID != "" && cfg.Config.PremintRoleName != ""
	completed := apiKeySet && roleSet

	// Check if the Premint API key is set
	apiKeyField := &discordgo.MessageEmbedField{}
	if apiKeySet {
		apiKeyField.Name = "✅ Connected project API Key"
		apiKeyField.Value = "`" + cfg.Config.PremintAPIKey + "`"
	} else {
		apiKeyField.Name = "❌ Missing project API Key"
		apiKeyField.Value = "Use `!premint-set-api-key PREMINT_API_KEY` to set it. You can find your API key on the Premint website: https://www.premint.xyz/dashboard/. Click on a project, then click Edit Settings > API."
	}

	// Check if the role is set
	roleField := &discordgo.MessageEmbedField{}
	if roleSet {
		roleField.Name = "✅ Role is set"
		roleField.Value = "`" + cfg.Config.PremintRoleName + "`"
	} else {
		roleField.Name = "❌ Role is not set"
		roleField.Value = "Use `!premint-set-role DISCORD_ROLE_ID` to set it. Create a role that you want your users to get when they use the `/premint` command. You can find your Discord role ID by going to Server Settings > Roles > Right click the role > Copy ID."
	}

	fields := []*discordgo.MessageEmbedField{apiKeyField, roleField}
	s.ChannelMessageSendEmbed(m.ChannelID, createSetupEmbed(fields, completed))
}

func createSetupEmbed(fields []*discordgo.MessageEmbedField, completed bool) *discordgo.MessageEmbed {
	color := 0x00ff00
	description := "Your guild is setup!"
	if !completed {
		color = 0xff0000
		description = "Complete the steps below."
	}
	return &discordgo.MessageEmbed{
		Type:        "",
		Title:       "Premint Setup",
		Description: description,
		Color:       color,
		Fields:      fields,
	}
}

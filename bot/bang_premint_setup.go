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

	p := GetConfig(ctx, logger, database, m.GuildID)

	// TODO: Make sure they have the admin role

	apiKeySet := p.Config.PremintAPIKey != ""
	roleSet := p.Config.PremintRoleID != "" && p.Config.PremintRoleName != ""
	completed := apiKeySet && roleSet

	// Check if the Premint API key is set
	apiKeyField := &discordgo.MessageEmbedField{}
	if apiKeySet {
		apiKeyField.Name = "✅ Connected project API Key"
		apiKeyField.Value = "`" + p.Config.PremintAPIKey + "`"
	} else {
		apiKeyField.Name = "❌ Missing project API Key"
		apiKeyField.Value = "Use `!premint-set-api-key <API Key>` to set it. You can find your API key on the Premint website: https://www.premint.xyz/dashboard/. Click on a project, then click Edit Settings, then API."
	}

	// Check if the role is set
	roleField := &discordgo.MessageEmbedField{}
	if roleSet {
		roleField.Name = "✅ Role is set"
		roleField.Value = "`" + p.Config.PremintRoleName + "`"
	} else {
		roleField.Name = "❌ Role is not set"
		roleField.Value = "Use `!premint-set-role <Role ID>` to set it. You can find your role ID but right clicking it > Copy ID."
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

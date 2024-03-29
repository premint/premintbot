package bot

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func premintCommand(
	ctx context.Context,
	logger *zap.SugaredLogger,
	database *firestore.Client,
	s *discordgo.Session,
	m *discordgo.MessageCreate,
) {
	if m.Content != "!premint" {
		return
	}

	logger.Infow("!premint called", zap.String("guild", m.GuildID), zap.String("user", m.Author.ID))

	cfg := GetConfig(ctx, logger, database, m.GuildID)
	g := getGuildFromMessage(s, m)

	// Find #premint-config channel
	premintConfigChannel := ""
	for _, channel := range g.Channels {
		if channel.Name == premintConfigChannelName {
			premintConfigChannel = channel.ID
		}
	}

	if isAdmin(cfg, m.Author) && m.ChannelID == premintConfigChannel {
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, CreateAdminHelpEmbed())
		if err != nil {
			logger.Errorw("Failed to send help message", "guild", m.GuildID, "error", err)
			return
		}
	} else {
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, CreatePublicHelpEmbed())
		if err != nil {
			logger.Errorw("Failed to send help message", "guild", m.GuildID, "error", err)
			return
		}
	}
}

// CreatePublicHelpEmbed creates a response for the !help command
func CreatePublicHelpEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "PREMINT Bot Help",
		Description: "Here is what you can do with your PREMINT Bot:",
		Color:       0x00ffff,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "❓ If you have any questions, please ask in the PREMINT Discord.",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Collector Commands",
				Value:  "`/premint` - This will check to see if the Discord user is registered on the PREMINT list. If they are, it will return their wallet ID.\n`/premint {ETH wallet address or ENS name}` - This will check if the wallet address is on the PREMINT list.\n`!premint` - Show this message",
				Inline: false,
			},
			// TODO: Support aliases
			// {
			// 	Name:   "Aliases",
			// 	Value:  "_Note: Users can also use `/accesslist`, `/allowlist` or `/whitelist`_",
			// 	Inline: false,
			// },
		},
	}
}

// CreateAdminHelpEmbed creates a response for the !help command
func CreateAdminHelpEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "PREMINT Bot Help",
		Description: "Here is what you can do with your PREMINT Bot:",
		Color:       0x00ffff,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "❓ If you have any questions, please ask in the PREMINT Discord. Full setup instructions on Github: https://github.com/premint/premintbot/blob/main/SETUP.md",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Admin Setup",
				Value:  "`!premint-setup` - Show the bot settings\n`!premint-set-api-key PREMINT_API_KEY` - This connects the bot to a specific PREMINT project. You can find your API Key in the projects Settings > API\n`!premint-set-role DISCORD_ROLE_ID` - Set the role you want your users to get when they are registered with Premint\n`!premint-nuke` - Delete all channels and set the guild to inactive",
				Inline: false,
			},
		},
	}
}

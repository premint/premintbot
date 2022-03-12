package bot

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"github.com/mager/premintbot/premint"
	"go.uber.org/zap"
)

func premintSlashCommand(ctx context.Context, logger *zap.SugaredLogger, database *firestore.Client, premintClient *premint.PremintClient) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		p := getConfig(ctx, logger, database, i.GuildID)

		if p.config.PremintAPIKey == "" {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Premint API key is not set. Please use `!set-premint` command to set it.",
				},
			})
			return
		}

		var status bool
		var err error
		var message string

		cmdData := i.ApplicationCommandData()
		if cmdData.Options == nil {
			status, err = premint.CheckPremintStatusForUser(p.config.PremintAPIKey, i.Interaction.Member.User.ID)
			if err != nil {
				logger.Errorw("Failed to check premint status", "guild", i.GuildID, "error", err)
				return
			}
		} else {
			// TODO: Validate ETH address
			address := i.ApplicationCommandData().Options[0].StringValue()
			status, err = premint.CheckPremintStatusForAddress(p.config.PremintAPIKey, address)
			if err != nil {
				logger.Errorw("Failed to check premint status", "guild", i.GuildID, "error", err)
				return
			}
		}

		if status {
			message = "You are registered for Premint!"
		} else {
			message = "You are not registered for Premint"
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		})
	}
}

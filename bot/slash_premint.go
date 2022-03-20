package bot

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"github.com/mager/premintbot/premint"
	"go.uber.org/zap"
)

func premintSlashCommand(ctx context.Context, logger *zap.SugaredLogger, database *firestore.Client, premintClient *premint.PremintClient) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		p := GetConfig(ctx, logger, database, i.GuildID)

		if p.Config.PremintAPIKey == "" {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Premint API key is not set. Please use `!premint-set-api-key <Premint API key>` command to set it.",
				},
			})
			return
		}

		var err error
		var message string
		cmdData := i.ApplicationCommandData()

		resp := premint.CheckPremintStatusResp{}
		if cmdData.Options == nil {
			logger.Info("Checking premint status with the Discord user ID")
			resp, err = premint.CheckPremintStatusForUser(p.Config.PremintAPIKey, i.Interaction.Member.User.ID)
			if err != nil {
				logger.Errorw("Failed to check premint status", "guild", i.GuildID, "error", err)
				return
			}
		} else {
			// TODO: Validate ETH address
			logger.Info("Checking premint status with the ETH wallet address")
			address := i.ApplicationCommandData().Options[0].StringValue()
			resp, err = premint.CheckPremintStatusForAddress(p.Config.PremintAPIKey, address)
			if err != nil {
				logger.Errorw("Failed to check premint status", "guild", i.GuildID, "error", err)
				return
			}
		}

		if resp.Registered {
			message = fmt.Sprintf("✅ Wallet %s is registered on the %s list. %s", resp.WalletAddress, resp.ProjectName, resp.ProjectURL)
		} else {
			message = fmt.Sprintf("❌ Wallet %s is not registered on the %s list. %s", resp.WalletAddress, resp.ProjectName, resp.ProjectURL)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
				// Ephemeral = this message is only visible to the user who invoked the Interaction,
				Flags: 64,
			},
		})

	}
}

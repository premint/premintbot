package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	bq "github.com/premint/premintbot/bigquery"
	"github.com/premint/premintbot/infura"
	"github.com/premint/premintbot/premint"
	"go.uber.org/zap"
)

func premintSlashCommand(ctx context.Context, logger *zap.SugaredLogger, database *firestore.Client, premintClient *premint.PremintClient, bqClient *bigquery.Client, infuraClient *infura.InfuraClient) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		p := GetConfig(ctx, logger, database, i.GuildID)

		if p.Config.PremintAPIKey == "" {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Premint API key is not set. Please ask an admin to run the `!premint-set-api-key PREMINT_API_KEY` command to set it.",
				},
			})
			return
		}

		var err error
		var message string
		cmdData := i.ApplicationCommandData()

		resp := premint.CheckPremintStatusResp{}
		withAddress := false
		address := ""
		if cmdData.Options == nil {
			logger.Info("Checking premint status with the Discord user ID")
			resp, err = premint.CheckPremintStatusForUser(logger, p.Config.PremintAPIKey, i.Interaction.Member.User.ID)
			if err != nil {
				logger.Errorw("Failed to check premint status", "guild", i.GuildID, "error", err)
				return
			}
		} else {
			// TODO: Validate ETH address
			withAddress = true
			logger.Info("Checking premint status with the ETH wallet address")
			addressOption := i.ApplicationCommandData().Options[0].StringValue()
			// Check if the address is actually an ENS name
			if !strings.HasPrefix(addressOption, "0x") {
				address = infuraClient.GetAddressFromENSName(addressOption)
				logger.Infow("Address is an ENS name", "resolved", address)
			} else {
				address = addressOption
			}

			resp, err = premint.CheckPremintStatusForAddress(logger, p.Config.PremintAPIKey, address)
			if err != nil {
				logger.Errorw("Failed to check premint status", "guild", i.GuildID, "error", err)
				return
			}
		}

		evt := &bq.BQSlashPremint{
			Address:     address,
			GuildID:     i.GuildID,
			UserID:      i.Interaction.Member.User.ID,
			Timestamp:   time.Now(),
			WithAddress: withAddress,
			Registered:  resp.Registered,
		}
		bq.RecordSlashPremint(bqClient, evt)

		if resp.Registered {
			message = fmt.Sprintf("✅ Wallet %s is registered on the %s list. %s", resp.WalletAddress, resp.ProjectName, resp.ProjectURL)
			roleSet := false
			for _, role := range i.Interaction.Member.Roles {
				if role == p.Config.PremintRoleID {
					roleSet = true
					break
				}
			}

			if p.Config.PremintRoleID != "" && !roleSet {
				err = s.GuildMemberRoleAdd(i.GuildID, i.Interaction.Member.User.ID, p.Config.PremintRoleID)
				if err != nil {
					logger.Errorw("Failed to add  role", "guild", i.GuildID, "error", err)
					return
				}
			}
		} else {
			if withAddress {
				message = fmt.Sprintf("❌ Wallet %s is not registered on the %s list. %s", resp.WalletAddress, resp.ProjectName, resp.ProjectURL)
			} else {
				message = fmt.Sprintf("❌ User %s is not registered on the %s list. %s", i.Interaction.Member.User.Username, resp.ProjectName, resp.ProjectURL)
			}
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

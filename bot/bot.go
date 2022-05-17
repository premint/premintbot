package bot

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"github.com/premint/premintbot/config"
	"github.com/premint/premintbot/infura"
	"github.com/premint/premintbot/premint"
	"go.uber.org/zap"
)

var (
	// /premint slash command
	premintCmd = &discordgo.ApplicationCommand{
		Name:        "premint",
		Description: "Check if your address is registered with Premint",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "address",
				Description: "Your ETH address or ENS name",
				Required:    false,
			},
		},
	}
)

func Start(
	cfg config.Config,
	dg *discordgo.Session,
	logger *zap.SugaredLogger,
	database *firestore.Client,
	premintClient *premint.PremintClient,
	bqClient *bigquery.Client,
	infuraClient *infura.InfuraClient,
) {
	ctx := context.Background()

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate(ctx, logger, database, premintClient, bqClient))

	// Register the guildCreate func as a callback for GuildCreate events.
	dg.AddHandler(guildCreate(ctx, logger, database, bqClient))

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logger.Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	dg.ApplicationCommandCreate(cfg.DiscordAppID, "", premintCmd)
	dg.AddHandler(premintSlashCommand(ctx, logger, database, premintClient, bqClient, infuraClient))

	// Open a websocket connection to Discord and begin listening.
	wsErr := dg.Open()
	if wsErr != nil {
		fmt.Println("error opening connection,", wsErr)
	}

	// DEBUGGING
	// Register the auditLogChange func as a callback for auditLog events.
	// dg.AddHandler(auditLogChange(ctx, logger, database, bqClient))
	// Register the auditLogUpdate func as a callback for auditLog events.
	// dg.AddHandler(auditLogUpdate(ctx, logger, database, bqClient))
}

func messageCreate(ctx context.Context, logger *zap.SugaredLogger, database *firestore.Client, premintClient *premint.PremintClient, bqClient *bigquery.Client) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Admin commands
		premintNukeCommand(ctx, logger, database, s, m)
		premintSetupCommand(ctx, logger, database, s, m)
		premintSetAPIKeyCommand(ctx, logger, database, bqClient, s, m)
		premintSetRoleCommand(ctx, logger, database, bqClient, s, m)

		// Public commands
		premintCommand(ctx, logger, database, s, m)
	}
}

func getGuildFromMessage(s *discordgo.Session, m *discordgo.MessageCreate) *discordgo.Guild {
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		guild, err = s.Guild(m.GuildID)
		if err != nil {
			return nil
		}
	}
	return guild
}

// isAdmin checks if the user is an admin
func isAdmin(p *ConfigParams, u *discordgo.User) bool {
	for _, admin := range p.Config.GuildAdmins {
		if admin == u.ID {
			return true
		}
	}

	return false
}

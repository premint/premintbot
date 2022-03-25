package bot

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"github.com/mager/premintbot/premint"
	"go.uber.org/zap"
)

var (
	appID                 = "950933570564800552"
	integerOptionMinValue = 1.0
	premintCmd            = &discordgo.ApplicationCommand{
		Name:        "premint",
		Description: "Check if your address is registered with Premint",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "address",
				Description: "Your ETH address",
				Required:    false,
			},
		},
	}
	setupPremintCmd = &discordgo.ApplicationCommand{
		Name:        "setup-premint",
		Description: "Add your Premint API key to the bot",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "API Key",
				Description: "Your Premint API key",
				Required:    true,
			},
		},
	}
)

func Start(
	dg *discordgo.Session,
	logger *zap.SugaredLogger,
	database *firestore.Client,
	premintClient *premint.PremintClient,
	bqClient *bigquery.Client,
) {
	ctx := context.Background()
	logger.Info("Registering Discord bot")

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate(ctx, logger, database, premintClient))

	// Register the guildCreate func as a callback for GuildCreate events.
	dg.AddHandler(guildCreate(ctx, logger, database, bqClient))

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logger.Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	dg.AddHandler(premintSlashCommand(ctx, logger, database, premintClient))
	dg.ApplicationCommandCreate(appID, "", premintCmd)

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
	// https://github.com/bwmarrin/discordgo/blob/v0.23.2/structs.go#L1295
	// dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds

	// Wait here until CTRL-C or other term signal is received.
	// fmt.Println("Bot is now running. Press CTRL-C to exit.")
	// sc := make(chan os.Signal, 1)
	// signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	// <-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(ctx context.Context, logger *zap.SugaredLogger, database *firestore.Client, premintClient *premint.PremintClient) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Admin commands
		premintNukeCommand(ctx, logger, database, s, m)
		premintSetupCommand(ctx, logger, database, s, m)
		premintSetAPIKeyCommand(ctx, logger, database, s, m)
		premintSetRoleCommand(ctx, logger, database, s, m)

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

func isAdmin(p *ConfigParams, a discordgo.Member) bool {
	return false
}

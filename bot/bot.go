package bot

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"github.com/mager/premintbot/premint"
	"go.uber.org/zap"
)

var (
	integerOptionMinValue = 1.0
	premintCmd            = &discordgo.ApplicationCommand{
		Name:        "pppp",
		Description: "Command for demonstrating options",
		Options: []*discordgo.ApplicationCommandOption{

			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "address",
				Description: "String option",
				Required:    false,
			},
		},
	}
	premintCmdHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		margs := []interface{}{
			// Here we need to convert raw interface{} value to wanted type.
			// Also, as you can see, here is used utility functions to convert the value
			// to particular type. Yeah, you can use just switch type,
			// but this is much simpler
			i.ApplicationCommandData().Options[0].StringValue(),
		}
		msgformat :=
			` Now you just learned how to use command options. Take a look to the value of which you've just entered:
				> string_option: %s
`
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Ignore type for now, we'll discuss them in "responses" part
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(
					msgformat,
					margs...,
				),
			},
		})
	}
)

func Start(
	dg *discordgo.Session,
	logger *zap.SugaredLogger,
	database *firestore.Client,
	premintClient *premint.PremintClient,
) {
	ctx := context.Background()
	logger.Info("Registering Discord bot")

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate(ctx, logger, database, premintClient))

	// Register the guildCreate func as a callback for GuildCreate events.
	dg.AddHandler(guildCreate(ctx, logger, database))

	// https://github.com/bwmarrin/discordgo/blob/v0.23.2/structs.go#L1295
	// dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds

	// Slash commands
	// dg.AddHandler(slashCommand(ctx, logger, database, premintClient))

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logger.Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.User == nil {
			logger.Infof("Interaction created: %v", i.ID)
			premintCmdHandler(s, i)
		}
	})

	ccmd, err := dg.ApplicationCommandCreate("950933570564800552", "", premintCmd)
	if err != nil {
		logger.Panicf("Cannot create '%v' command: %v", premintCmd.Name, err)
	} else {
		logger.Infof("Created '%v' command", premintCmd.Name)
	}

	// Open a websocket connection to Discord and begin listening.
	wsErr := dg.Open()
	if wsErr != nil {
		fmt.Println("error opening connection,", wsErr)
	}

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
		nukeCommand(ctx, logger, database, s, m)
		setupCommand(ctx, logger, database, s, m)
		setPremintCommand(ctx, logger, database, s, m)

		// Public commands
		helpCommand(ctx, logger, database, s, m)
		premintCommand(ctx, logger, database, premintClient, s, m)
	}
}

func createGeneralEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Premint Bot",
		Description: "Hello! My name is Premint Bot. I am a bot that helps you manage your Discord server.",
		Color:       0x00ff00,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://cdn.discordapp.com/avatars/420864490981227266/b7f9f9a9c7b6e5e6f7e8f8c1b7f1f2d6.png?size=2048",
		},
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

func slashCommand(ctx context.Context, logger *zap.SugaredLogger, database *firestore.Client, premintClient *premint.PremintClient) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		address := i.ApplicationCommandData().Options[0].StringValue()
		logger.Info("Slash command", zap.String("address", address))
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Ignore type for now, we'll discuss them in "responses" part
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Testing",
			},
		})
	}
}

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
	commands              = []*discordgo.ApplicationCommand{
		{
			Name:        "premint",
			Description: "Command for demonstrating options",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "string-option",
					Description: "String option",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "integer-option",
					Description: "Integer option",
					MinValue:    &integerOptionMinValue,
					MaxValue:    10,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "number-option",
					Description: "Float option",
					MaxValue:    10.1,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "bool-option",
					Description: "Boolean option",
					Required:    true,
				},

				// The same concept applies to Discord's Slash-commands API

				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel-option",
					Description: "Channel option",
					// Channel type mask
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
						discordgo.ChannelTypeGuildVoice,
					},
					Required: false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-option",
					Description: "User option",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role-option",
					Description: "Role option",
					Required:    false,
				},
			},
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"premint": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			margs := []interface{}{
				// Here we need to convert raw interface{} value to wanted type.
				// Also, as you can see, here is used utility functions to convert the value
				// to particular type. Yeah, you can use just switch type,
				// but this is much simpler
				i.ApplicationCommandData().Options[0].StringValue(),
				i.ApplicationCommandData().Options[1].IntValue(),
				i.ApplicationCommandData().Options[2].FloatValue(),
				i.ApplicationCommandData().Options[3].BoolValue(),
			}
			msgformat :=
				` Now you just learned how to use command options. Take a look to the value of which you've just entered:
				> string_option: %s
				> integer_option: %d
				> number_option: %f
				> bool_option: %v
`
			if len(i.ApplicationCommandData().Options) >= 5 {
				margs = append(margs, i.ApplicationCommandData().Options[4].ChannelValue(nil).ID)
				msgformat += "> channel-option: <#%s>\n"
			}
			if len(i.ApplicationCommandData().Options) >= 6 {
				margs = append(margs, i.ApplicationCommandData().Options[5].UserValue(nil).ID)
				msgformat += "> user-option: <@%s>\n"
			}
			if len(i.ApplicationCommandData().Options) >= 7 {
				margs = append(margs, i.ApplicationCommandData().Options[6].RoleValue(nil, "").ID)
				msgformat += "> role-option: <@&%s>\n"
			}
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
		},
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
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds

	// Slash commands
	// dg.AddHandler(slashCommand(ctx, logger, database, premintClient))

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logger.Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	for _, v := range commands {
		_, err := dg.ApplicationCommandCreate("950933570564800552", "", v)
		if err != nil {
			logger.Panicf("Cannot create '%v' command: %v", v.Name, err)
		} else {
			logger.Infof("Created '%v' command", v.Name)
		}
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

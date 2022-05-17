package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
	magicClient "github.com/magiclabs/magic-admin-go/client"
	bq "github.com/premint/premintbot/bigquery"

	"github.com/premint/premintbot/bot"
	"github.com/premint/premintbot/config"
	"github.com/premint/premintbot/database"
	"github.com/premint/premintbot/handler"
	"github.com/premint/premintbot/infura"
	"github.com/premint/premintbot/logger"
	"github.com/premint/premintbot/magic"
	"github.com/premint/premintbot/premint"
	"github.com/premint/premintbot/router"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(
			bq.Options,
			config.Options,
			database.Options,
			infura.Options,
			logger.Options,
			magic.Options,
			premint.Options,
			router.Options,
		),
		fx.Invoke(Register),
	).Run()
}

func Register(
	lc fx.Lifecycle,
	bqClient *bigquery.Client,
	cfg config.Config,
	database *firestore.Client,
	infuraClient *infura.InfuraClient,
	logger *zap.SugaredLogger,
	magic *magicClient.API,
	premintClient *premint.PremintClient,
	router *mux.Router,
) {
	// Setup Discord Bot
	token := fmt.Sprintf("Bot %s", cfg.DiscordAuthToken)
	dg, err := discordgo.New(token)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Start the bot
	bot.Start(cfg, dg, logger, database, premintClient, bqClient, infuraClient)

	// Route handler
	handler.New(logger, router, database, dg)

	// Cleanly close down the Discord session.
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing Discord session")
			defer dg.Close()
			if err != nil {
				logger.Errorf("Failed to close Discord session: %v", err)
			}
			return err
		},
	})
}

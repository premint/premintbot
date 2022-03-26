package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
	bq "github.com/mager/premintbot/bigquery"

	"github.com/mager/premintbot/bot"
	"github.com/mager/premintbot/config"
	"github.com/mager/premintbot/database"
	"github.com/mager/premintbot/handler"
	"github.com/mager/premintbot/logger"
	"github.com/mager/premintbot/premint"
	"github.com/mager/premintbot/router"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(
			bq.Options,
			config.Options,
			database.Options,
			logger.Options,
			router.Options,
			premint.Options,
		),
		fx.Invoke(Register),
	).Run()
}

func Register(
	lc fx.Lifecycle,
	bqClient *bigquery.Client,
	cfg config.Config,
	database *firestore.Client,
	logger *zap.SugaredLogger,
	router *mux.Router,
	premintClient *premint.PremintClient,
) {
	// Setup Discord Bot
	token := fmt.Sprintf("Bot %s", cfg.DiscordAuthToken)
	dg, err := discordgo.New(token)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	bot.Start(dg, logger, database, premintClient, bqClient)
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

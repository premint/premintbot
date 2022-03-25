package bot

import (
	"context"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"github.com/kr/pretty"
	"go.uber.org/zap"
)

// auditLogUpdate is a function that is called when audit log is updated
func auditLogUpdate(ctx context.Context, logger *zap.SugaredLogger, database *firestore.Client, bqClient *bigquery.Client) func(s *discordgo.Session, a *discordgo.AuditLogAction) {
	return func(s *discordgo.Session, a *discordgo.AuditLogAction) {
		logger.Info("AUDIT LOG UPDATED")

		pretty.Print(a)
	}
}

// auditLogChange is a function that is called when audit log is changed
func auditLogChange(ctx context.Context, logger *zap.SugaredLogger, database *firestore.Client, bqClient *bigquery.Client) func(s *discordgo.Session, a *discordgo.AuditLogChange) {
	return func(s *discordgo.Session, a *discordgo.AuditLogChange) {
		logger.Info("AUDIT LOG CHANGED")

		pretty.Print(a)
	}
}

package bigquery

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
)

type BQGuildsCreate struct {
	GuildID          string
	UserID           string
	Timestamp        time.Time
	GuildName        string
	GuildAdminRoleID string
	GuildAdmins      []string
	OwnerID          string
}

type BQAdminErrors struct {
	GuildID      string
	UserID       string
	Timestamp    time.Time
	ActionType   string
	ErrorMessage string
}

func RecordGuildsCreate(
	bq *bigquery.Client,
	evt *BQGuildsCreate,
) {
	var (
		ctx   = context.Background()
		table = bq.DatasetInProject(projectID, "guilds").Table("create")
		u     = table.Inserter()
		items = []*BQGuildsCreate{evt}
	)
	if err := u.Put(ctx, items); err != nil {
		log.Fatalf("Failed to insert: %v", err)
	}
}

func RecordAdminErrors(
	bq *bigquery.Client,
	evt *BQAdminErrors,
) {
	var (
		ctx   = context.Background()
		table = bq.DatasetInProject(projectID, "guilds").Table("admin_errors")
		u     = table.Inserter()
		items = []*BQAdminErrors{evt}
	)
	if err := u.Put(ctx, items); err != nil {
		log.Fatalf("Failed to insert: %v", err)
	}
}

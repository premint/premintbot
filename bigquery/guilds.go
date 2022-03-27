package bigquery

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/bwmarrin/discordgo"
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

type BQAdminAction struct {
	GuildID    string
	UserID     string
	Timestamp  time.Time
	ActionType string
	Message    string
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

func RecordAdminAction(
	bq *bigquery.Client,
	m *discordgo.MessageCreate,
	actionType, msg string,
) {
	var (
		ctx   = context.Background()
		table = bq.DatasetInProject(projectID, "guilds").Table("admin_actions")
		u     = table.Inserter()
		items = []*BQAdminAction{
			{
				GuildID:    m.GuildID,
				UserID:     m.Author.ID,
				Timestamp:  time.Now(),
				ActionType: actionType,
				Message:    msg,
			},
		}
	)
	if err := u.Put(ctx, items); err != nil {
		log.Fatalf("Failed to insert: %v", err)
	}
}

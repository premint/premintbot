package bigquery

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
)

type BQGuildsCreate struct {
	GuildID          string
	GuildName        string
	GuildAdminRoleID string
	GuildAdmins      []string
	OwnerID          string
	JoinedAt         time.Time
}

var projectID = "premint-343516"

// ProvideBQ provides a bigquery client
func ProvideBQ() *bigquery.Client {
	client, err := bigquery.NewClient(context.TODO(), projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

var Options = ProvideBQ

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

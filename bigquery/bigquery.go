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
	OwnerID          string
	JoinedAt         time.Time
}

// ProvideBQ provides a bigquery client
func ProvideBQ() *bigquery.Client {
	projectID := "premint-343516"

	client, err := bigquery.NewClient(context.TODO(), projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

var Options = ProvideBQ

func RecordGuildsCreate(
	bq *bigquery.Client,
	guildID string,
	guildName string,
	guildAdminRoleID string,
	ownerID string,
	joinedAt time.Time,
) {
	var (
		ctx   = context.Background()
		table = bq.DatasetInProject("premint-343516", "guilds").Table("create")
		u     = table.Inserter()

		items = []*BQGuildsCreate{
			{
				GuildID:          guildID,
				GuildName:        guildName,
				GuildAdminRoleID: guildAdminRoleID,
				OwnerID:          ownerID,
				JoinedAt:         joinedAt,
			},
		}
	)
	if err := u.Put(ctx, items); err != nil {
		log.Fatalf("Failed to insert: %v", err)
	}
}

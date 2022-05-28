package bigquery

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
)

type BQSlashPremint struct {
	Address     string
	GuildID     string
	UserID      string
	Timestamp   time.Time
	Registered  bool
	WithAddress bool
}

func RecordSlashPremint(
	bq *bigquery.Client,
	evt *BQSlashPremint,
) {
	var (
		ctx   = context.Background()
		table = bq.DatasetInProject(projectID, "commands").Table("slash_premint")
		u     = table.Inserter()
		items = []*BQSlashPremint{evt}
	)
	if err := u.Put(ctx, items); err != nil {
		log.Fatalf("Failed to insert: %v", err)
	}
}

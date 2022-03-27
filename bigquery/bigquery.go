package bigquery

import (
	"context"
	"log"

	"cloud.google.com/go/bigquery"
)

// TODO: Move to config
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

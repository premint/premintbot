package database

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/mager/premintbot/config"
)

// ProvideDB provides a firestore client
func ProvideDB(cfg config.Config) *firestore.Client {
	projectID := cfg.GoogleCloudProject

	client, err := firestore.NewClient(context.TODO(), projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

var Options = ProvideDB

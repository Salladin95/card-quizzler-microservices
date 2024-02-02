package fireBase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/Salladin95/card-quizzler-microservices/auth-service/cmd/api/config"
	"google.golang.org/api/option"
	"log"
	"os"
	"path/filepath"
)

// NewFireBaseApp creates a new Firebase App based on the provided configuration.
func NewFireBaseApp(cfg config.FireBaseCfg) *firebase.App {
	// Get the current working directory
	root, err := os.Getwd()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Create a background context
	ctx := context.Background()

	// Specify the options for initializing the Firebase App, including the path to the service account key file
	opt := option.WithCredentialsFile(filepath.Join(root, cfg.FireBaseAccKey))

	// Initialize the Firebase App
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: cfg.FireBaseProjectId,
	}, opt)
	if err != nil {
		// Log and exit if there's an error initializing the Firebase App
		log.Fatalf("error initializing Firebase App: %v\n", err)
	}

	// Return the initialized Firebase App
	return app
}

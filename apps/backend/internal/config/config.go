package config

import "os"

const defaultProjectID = "sns-only-event-local"

type Config struct {
	ProjectID string
}

func Load() Config {
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		projectID = os.Getenv("PUBSUB_PROJECT_ID")
	}
	if projectID == "" {
		projectID = defaultProjectID
	}

	return Config{
		ProjectID: projectID,
	}
}

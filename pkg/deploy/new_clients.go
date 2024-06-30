package deploy

import (
	"context"
	"fmt"

	artifactregistry "cloud.google.com/go/artifactregistry/apiv1"
	run "cloud.google.com/go/run/apiv2"
	"google.golang.org/api/option"
)

type client struct {
	run              *run.ExecutionsClient
	artifactRegistry *artifactregistry.Client
	project          string
	region           string
	credentialsFile  string
}

// NewClient creates a new client struct with the provided region, project, and credentials file.
// If the credentials file is provided, the clients will be created with the credentials file.
// If the credentials file is not provided, the clients will be created with the application default credentials.
// If the clients cannot be created, the function will log a fatal error.
// The function returns a pointer to the client struct.
func NewClient(project, region, credentialsFile string) (*client, error) {
	clients := &client{
		project:         project,
		region:          region,
		credentialsFile: credentialsFile,
	}

	if err := clients.setClients(); err != nil {
		return nil, fmt.Errorf("setClients(): %w", err)
	}

	return clients, nil
}

func (c *client) setClients() error {
	var err error
	ctx := context.TODO()

	// Generate the client options with the credentials file
	// if the credentials file is provided
	var clientOptions []option.ClientOption
	if c.credentialsFile != "" {
		clientOptions = append(clientOptions, option.WithCredentialsFile(c.credentialsFile))
	}

	// Set the run client with the correct client options
	c.run, err = run.NewExecutionsClient(ctx, clientOptions...)
	if err != nil {
		return err
	}

	c.artifactRegistry, err = artifactregistry.NewClient(ctx, clientOptions...)
	if err != nil {
		return err
	}

	return nil
}

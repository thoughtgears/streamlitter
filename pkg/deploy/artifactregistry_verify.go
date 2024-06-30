package deploy

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"
)

func (c *Client) VerifyRepository(ctx context.Context, registryName string) error {
	request := &artifactregistrypb.GetRepositoryRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/repositories/%s", c.project, c.region, registryName),
	}

	_, err := c.artifactRegistry.GetRepository(ctx, request)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return errors.New("repository not found, please specify a valid repository")
		}
		return fmt.Errorf("GetRepository(): %w", err)
	}

	return nil
}

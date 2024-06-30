package deploy

import (
	"fmt"

	"github.com/thoughtgears/streamlitter/config"

	"cloud.google.com/go/run/apiv2/runpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) DeployApplication(project, region string, app config.AppConfig) (string, error) {
	exists, err := c.serviceExists(project, region, app.ServiceName)
	if err != nil {
		return "", fmt.Errorf("serviceExists(): %w", err)
	}

	switch exists {
	case true:
		fmt.Println("Service exists, updating service")
	case false:
		fmt.Println("Service does not exist, creating service")
		uri, err := c.createService(app)
		if err != nil {
			return "", fmt.Errorf("createService(): %w", err)
		}
		return uri, nil
	}

	return "", nil
}

func (c *Client) serviceExists(project, region, serviceName string) (bool, error) {
	req := &runpb.GetServiceRequest{Name: fmt.Sprintf("projects/%s/locations/%s/services/%s", project, region, serviceName)}
	_, err := c.run.GetService(c.ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return false, nil
		}
		return false, fmt.Errorf("GetService(): %w", err)
	}

	return true, nil
}

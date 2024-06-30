package deploy

import (
	"fmt"

	"github.com/thoughtgears/streamlit-hoster/config"

	"cloud.google.com/go/run/apiv2/runpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) DeployApplication(appName, version string, public bool) {
	if c.serviceExists(appName) {
		// update service
	} else {
		// create service
	}

}

func (c *Client) serviceExists(serviceName string) bool {
	req := &runpb.GetServiceRequest{Name: serviceName}
	_, err := c.run.GetService(c.ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return false
		}
	}

	return true
}

func (c *Client) createService(app config.AppConfig) (string, error) {
	var service = app.Name
	if app.Version != "" {
		service = fmt.Sprintf("%s-%s", app.Name, app.Version)
	}

	req := &runpb.CreateServiceRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", c.project, c.region),
		Service: &runpb.Service{
			Name:        service,
			Description: app.Description,
			Labels: map[string]string{
				"streamlitter": "true",
			},
			Ingress: runpb.IngressTraffic_INGRESS_TRAFFIC_ALL,
			Scaling: &runpb.ServiceScaling{
				MinInstanceCount: app.Scaling.Min,
			},
			Template: &runpb.RevisionTemplate{
				Labels:                        nil,
				Annotations:                   nil,
				Scaling:                       nil,
				VpcAccess:                     nil,
				Timeout:                       nil,
				ServiceAccount:                "",
				Containers:                    nil,
				Volumes:                       nil,
				ExecutionEnvironment:          0,
				EncryptionKey:                 "",
				MaxInstanceRequestConcurrency: 0,
				SessionAffinity:               false,
				HealthCheckDisabled:           false,
			},
		},
		ServiceId: service,
	}

	op, err := c.run.CreateService(c.ctx, req)
	if err != nil {
		return "", fmt.Errorf("run.CreateService(): %w", err)
	}

	resp, err := op.Wait(c.ctx)
	if err != nil {
		return "", fmt.Errorf("op.Wait(): %w", err)
	}

	fmt.Println(resp.Template)

	return resp.Uri, nil
}

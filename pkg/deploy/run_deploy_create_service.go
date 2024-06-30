package deploy

import (
	"fmt"
	"strings"

	"github.com/thoughtgears/streamlitter/config"

	"cloud.google.com/go/run/apiv2/runpb"
	"google.golang.org/api/iam/v1"
	"google.golang.org/protobuf/types/known/durationpb"
)

const defaultContainerTimeout = 300

func (c *Client) createService(app config.AppConfig) (string, error) {
	var service = app.Name
	if app.Version != "" {
		service = fmt.Sprintf("%s-%s", app.Name, app.Version)
	}

	serviceAccount, err := c.createServiceAccount(app.Name)
	if err != nil {
		return "", fmt.Errorf("createServiceAccount(): %w", err)
	}

	var containerEnv []*runpb.EnvVar
	if app.Env != nil {
		for _, env := range app.Env {
			containerEnv = append(containerEnv, &runpb.EnvVar{
				Name:   env.Name,
				Values: &runpb.EnvVar_Value{Value: env.Value},
			})
		}
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
			Template: &runpb.RevisionTemplate{
				Scaling: &runpb.RevisionScaling{
					MinInstanceCount: app.Scaling.Min,
					MaxInstanceCount: app.Scaling.Max,
				},
				Timeout:        &durationpb.Duration{Seconds: defaultContainerTimeout},
				ServiceAccount: serviceAccount,
				Containers: []*runpb.Container{
					{
						Image: app.ImageURL,
						Resources: &runpb.ResourceRequirements{
							Limits: map[string]string{
								"cpu":    app.Limits.CPU,
								"memory": app.Limits.Memory},
							CpuIdle:         true,
							StartupCpuBoost: false,
						},
						Env: containerEnv,
					},
				},
				ExecutionEnvironment:          runpb.ExecutionEnvironment_EXECUTION_ENVIRONMENT_GEN1,
				MaxInstanceRequestConcurrency: app.Scaling.Concurrency,
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

func (c *Client) createServiceAccount(appName string) (string, error) {
	req := &iam.CreateServiceAccountRequest{
		AccountId: fmt.Sprintf("run-%s", appName),
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: fmt.Sprintf("[RUN] %s", strings.ToTitle(appName)),
		},
	}

	account, err := c.iam.Projects.ServiceAccounts.Create(fmt.Sprintf("projects/%s", c.project), req).Do()
	if err != nil {
		return "", fmt.Errorf("serviceAccounts.Create(): %w", err)
	}

	return account.Email, nil
}

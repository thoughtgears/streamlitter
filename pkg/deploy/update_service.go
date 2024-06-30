package deploy

import (
	"fmt"

	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/thoughtgears/streamlitter/config"
)

func (c *Client) updateService(app config.AppConfig) (string, error) {
	getRequest := &runpb.GetServiceRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/services/%s", c.project, c.region, app.ServiceName),
	}
	service, err := c.run.GetService(c.ctx, getRequest)
	if err != nil {
		return "", fmt.Errorf("GetService(): %w", err)
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

	service.Template.Containers[0].Image = app.ImageURL
	service.Template.Containers[0].Resources.Limits["cpu"] = app.Limits.CPU
	service.Template.Containers[0].Resources.Limits["memory"] = app.Limits.Memory
	service.Template.Scaling.MinInstanceCount = app.Scaling.Min
	service.Template.Scaling.MaxInstanceCount = app.Scaling.Max
	service.Template.MaxInstanceRequestConcurrency = app.Scaling.Concurrency
	service.Template.Containers[0].Env = containerEnv

	updateRequest := &runpb.UpdateServiceRequest{
		Service: service,
	}

	ops, err := c.run.UpdateService(c.ctx, updateRequest)
	if err != nil {
		return "", fmt.Errorf("UpdateService(): %w", err)
	}

	updatedService, err := ops.Wait(c.ctx)
	if err != nil {
		return "", fmt.Errorf("ops.Wait(): %w", err)
	}

	if updatedService.TerminalCondition.State != runpb.Condition_CONDITION_SUCCEEDED {
		return "", fmt.Errorf("error updating service: %s", updatedService.TerminalCondition.Message)
	}

	return updatedService.Template.Revision, nil
}

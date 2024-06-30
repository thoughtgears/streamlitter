package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	defaultScalingMax  = 1
	defaultConcurrency = 80
)

// Config represents the configuration for the deployment.
// It should contain all the information needed to deploy the apps.
type Config struct {
	Project              string      `yaml:"-"`
	Region               string      `yaml:"-"`
	ArtifactRegistryName string      `yaml:"-"`
	Apps                 []AppConfig `yaml:"apps"`
}

// AppConfig represents the configuration for a single app.
// It should contain all the information needed to deploy a single
// streamlit app to Google Cloud run.
type AppConfig struct {
	Name        string     `yaml:"name"`
	Public      bool       `yaml:"public,omitempty"`
	Image       string     `yaml:"image,omitempty"`
	ImageTag    string     `yaml:"image-tag,omitempty"`
	Version     string     `yaml:"version,omitempty"`
	Scaling     AppScaling `yaml:"scaling,omitempty"`
	Limits      AppLimits  `yaml:"limits,omitempty"`
	Env         []AppEnv   `yaml:"env,omitempty"`
	ServiceName string     `yaml:"-"`
	ImageURL    string     `yaml:"-"`
}

type AppScaling struct {
	Min         int32 `yaml:"min,omitempty"`
	Max         int32 `yaml:"max,omitempty"`
	Concurrency int32 `yaml:"concurrency,omitempty"`
}

type AppLimits struct {
	CPU    string `yaml:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

type AppEnv struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func NewConfig(project, region, repo, filePath string) (Config, error) {
	config := Config{
		Project:              project,
		Region:               region,
		ArtifactRegistryName: repo,
	}

	if err := config.parseAppConfig(filePath); err != nil {
		return Config{}, err
	}

	return config, nil
}

// parseAppConfig reads the configuration file at the provided file path and
// parses it into the Config struct.
// If the file cannot be read or the YAML cannot be parsed, the function returns an error.
// The function returns nil if the file is successfully parsed.
func (c *Config) parseAppConfig(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("os.ReadFile(): %w", err)
	}

	if err := c.parseYamlFile(data); err != nil {
		return fmt.Errorf("parseYamlFile(): %w", err)
	}

	return nil
}

// parseYamlFile unmarshals the provided byte slice into a Config struct.
func (c *Config) parseYamlFile(data []byte) error {
	if err := yaml.Unmarshal(data, &c); err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	// Set default values for certain fields in the app configs
	for idx, app := range c.Apps {
		var imageTag = app.ImageTag
		if app.ImageTag == "" {
			c.Apps[idx].ImageTag = "latest"
			imageTag = "latest"
		}

		if app.Image == "" {
			c.Apps[idx].Image = app.Name
		}
		c.Apps[idx].ImageURL = fmt.Sprintf("%s-docker.pkg.dev/%s/%s/%s:%s", c.Region, c.Project, c.ArtifactRegistryName, app.Image, imageTag)

		if app.Scaling.Max == 0 {
			c.Apps[idx].Scaling.Max = defaultScalingMax
		}

		if app.Scaling.Concurrency == 0 {
			c.Apps[idx].Scaling.Concurrency = defaultConcurrency
		}

		if app.Version != "" {
			c.Apps[idx].ServiceName = fmt.Sprintf("%s-%s", app.Name, app.Version)
		}

		if app.Limits.CPU == "" {
			c.Apps[idx].Limits.CPU = "1000m"
		}

		if app.Limits.Memory == "" {
			c.Apps[idx].Limits.Memory = "256Mi"
		}
	}

	if err := c.validate(); err != nil {
		return fmt.Errorf("validate(): %w", err)
	}

	return nil
}

// validate checks the configuration for any missing or invalid fields.
func (c *Config) validate() error {
	if len(c.Apps) == 0 {
		return errors.New("at least one app must be defined")
	}

	appNames := make(map[string]bool)
	for _, app := range c.Apps {
		if app.Name == "" {
			return errors.New("app name is a required field")
		}
		if appNames[app.Name] {
			return errors.New("app names must be unique")
		}
		appNames[app.Name] = true
	}

	return nil
}

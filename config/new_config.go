package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration for the deployment.
// It should contain all the information needed to deploy the apps.
type Config struct {
	Project              string
	Region               string
	ArtifactRegistryName string
	Apps                 []AppConfig `yaml:"apps"`
}

// AppConfig represents the configuration for a single app.
// It should contain all the information needed to deploy a single
// streamlit app to Google Cloud run.
type AppConfig struct {
	Name         string `yaml:"name"`
	Description  string `yaml:"description,omitempty"`
	Public       bool   `yaml:"public,omitempty"`
	Image        string `yaml:"image,omitempty"`
	fullImageURL string
	Version      string     `yaml:"version,omitempty"`
	Scaling      AppScaling `yaml:"scaling,omitempty"`
}

type AppScaling struct {
	Min         int32 `yaml:"min,omitempty"`
	Max         int32 `yaml:"max,omitempty"`
	Concurrency int32 `yaml:"concurrency,omitempty"`
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
	for i, app := range c.Apps {
		if app.Image == "" {
			c.Apps[i].Image = app.Name
		}
		c.Apps[i].fullImageURL = fmt.Sprintf("%s-docker.pkg.dev/%s/%s/%s", c.Region, c.Project, c.ArtifactRegistryName, app.Image)

		if app.Scaling.Max == 0 {
			c.Apps[i].Scaling.Max = 1
		}

		if app.Scaling.Concurrency == 0 {
			c.Apps[i].Scaling.Concurrency = 80
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
		return fmt.Errorf("at least one app must be defined")
	}

	appNames := make(map[string]bool)
	for _, app := range c.Apps {
		if app.Name == "" {
			return fmt.Errorf("app name is a required field")
		}
		if appNames[app.Name] {
			return fmt.Errorf("app names must be unique")
		}
		appNames[app.Name] = true
	}

	return nil
}

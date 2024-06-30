package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var c = Config{
	Project:              "my-project",
	Region:               "us-central1",
	ArtifactRegistryName: "my-repo",
}

func TestParse_Yaml(t *testing.T) {
	data := `
apps:
  - name: "app1"
    description: "test app"
    public: false
    image: "test"
    version: "v1"
    scaling:
      min: 1
      max: 10
      concurrency: 5
  - name: "app2"
    public: true

`

	err := c.parseYamlFile([]byte(data))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(c.Apps))
	assert.Equal(t, "app1", c.Apps[0].Name)
	assert.False(t, c.Apps[0].Public)
	assert.Equal(t, "test", c.Apps[0].Image)
	assert.Equal(t, "v1", c.Apps[0].Version)
	assert.Equal(t, 1, c.Apps[0].Scaling.Min)
	assert.Equal(t, 10, c.Apps[0].Scaling.Max)
	assert.Equal(t, 5, c.Apps[0].Scaling.Concurrency)
	assert.Equal(t, c.Apps[0].fullImageURL, fmt.Sprintf("%s-docker.pkg.dev/%s/%s/%s", c.Region, c.Project, c.ArtifactRegistryName, c.Apps[0].Image))

	assert.Equal(t, "app2", c.Apps[1].Name)
	assert.True(t, c.Apps[1].Public)
}

func TestParse_ValidationNoApps(t *testing.T) {
	data := `
apps:
`
	err := c.parseYamlFile([]byte(data))
	assert.Error(t, err)
	assert.Equal(t, "validate(): at least one app must be defined", err.Error())
}

func TestParse_ValidationNoAppName(t *testing.T) {
	data := `
apps:
  - name: "app1"
    public: false
  - public: true
`
	err := c.parseYamlFile([]byte(data))
	assert.Error(t, err)
	assert.Equal(t, "validate(): app name is a required field", err.Error())
}

func TestParse_ValidationNonUniqueAppNames(t *testing.T) {
	data := `
apps:
  - name: "app1"
    public: false
  - name: "app1"
    public: true
`
	err := c.parseYamlFile([]byte(data))
	assert.Error(t, err)
	assert.Equal(t, "validate(): app names must be unique", err.Error())
}

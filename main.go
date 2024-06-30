package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/thoughtgears/streamlitter/config"
	"github.com/thoughtgears/streamlitter/pkg/deploy"
)

var (
	filePath        string
	region          string
	project         string
	credentialsFile string
	repository      string
	debug           bool
	version         bool
	Version         string
)

func init() {
	flag.StringVar(&filePath, "file-path", "apps.yaml", "Path to the configuration file")
	flag.StringVar(&region, "region", "europe-west1", "GCP Region")
	flag.StringVar(&project, "project", "", "GCP Project")
	flag.StringVar(&repository, "repository", "", "Artifact Registry Repository")
	flag.StringVar(&credentialsFile, "credentials-file", "", "Google JSON key file")
	flag.BoolVar(&debug, "debug", false, "For debug purposes only")
	flag.BoolVar(&version, "version", false, "Display app version")
	flag.Parse()

	if version {
		fmt.Printf("Version: %s\n", Version)
		os.Exit(0)
	}

	if project == "" {
		project = os.Getenv("GPC_PROJECT_ID")
	}

	if region == "" {
		region = os.Getenv("GPC_REGION")
	}

	if repository == "" {
		repository = os.Getenv("REPO_NAME")
	}
}

func main() {
	cfg, err := config.NewConfig(project, region, repository, filePath)
	if err != nil {
		log.Fatal(err)
	}

	if debug {
		fmt.Println("GCP config")
		fmt.Printf("  Project: %s\n  Region: %s\n  Repo: %s\n", cfg.Project, cfg.Region, cfg.ArtifactRegistryName)
		fmt.Println("App Configs:")
		for _, app := range cfg.Apps {
			fmt.Printf("  Name: %s\n  Public App: %v\n", app.Name, app.Public)
			fmt.Printf("  Image: %s\n  Version: %s\n", app.Image, app.Version)
			fmt.Printf("  ImageURL: %s\n  ServiceName: %s\n\n", app.ImageURL, app.ServiceName)
		}
	}

	clients, err := deploy.NewClient(project, region, credentialsFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := clients.VerifyRepository(context.TODO(), repository); err != nil {
		log.Fatal(err)
	}

	for _, app := range cfg.Apps {
		uri, err := clients.DeployApplication(project, region, app)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("App deployed to: %s\n", uri)
	}
}

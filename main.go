package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/cesarokuti/releases-monitoring/gchat"
	"github.com/cesarokuti/releases-monitoring/helm"
	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"
)

var (
	release       Releases
	latestVersion string
)

type Release struct {
	Name       string `json:"name"`
	Provider   string `json:"provider"`
	Repository string `json:"repository"`
}

type Releases struct {
	Dependencies []Release `json:"dependencies"`
}

func main() {
	owner := os.Getenv("GH_OWNER")
	repo := os.Getenv("GH_REPO")
	path := os.Getenv("GH_PATH")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	_, directories, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if len(directories) > 0 {
		for _, directory := range directories {
			filePath := *directory.Path + "/Chart.yaml"
			contentYaml, _, _, err := client.Repositories.GetContents(ctx, owner, repo, filePath, nil)
			if err != nil {
				fmt.Printf("Chart.yaml not found: %v\n", err)
				continue
			}
			y, _ := contentYaml.GetContent()
			c := helm.GetChartFile(y)

			filePath = *directory.Path + "/releases.json"
			contentJson, _, _, err := client.Repositories.GetContents(ctx, owner, repo, filePath, nil)
			if err != nil {
				fmt.Printf("Releases file not found: %v trying to get Chart url\n", filePath)
			} else {
				j, _ := contentJson.GetContent()
				json.Unmarshal([]byte(j), &release)
			}

			for _, depChart := range c.Dependencies {
				if contentJson != nil {
					for _, depJson := range release.Dependencies {
						if depChart.Name == depJson.Name {
							if depJson.Provider == "artifacthub" {
								latestVersion = helm.ArtifactHub(depJson.Repository)
								if helm.VersionCompare(latestVersion, depChart.Version) != "" {
									fmt.Printf("%s %s %s\n", depChart.Name, helm.VersionCompare(latestVersion, depChart.Version), latestVersion)
									gchat.SendAlert(depChart.Name, latestVersion)
								}
							} else {
								fmt.Printf("Provider not supported: %s for package %s\n", depJson.Provider, depJson.Name)
								latestVersion = ""
							}
						}
					}
				} else if strings.HasPrefix(depChart.Repository, "https://") {
					latestVersion, err = helm.ChartVersion(depChart.Repository, depChart.Name)
					if err != nil {
						fmt.Println(err)
						continue
					}
					if helm.VersionCompare(latestVersion, depChart.Version) != "" {
						fmt.Printf("%s %s %s\n", depChart.Name, helm.VersionCompare(latestVersion, depChart.Version), latestVersion)
						gchat.SendAlert(depChart.Name, latestVersion)
					}

				} else if strings.HasPrefix(depChart.Repository, "oci://") {
					fmt.Printf("OCI repository not supported: %s\n", depChart.Repository)
					latestVersion = ""
				}

			}
		}
	}
}

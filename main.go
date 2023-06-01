package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/cesarokuti/releases-monitoring/helm"
	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"
)

type Release struct {
	Name       string `json:"name"`
	Provider   string `json:"provider"`
	Repository string `json:"repository"`
}

type Releases struct {
	Dependencies []Release `json:"dependencies"`
}

func compareFile(r Releases, c helm.Chart) (string, string, string) {
	for _, depJson := range r.Dependencies {
		if depJson.Provider == "artifacthub" {
			latestVersion := helm.ArtifactHub(depJson.Name, depJson.Repository)

			v1, _ := semver.NewVersion(latestVersion)

			for _, depChart := range c.Dependencies {
				if depChart.Name == depJson.Name {
					v2, _ := semver.NewVersion(depChart.Version)
					if v1.GreaterThan(v2) {
						return "Package have new release", depJson.Name, latestVersion
					} else {
						return "Package is up-to-date", depJson.Name, latestVersion
					}

				}
			}

		} else {
			return "provide not supported", depJson.Name, depJson.Provider
		}
	}
	return "nothing", "to", "do"
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
			chart := helm.GetChartFile(contentYaml.GetContent())

			filePath = *directory.Path + "/releases.json"
			var release Releases
			contentJson, _, _, err := client.Repositories.GetContents(ctx, owner, repo, filePath, nil)
			if err != nil {
				fmt.Printf("Releases file not found: %v\n", filePath)
				fmt.Println(helm.ChartVersion(chart))
				continue
			}
			j, _ := contentJson.GetContent()
			json.Unmarshal([]byte(j), &release)
			fmt.Println(compareFile(release, chart))

		}
	}
}

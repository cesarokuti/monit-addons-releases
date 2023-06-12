package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/cesarokuti/releases-monitoring/gchat"
	"github.com/cesarokuti/releases-monitoring/git"
	"github.com/cesarokuti/releases-monitoring/helm"
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

	client, err := git.GetGitHubClient()
	if err != nil {
		fmt.Printf("Erro: %v\n", err)
		os.Exit(1)
	}

	directories, err := git.GetAddons(client, owner, repo, path)
	if err != nil {
		fmt.Printf("Erro ao obter addons: %v\n", err)
		os.Exit(1)
	}

	for _, directory := range directories {
		filePath := *directory.Path + "/Chart.yaml"

		yaml, err := git.GetFile(client, owner, repo, filePath)
		if err != nil {
			fmt.Printf("Chart.yaml not found: %v\n", err)
			continue
		}
		y, _ := yaml.GetContent()
		c := helm.GetChartFile(y)

		filePath = *directory.Path + "/releases.json"
		j, err := git.GetFile(client, owner, repo, filePath)
		if err != nil {
			fmt.Printf("Releases file not found: %v trying to get Chart url\n", filePath)
		} else {
			j, _ := j.GetContent()
			json.Unmarshal([]byte(j), &release)
		}

		for _, depChart := range c.Dependencies {
			if j != nil {
				for _, depJson := range release.Dependencies {
					if depChart.Name == depJson.Name {
						if depJson.Provider == "artifacthub" {
							latestVersion = helm.ArtifactHub(depJson.Repository)
							if helm.VersionCompare(latestVersion, depChart.Version) {
								fmt.Printf("%s have a new release %s\n", depChart.Name, latestVersion)
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
				if helm.VersionCompare(latestVersion, depChart.Version) {
					fmt.Printf("%s have a new release %s\n", depChart.Name, latestVersion)
					gchat.SendAlert(depChart.Name, latestVersion)
				}

			} else if strings.HasPrefix(depChart.Repository, "oci://") {
				fmt.Printf("OCI repository not supported: %s\n", depChart.Repository)
				latestVersion = ""
			}

		}
	}
}

package helm

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"gopkg.in/yaml.v3"
)

type Chart struct {
	APIVersion   string            `yaml:"apiVersion"`
	Name         string            `yaml:"name"`
	Description  string            `yaml:"description"`
	Type         string            `yaml:"type"`
	Version      string            `yaml:"version"`
	Dependencies []DependencyChart `yaml:"dependencies"`
}

type DependencyChart struct {
	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
	Repository string `yaml:"repository"`
}

func GetChartFile(y string) Chart {
	var chart Chart
	yaml.Unmarshal([]byte(y), &chart)
	return chart

}

func ChartVersion(c Chart) string {
	for _, depChart := range c.Dependencies {
		if strings.HasPrefix(depChart.Repository, "https://") {
			url := depChart.Repository + "/index.yaml"

			resp, _ := http.Get(url)

			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			regexStr := depChart.Name + `-(\d+\.\d+\.\d+)`

			regex := regexp.MustCompile(regexStr)
			matches := regex.FindStringSubmatch(string(body))
			if len(matches) < 2 {
				return "Not found Helm Chart release name"
			} else {
				v1, _ := semver.NewVersion(matches[1])
				v2, _ := semver.NewVersion(depChart.Version)
				if v1.GreaterThan(v2) {
					fmt.Printf("Chart %s have new release %s the installed release is %s\n", depChart.Name, v1, v2)
				} else {
					fmt.Printf("Chart %s is up-to-date %s\n", depChart.Name, v1)
				}
			}
		} else {
			fmt.Printf("Repo not supported %s, %s\n", depChart.Name, depChart.Repository)
		}
	}

	return ""
}

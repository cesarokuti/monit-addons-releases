package helm

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

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

func ChartVersion(r string, n string) (string, error) {
	url := r + "/index.yaml"

	resp, _ := http.Get(url)

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	regexStr := n + `-(\d+\.\d+\.\d+)`

	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return "", fmt.Errorf("Not found Helm Chart release name")
	} else {
		return matches[1], nil
	}
	return "", fmt.Errorf("Repo not supported %s, %s\n", n, r)

}

func VersionCompare(l string, a string) string {
	v1, _ := semver.NewVersion(l)
	v2, _ := semver.NewVersion(a)
	if v1.GreaterThan(v2) {
		return "have new release"
	} else {
		return "is up-to-date"
	}

}

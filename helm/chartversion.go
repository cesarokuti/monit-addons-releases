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
	Description  string            `yaml:"description,omitempty"`
	Type         string            `yaml:"type,omitempty"`
	Version      string            `yaml:"version"`
	Dependencies []DependencyChart `yaml:"dependencies"`
}

type DependencyChart struct {
	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
	Repository string `yaml:"repository"`
}

func GetChartFile(y string) (Chart, error) {
	var chart Chart

	err := yaml.Unmarshal([]byte(y), &chart)
	if err != nil {
		return chart, fmt.Errorf("failed to unmarshal de YAML %v, %v", y, err)
	}
	return chart, nil

}

func ChartVersion(r string, n string) (string, error) {
	url := r + "/index.yaml"

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error to get %s", url)
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	regexStr := n + `-(\d+\.\d+\.\d+)`

	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return "", fmt.Errorf("not found Helm Chart release name")
	} else {
		return matches[1], nil
	}

}

func VersionCompare(l string, a string) bool {
	v1, _ := semver.NewVersion(l)
	v2, _ := semver.NewVersion(a)
	if v1.GreaterThan(v2) {
		return true
	} else {
		return false
	}

}

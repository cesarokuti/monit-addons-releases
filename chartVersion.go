package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
)

func chartVersion(c Chart) string {
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

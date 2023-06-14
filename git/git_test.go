package git

import (
	"os"
	"testing"
)

func TestGitHub(t *testing.T) {
	owner := os.Getenv("GH_OWNER")
	repo := os.Getenv("GH_REPO")

	client, err := GetGitHubClient()
	if err != nil {
		t.Errorf("error to make github connection: %v", err)
	}

	_, err = GetAddons(client, owner, repo, "chart/default-add-ons")
	if err != nil {
		t.Errorf("error to get addons: %v", err)
	}

	_, err = GetFile(client, owner, repo, "chart/default-add-ons/argo-cd/Chart.yaml")
	if err != nil {
		t.Errorf("Chart.yaml not found: %v\n", err)

	}
}

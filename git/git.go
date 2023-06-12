package git

import (
	"context"
	"fmt"
	"os"

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

func GetGitHubClient() (*github.Client, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// Testar a conexão
	_, _, err := client.Users.Get(ctx, "") // Use uma chamada de API simples como exemplo
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao GitHub: %v", err)
	}

	return client, nil
}

func GetAddons(client *github.Client, o, r, p string) ([]*github.RepositoryContent, error) {
	ctx := context.Background()

	_, d, _, err := client.Repositories.GetContents(ctx, o, r, p, nil)
	if err != nil {
		return nil, fmt.Errorf("error to connect to GitHub: %v", err)
	}
	if len(d) > 0 {
		return d, nil
	}
	return nil, fmt.Errorf("error addons not founded: %v", err)
}

func GetFile(client *github.Client, o, r, p string) (*github.RepositoryContent, error) {
	ctx := context.Background()

	content, _, _, err := client.Repositories.GetContents(ctx, o, r, p, nil)
	if err != nil {
		return nil, fmt.Errorf("not found Chart.yaml: %v", err)
	}
	return content, nil
}

package helm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Package struct {
	Version string `json:"version"`
}

func ArtifactHub(r string) (string, error) {

	apiURL := "https://artifacthub.io/api/v1/packages/helm/" + r

	response, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("error to reach ArtifactHub: %v", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error to get Body data: %v", err)
	}

	var pkg Package
	err = json.Unmarshal(body, &pkg)
	if err != nil {
		return "", fmt.Errorf("error to unmarshal Json: %v", err)
	}

	return pkg.Version, nil
}

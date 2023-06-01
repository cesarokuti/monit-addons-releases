package helm

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Package struct {
	Version string `json:"version"`
}

func ArtifactHub(r string) string {

	apiURL := "https://artifacthub.io/api/v1/packages/helm/" + r

	response, err := http.Get(apiURL)
	if err != nil {
		return err.Error()
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err.Error()
	}

	var pkg Package
	err = json.Unmarshal(body, &pkg)
	if err != nil {
		return err.Error()
	}

	return pkg.Version
}

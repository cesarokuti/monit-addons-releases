package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Package struct {
	Version string `json:"version"`
}

func artifactHub(depName, depRepo string) {

	apiURL := "https://artifacthub.io/api/v1/packages/helm/" + depRepo

	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("artifactHub request error: %s", err.Error())
		return
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error to read requested Body: %s", err.Error())
		return
	}

	var pkg Package
	err = json.Unmarshal(body, &pkg)
	if err != nil {
		fmt.Printf("Error to decode JSON: %s", err.Error())
		return
	}

	fmt.Printf("The latest stable version of %s: %s\n", depName, pkg.Version)
}

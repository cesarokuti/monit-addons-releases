package main

type Dependency struct {
	Name       string `json:"name"`
	Provider   string `json:"provider"`
	Repository string `json:"repository"`
}

type Dependencies struct {
	Dependencies []Dependency `json:"dependencies"`
}

func main() {

	gitHub()
	/*
		file, err := ioutil.ReadFile("./monitoring.json")
		if err != nil {
			fmt.Printf("Error when opening file: %s", err)
		}

		var releases Dependencies
		err = json.Unmarshal(file, &releases)
		if err != nil {
			fmt.Printf("Error during Unmarshal(): %s", err)
		}

		for _, dep := range releases.Dependencies {
			if dep.Provider == "artifacthub" {
				artifactHub(dep.Name, dep.Repository)

			} else {
				fmt.Println("We not support the Provider:", dep.Provider)

			}
		}
	*/
}

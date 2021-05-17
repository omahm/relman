package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

type ProjectYaml struct {
	ID           string
	Version      string
	Repositories []Repositories
}

type Repositories struct {
	Name string
	Path string
}

func main() {
	projectYaml := ProjectYaml{}

	resp, err := http.Get("https://raw.githubusercontent.com/omahm/relman/main/relman_project.yaml")
	if err != nil {
		log.Fatalln("Unable to fetch YAML file: %v", err)
	}

	yamlFile, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	err = yaml.Unmarshal(yamlFile, &projectYaml)
	if err != nil {
		log.Fatalf("Unable to parse YAML file: %v", err)
	}

	fmt.Printf("%+v\n", projectYaml)
}

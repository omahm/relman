package main

import (
	"fmt"
	"io/ioutil"
	"log"

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

	fileName := "relman_projects.yaml"

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Reading of project file failed with:  #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &projectYaml)
	if err != nil {
		log.Fatalf("Unable to parse YAML file: %v", err)
	}

	fmt.Printf("%+v\n", projectYaml)
}

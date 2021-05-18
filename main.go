package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/go-github/v35/github"
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

const repoOwner string = "omahm"

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

	githubClient := github.NewClient(nil)

	for _, repo := range projectYaml.Repositories {
		release, _, _ := githubClient.Repositories.GetLatestRelease(context.Background(), repoOwner, repo.Path)
		current_release := release.GetTagName()

		software_version := strings.Split(current_release, ".")

		current_minor_version, _ := strconv.Atoi(software_version[1])
		release_minor_version := current_minor_version + 1
		software_version[1] = strconv.Itoa(release_minor_version)
		new_release := strings.Join(software_version, ".")

		fmt.Printf("Current release: %v \n", current_release)
		fmt.Printf("New release: %v \n", new_release)

		fmt.Printf("Repo: %v @ %v (%v) \n", repo.Name, repo.Path, release.GetTagName())
	}

}

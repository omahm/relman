package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v35/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

	token := os.Getenv("GITHUB_TOKEN")

	if len(token) < 1 {
		log.Fatalf("Github token doesn't appear to be set. Please ensure a valid token is present in your .env file")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	githubClient := github.NewClient(tc)

	for _, repo := range projectYaml.Repositories {
		release, _, _ := githubClient.Repositories.GetLatestRelease(context.Background(), repoOwner, repo.Path)
		current_release_tag := release.GetTagName()

		software_version := strings.Split(current_release_tag, ".")

		current_minor_version, _ := strconv.Atoi(software_version[1])
		release_minor_version := current_minor_version + 1
		software_version[1] = strconv.Itoa(release_minor_version)
		new_release_tag := strings.Join(software_version, ".")

		branch := "main"

		new_release := github.RepositoryRelease{TagName: &new_release_tag, TargetCommitish: &branch}

		githubClient.Repositories.CreateRelease(context.Background(), repoOwner, repo.Path, &new_release)

		fmt.Printf("Current release: %v \n", current_release_tag)
		fmt.Printf("New release: %v \n", new_release_tag)
	}

}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	zip           = "master.zip"
	configFile    = ".conf.yml"
	imageRegistry = ""
	publish       = false
)

func main() {
	log.Println("Starting registrar")

	if len(os.Args) != 2 {
		log.Fatal("Not enough args")
	}

	repo := os.Args[1]

	if err := GetArchive(repo); err != nil {
		log.Fatal(err)
	}

	if repoDir := RepoName(repo); Verify(repoDir) {
		log.Println("Verified repo")
		if err := Build(repoDir); err != nil {
			log.Fatal(err)
		}

		if err := Register(repoDir); err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Finished registration")
}

func GetArchive(repo string) error {
	resp, err := http.Get(fmt.Sprintf("%v/archive/master.zip", repo))

	defer resp.Body.Close()

	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(zip, data, 0755)

	if err != nil {
		return err
	}

	err = exec.Command("unzip", "-o", zip).Run()

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func RepoName(repo string) string {
	s := strings.Split(repo, "/")
	return fmt.Sprintf("%v-master", s[len(s)-1])
}

func Verify(repoName string) bool {
	// Validate config file
	if _, err := os.Stat(path.Join(repoName, configFile)); os.IsNotExist(err) {
		return false
	}

	//Validate Dockerfile
	if _, err := os.Stat(path.Join(repoName, "Dockerfile")); os.IsNotExist(err) {
		return false
	}

	return true
}

func Build(dir string) error {
	return exec.Command("docker", "build", dir, "-t", fmt.Sprintf("%v/%v", imageRegistry, dir)).Run()
}

// TODO: TEST THIS
func Register(repoName string) error {
	if publish {
		if err := exec.Command("docker", "push", fmt.Sprintf("%v/%v", imageRegistry, repoName)).Run(); err != nil {
			return err
		}
		RegisterConfig()
	}

	return nil
}

// TODO: Fill this out once we figure out how to talk to the rest of the app server
func RegisterConfig() {}

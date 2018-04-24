package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/convox/rack/manifest1"
	"github.com/convox/rack/structs"
	yaml "gopkg.in/yaml.v2"
)

//
// Here are all the extra functions that are used in the main.go source file
//
// SaveManifest renames the old docker-compose.convox.yml file
// to a timestamped file and writes the new, modified file to disk
func SaveManifest(fileName string, data []byte) error {
	return ioutil.WriteFile(fileName, data, 0644)
}

// GenerateBuildJSONfile creates the structure needed to save the modifed docker-compose.convox.yml file
func GenerateBuildJSONFile(m *manifest1.Manifest, appName, buildID, description string) ([]byte, error) {
	if m == nil {
		return nil, fmt.Errorf("m cannot be nil")
	}

	build := &structs.Build{
		Id:          buildID,
		App:         appName,
		Description: description,
		Release:     buildID,
	}

	manifestYaml, err := yaml.Marshal(m)
	if err != nil {
		return []byte{}, err
	}

	build.Manifest = string(manifestYaml)

	rez, err := json.Marshal(build)
	if err != nil {
		return []byte{}, err
	}

	return rez, nil
}

// GetGitHash gets the current git hash
// TODO: Hrm, not sure what this was for ... but it is here now
func GetGitHash() string {
	result := RunCommand("git rev-parse --short HEAD", true)
	return strings.Trim(result, "\n")
}

// GetRepo() returns the ECS_REPO from an env var
func GetRepo(account *string, region *string) string {
	data := os.Getenv("AWS_ECS_REPO_PREFIX")
	if data == "" {
		if *account == "" {
			log.Fatal("Please set AWS_ACCOUNT")
		}
		if *region == "" {
			log.Fatal("Please set AWS_REGION")
		}
		return fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", *account, *region)
	} else {
		return data
	}
}

// GetRegion returns the AWS_REGION from an env var
func GetRegion() string {
	return os.Getenv("AWS_REGION")
}

// GetAccount returns the AWS_ACCOUNT from an env var
func GetAccount() string {
	return os.Getenv("AWS_ACCOUNT")
}

// GetDefaultFile returns the default file name
func GetDefaultFile() string {
	return "docker-compose.convox.yml"
}

// RunCommand is the exported version of the function below
func RunCommand(cmd string, quiet bool) string {
	if !quiet {
		log.Printf("Running: %s\n", cmd)
	}
	result, err := runCommand(cmd)
	if err != nil {
		fmt.Errorf("ERROR (RunCommand): %s\n", err)
		return ""
	}
	return result
}

// This function executes the commands as needed and returns the result or an error
func runCommand(cmd string) (string, error) {
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	// Execute the command
	out, err := exec.Command(head, parts...).Output()
	return string(out), err
}

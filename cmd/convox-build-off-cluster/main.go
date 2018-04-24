package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/convox/rack/manifest1"
	yaml "gopkg.in/yaml.v2"
)

const (
	defaultDescription = "Default description here"
)

// This is used to extract the YAML structure to save it as the new docker-compose.convox.yml file
type myType struct {
	Manifest string `yaml:"manifest,omitempty"`
}

func main() {
	var inputComposeFileNameFlag = flag.String("compose-file", GetDefaultFile(), "the path to the docker-compose file to build")
	var outputComposeFileNameFlag = flag.String("output-compose-file", GetDefaultFile(), "file to write new compose file to")
	var appNameFlag = flag.String("app", "", "the name of the service")
	var buildIDFlag = flag.String("build-id", GetGitHash(), "the build ID to identify this build")
	var descriptionFlag = flag.String("description", defaultDescription, "(optional) a description of this build")

	// adding a few more flags
	var defaultAccountFlag = flag.String("defaultAccount", GetAccount(), "Account where the image should be pushed to")
	var regionFlag = flag.String("region", GetRegion(), "This is the region of the repo")
	var gitSHAFlag = flag.String("gitsha", "", "Please specify a git SHA to tag the image with")
	var repoNameFlag = flag.String("repo", GetRepo(defaultAccountFlag, regionFlag), "The REPO to be used")
	// Overwriting the default flag.Usage() function as we need to add some extra information
	// about ENV variables that need to be set accordingly
	// Any call to the program with -h or invalid parameters will show that information
	oldFlag := flag.Usage
	flag.Usage = func() {
		oldFlag()
		extraInfo()
	}

	flag.Parse()

	// Dereferencing these pointer variables
	appName := *appNameFlag
	buildID := *buildIDFlag
	inputComposeFileName := *inputComposeFileNameFlag
	outputComposeFileName := *outputComposeFileNameFlag
	description := *descriptionFlag
	repo := *repoNameFlag
	gitSHA := *gitSHAFlag

	// Mandatory flags to have; if any one is missing, stop
	if repo == "" {
		log.Printf("Please set ENV variables 'AWS_ACCOUNT' and 'AWS_REGION'")
		log.Fatal("Stopped.")
	}

	// Mandatory: Application name
	if appName == "" {
		flag.Usage()
		log.Fatal("missing '-app <app-name>'")
	}

	// Mandatory: gitsha name
	if gitSHA == "" {
		flag.Usage()
		log.Fatal("missing '-gitsha <git hash>'")
	}

	// A few more checks
	// Load the docker-compose.convox.yml (or differently named) config file
	m, err := manifest1.LoadFile(inputComposeFileName)
	if err != nil {
		flag.Usage()
		log.Fatal(err)
	}

	//
	// At this point in time we are done checking everything,
	// let's get to work
	//

	// Creating a new Manifest and build stream
	output := manifest1.NewOutput(false)

	opt := manifest1.BuildOptions{}
	buildStream := output.Stream("local-build")

	// Cycling through all service descriptions
	if err = m.Build(".", appName, buildStream, opt); err != nil {
		switch (err).(type) {
		case *os.PathError:
			e := err.(*os.PathError)
			fmt.Printf("While trying to %s '%s': Error: %s\n", e.Op, e.Path, e.Err)
		default:
			fmt.Printf("Error: %s\n", err.Error())
		}

		log.Fatal(err)
	}

	// Using this variable to extract YAML data ... to save it
	var yamlMe myType

	for key, service := range m.Services {
		// Creating the proper image name for tagging and pushing
		imageName := fmt.Sprintf("%s/%s:%s", repo, appName, key+"_"+gitSHA)

		// Creating the proper tagCmd
		// The 'latest' tagname is from the Build process and can't be changed w/o pain
		tagCmd := fmt.Sprintf("docker --config ./ tag %s/%s:%s %s", appName, key, "latest", imageName)
		// Creating the proper pushCmd
		pushCmd := fmt.Sprintf("--config ./ push %s", imageName)

		// Manipulating the service for no Build information but Image information
		// If there is no image in the manifest, create an image and safe it
		if service.Image == "" {
			fmt.Println("No image tag detected, modifying manifest")
			service.Build = manifest1.Build{}
			service.Image = imageName
			m.Services[key] = service
		} else {
			fmt.Printf("Detected image '%s', moving on without modification\n", service.Image)
		}

		// generating the needed YAML
		tempData, err := GenerateBuildJSONFile(m, appName, buildID, description)
		if err != nil {
			log.Fatal(err)
		}

		// Unmarshaling data
		yamlErr := yaml.Unmarshal(tempData, &yamlMe)
		if yamlErr != nil {
			log.Fatal(yamlErr)
		}

		// Executing the prepared tag and push commands
		RunCommand(tagCmd, false)
		params := strings.Fields(pushCmd)
		if err := manifest1.DefaultRunner.Run(buildStream, manifest1.Docker(params...), manifest1.RunnerOptions{Verbose: true}); err != nil {
			log.Fatalf("export error: %s", err)
		}
	}
	// Save the modified docker-compose.convox.yml file
	SaveManifest(outputComposeFileName, []byte(yamlMe.Manifest))
	os.Exit(0)
}

func extraInfo() {
	x := "\n" +
		"There are certain ENV variables that can be set to provide extra information to this program:" +
		"\n" +
		"AWS_REGION : specifies the region" +
		"\n" +
		"AWS_ACCOUNT: specifies the account to be used (Ninja/Production) by its ID" +
		"\n\n"
	fmt.Fprintf(os.Stderr, x)
}

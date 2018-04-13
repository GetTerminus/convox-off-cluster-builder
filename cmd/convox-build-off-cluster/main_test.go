package main

import (
	"fmt"
	"testing"
)

func TestRunCommand(t *testing.T) {
	if RunCommand("ls -l", true) == "" {
		fmt.Printf("Stuff went wrong")
		t.Error("expected string (inverted logic now)")
	}
}

func TestGetGitHash(t *testing.T) {
	if GetGitHash() == "" {
		t.Error("expected a githash, got something else")
	}
}

func TestGetRepo(t *testing.T) {
	if GetRepo(&"1234567890", &"us-east-2") == "" {
		t.Error("expected a repo. got something different")
	}
}

func TestGetRegion(t *testing.T) {
	if GetRegion() == "" {
		t.Error("expected whatever but not nothing")
	}
}

func TestInvalidGenerateBuildJSONFileCallNoManifest(t *testing.T) {
	_, err := GenerateBuildJSONFile(nil, "appname", "buildid", "description")
	if err == nil {
		t.Errorf("GenerateBuildJSONFile() call failed")
	}
}

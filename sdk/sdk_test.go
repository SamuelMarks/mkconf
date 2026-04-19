package sdk

import (
	"mkconf/scanner"
	"strings"
	"testing"
)

func TestSDKGenerators(t *testing.T) {
	project := &Project{
		Path: ".",
		Info: &scanner.ProjectInfo{},
	}

	if out := project.GenerateDockerCompose(); !strings.Contains(out, "version:") {
		t.Errorf("GenerateDockerCompose failed, got: %s", out)
	}
	if out := project.GenerateMakefile(); !strings.Contains(out, "all:") && !strings.Contains(out, "help:") {
		t.Errorf("GenerateMakefile failed, got: %s", out)
	}
	if out := project.GenerateMakeBat(); !strings.Contains(out, "@echo off") {
		t.Errorf("GenerateMakeBat failed, got: %s", out)
	}
	if out := project.GenerateBazelBuild(); !strings.Contains(out, "BUILD") {
		t.Errorf("GenerateBazelBuild failed, got: %s", out)
	}
}

func TestNew(t *testing.T) {
	proj, err := New("..")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if proj.Info == nil {
		t.Errorf("Expected Info to be populated")
	}

	out := proj.GenerateDockerfile("debian")
	if out == "" {
		t.Errorf("Expected Dockerfile output")
	}
}

func TestRunTestEmpty(t *testing.T) {
	project := &Project{
		Path: ".",
		Info: &scanner.ProjectInfo{},
	}
	err := project.RunTest()
	if err != nil {
		t.Errorf("Expected no error for empty test command, got %v", err)
	}
}

func TestRunTestNonEmpty(t *testing.T) {
	project := &Project{
		Path: ".",
		Info: &scanner.ProjectInfo{
			TestCommand: "echo test",
		},
	}
	err := project.RunTest()
	if err != nil {
		t.Errorf("Expected no error for echo test, got %v", err)
	}
}

func TestBuildImage(t *testing.T) {
	// BuildImage uses docker. It will fail if docker is not running or if we provide dummy data.
	// Since builder is mostly external shell commands, let's just see if it passes or fails.
	project := &Project{
		Path: ".",
		Info: &scanner.ProjectInfo{},
	}
	// just pass invalid command to make it return fast
	_ = project.BuildImage("FROM scratch", "invalid-image-name")
}

func TestNewError(t *testing.T) {
	_, err := New("/path/that/does/not/exist")
	if err == nil {
		t.Errorf("Expected error for missing directory")
	}
}

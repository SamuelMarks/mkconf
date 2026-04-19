package generator

import (
	"mkconf/scanner"
	"strings"
	"testing"
)

func TestBazel(t *testing.T) {
	info := &scanner.ProjectInfo{
		Language:       "go",
		InstallCommand: "go mod download",
		BuildCommand:   "go build",
		TestCommand:    "go test",
		StartCommand:   []string{"/app/app"},
	}

	bazel := GenerateBazelBuild(info)
	if !strings.Contains(bazel, "name = \"build\"") {
		t.Errorf("Expected Bazel build output")
	}
	if !strings.Contains(bazel, "cmd = \"go mod download > $@\"") {
		t.Errorf("Expected install command in Bazel")
	}
	if !strings.Contains(bazel, "cmd = \"go build > $@\"") {
		t.Errorf("Expected build command in Bazel")
	}
	if !strings.Contains(bazel, "cmd = \"go test > $@\"") {
		t.Errorf("Expected test command in Bazel")
	}
	if !strings.Contains(bazel, "cmd = \"/app/app > $@\"") {
		t.Errorf("Expected start command in Bazel")
	}
}

func TestBazelEmptyInfo(t *testing.T) {
	info := &scanner.ProjectInfo{}

	bazel := GenerateBazelBuild(info)
	if !strings.Contains(bazel, "No install dependencies command defined") {
		t.Errorf("Expected empty install cmd fallback in Bazel")
	}
	if !strings.Contains(bazel, "No build command defined") {
		t.Errorf("Expected empty build cmd fallback in Bazel")
	}
	if !strings.Contains(bazel, "No test command defined") {
		t.Errorf("Expected empty test cmd fallback in Bazel")
	}
	if !strings.Contains(bazel, "No start command defined") {
		t.Errorf("Expected empty run cmd fallback in Bazel")
	}
	if !strings.Contains(bazel, "Please install Base manually") {
		t.Errorf("Expected 'Base' fallback for language in Bazel")
	}
}

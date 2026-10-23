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
		BuildCommand:   `go build -ldflags "-s -w"`,
		TestCommand:    "go test",
		StartCommand:   []string{"/app/app"},
	}

	bazel := GenerateBazelBuild(info)
	if !strings.Contains(bazel, `name = "build"`) {
		t.Error("Expected Bazel build output")
	}
	if !strings.Contains(bazel, `srcs = glob(["**/*"])`) {
		t.Error("Expected srcs glob in Bazel genrule")
	}
	if !strings.Contains(bazel, `cmd = "go mod download > $@"`) {
		t.Error("Expected install command in Bazel")
	}
	expectedEscaped := `cmd = "go build -ldflags ` + string(rune(92)) + `"-s -w` + string(rune(92)) + `" > $@"`
	if !strings.Contains(bazel, expectedEscaped) {
		t.Errorf("Expected escaped quotes in build command, got: %s", bazel)
	}
	if !strings.Contains(bazel, `cmd = "go test > $@"`) {
		t.Error("Expected test command in Bazel")
	}
	if !strings.Contains(bazel, `cmd = "/app/app > $@"`) {
		t.Error("Expected start command in Bazel")
	}
	if !strings.Contains(bazel, `docker run --rm app-debian > $@`) {
		t.Error("Expected non-interactive docker run in Bazel")
	}
}

func TestBazelEmptyInfo(t *testing.T) {
	info := &scanner.ProjectInfo{}

	bazel := GenerateBazelBuild(info)
	if !strings.Contains(bazel, `cmd = "echo 'No install dependencies command defined' > $@"`) {
		t.Errorf("Expected valid single-quoted empty install cmd fallback in Bazel, got: %s", bazel)
	}
	if !strings.Contains(bazel, `cmd = "echo 'No build command defined' > $@"`) {
		t.Error("Expected empty build cmd fallback in Bazel")
	}
	if !strings.Contains(bazel, `cmd = "echo 'No test command defined' > $@"`) {
		t.Error("Expected empty test cmd fallback in Bazel")
	}
	if !strings.Contains(bazel, `cmd = "echo 'No start command defined' > $@"`) {
		t.Error("Expected empty run cmd fallback in Bazel")
	}
	if !strings.Contains(bazel, "Please install Base manually") {
		t.Error("Expected 'Base' fallback for language in Bazel")
	}
}

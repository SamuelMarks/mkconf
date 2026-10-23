package generator

import (
	"mkconf/scanner"
	"strings"
	"testing"
)

func TestMakeBat(t *testing.T) {
	info := &scanner.ProjectInfo{
		Language:       "go",
		InstallCommand: "go mod download",
		BuildCommand:   "go build",
		TestCommand:    "go test",
		StartCommand:   []string{"/app/app"},
	}

	makebat := GenerateMakeBat(info)
	expectedSnippets := []string{
		"make.bat install_base  - Install Go",
		":install_deps" + newline + "go mod download",
		":build" + newline + "go build",
		":test" + newline + "go test",
		":run" + newline + "/app/app",
	}
	for _, snippet := range expectedSnippets {
		if !strings.Contains(makebat, snippet) {
			t.Errorf("Expected make.bat to contain %q, got: %s", snippet, makebat)
		}
	}
}

func TestMakeBatEmptyInfo(t *testing.T) {
	info := &scanner.ProjectInfo{}

	makebat := GenerateMakeBat(info)
	if !strings.Contains(makebat, "echo No install dependencies command defined") {
		t.Error("Expected empty install cmd fallback in make.bat")
	}
	if !strings.Contains(makebat, "echo No build command defined") {
		t.Error("Expected empty build cmd fallback in make.bat")
	}
	if !strings.Contains(makebat, "echo No test command defined") {
		t.Error("Expected empty test cmd fallback in make.bat")
	}
	if !strings.Contains(makebat, "echo No start command defined") {
		t.Error("Expected empty run cmd fallback in make.bat")
	}
	if !strings.Contains(makebat, "make.bat install_base  - Install Base") {
		t.Error("Expected 'Base' fallback for language in make.bat")
	}
}

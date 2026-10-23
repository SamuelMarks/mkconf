package generator

import (
	"mkconf/scanner"
	"strings"
	"testing"
)

func TestMakefile(t *testing.T) {
	info := &scanner.ProjectInfo{
		Language:       "go",
		InstallCommand: "go mod download",
		BuildCommand:   "go build",
		TestCommand:    "go test",
		StartCommand:   []string{"/app/app"},
	}

	makefile := GenerateMakefile(info)
	expectedSnippets := []string{
		"make install_base  - Install Go",
		"install_deps:" + newline + tab + "go mod download",
		"build:" + newline + tab + "go build",
		"test:" + newline + tab + "go test",
		"run:" + newline + tab + "/app/app",
	}
	for _, snippet := range expectedSnippets {
		if !strings.Contains(makefile, snippet) {
			t.Errorf("Expected Makefile to contain %q, got: %s", snippet, makefile)
		}
	}
}

func TestMakefileEmptyInfo(t *testing.T) {
	info := &scanner.ProjectInfo{}

	makefile := GenerateMakefile(info)
	if !strings.Contains(makefile, "No install dependencies command defined") {
		t.Error("Expected empty install cmd fallback in Makefile")
	}
	if !strings.Contains(makefile, "No build command defined") {
		t.Error("Expected empty build cmd fallback in Makefile")
	}
	if !strings.Contains(makefile, "No test command defined") {
		t.Error("Expected empty test cmd fallback in Makefile")
	}
	if !strings.Contains(makefile, "No start command defined") {
		t.Error("Expected empty run cmd fallback in Makefile")
	}
	if !strings.Contains(makefile, "make install_base  - Install Base") {
		t.Error("Expected 'Base' fallback for language in Makefile")
	}
}

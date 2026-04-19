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
	if !strings.Contains(makefile, "build:") {
		t.Errorf("Expected Makefile output")
	}
}

func TestMakefileEmptyInfo(t *testing.T) {
	info := &scanner.ProjectInfo{}

	makefile := GenerateMakefile(info)
	if !strings.Contains(makefile, "No install dependencies command defined") {
		t.Errorf("Expected empty install cmd fallback in Makefile")
	}
	if !strings.Contains(makefile, "No build command defined") {
		t.Errorf("Expected empty build cmd fallback in Makefile")
	}
	if !strings.Contains(makefile, "No test command defined") {
		t.Errorf("Expected empty test cmd fallback in Makefile")
	}
	if !strings.Contains(makefile, "No start command defined") {
		t.Errorf("Expected empty run cmd fallback in Makefile")
	}
	if !strings.Contains(makefile, "make install_base  - Install Base") {
		t.Errorf("Expected 'Base' fallback for language in Makefile")
	}
}

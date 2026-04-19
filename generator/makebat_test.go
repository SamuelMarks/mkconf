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
	if !strings.Contains(makebat, ":build") {
		t.Errorf("Expected MakeBat output")
	}
}

func TestMakeBatEmptyInfo(t *testing.T) {
	info := &scanner.ProjectInfo{}

	makebat := GenerateMakeBat(info)
	if !strings.Contains(makebat, "echo No install dependencies command defined") {
		t.Errorf("Expected empty install cmd fallback in make.bat")
	}
	if !strings.Contains(makebat, "echo No build command defined") {
		t.Errorf("Expected empty build cmd fallback in make.bat")
	}
	if !strings.Contains(makebat, "echo No test command defined") {
		t.Errorf("Expected empty test cmd fallback in make.bat")
	}
	if !strings.Contains(makebat, "echo No start command defined") {
		t.Errorf("Expected empty run cmd fallback in make.bat")
	}
	if !strings.Contains(makebat, "make.bat install_base  - Install Base") {
		t.Errorf("Expected 'Base' fallback for language in make.bat")
	}
}

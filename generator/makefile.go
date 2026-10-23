// Package generator provides functionality for generating configuration files.
package generator

import (
	"fmt"
	"mkconf/scanner"
	"strings"
)

// tab holds a single horizontal tab character for Makefile recipe indentation.
const tab = "	"

// GenerateMakefile prepares a GNU Makefile from the scanned project metadata.
//
// Parameters:
//   - info: The scanned *scanner.ProjectInfo containing language and lifecycle commands.
func GenerateMakefile(info *scanner.ProjectInfo) string {
	var sb strings.Builder
	rc := resolveProjectCommands(info, "@echo", `"`)

	lines := []string{
		".PHONY: help install_base install_deps build test run build_docker run_docker",
		"",
		"help:",
		fmt.Sprintf(tab+`@echo "  make install_base  - Install %s"`, rc.lang),
		tab + `@echo "  make install_deps  - Install dependencies"`,
		tab + `@echo "  make build         - Build the application"`,
		tab + `@echo "  make test          - Run tests locally"`,
		tab + `@echo "  make run           - Run the application"`,
		tab + `@echo "  make build_docker  - Build Docker images"`,
		tab + `@echo "  make run_docker    - Run Docker images"`,
		"",
		"install_base:",
		fmt.Sprintf(tab+`@echo "Please install %s manually"`, rc.lang),
		"",
		"install_deps:",
		tab + rc.installCmd,
		"",
		"build:",
		tab + rc.buildCmd,
		"",
		"test:",
		tab + rc.testCmd,
		"",
		"run:",
		tab + rc.runCmd,
		"",
		"build_docker:",
		tab + "docker build -t app-debian -f debian.Dockerfile .",
		tab + "docker build -t app-alpine -f alpine.Dockerfile .",
		tab + "docker build -t app-distroless -f distroless.Dockerfile .",
		"",
		"run_docker:",
		tab + "docker run --rm -it app-debian",
	}

	for _, line := range lines {
		sb.WriteString(line)
		sb.WriteString(newline)
	}

	return sb.String()
}

// Package generator provides functionality for generating configuration files.
package generator

import (
	"fmt"
	"mkconf/scanner"
	"strings"
)

// GenerateMakefile prepares a Makefile.
func GenerateMakefile(info *scanner.ProjectInfo) string {
	var sb strings.Builder

	installCmd := info.InstallCommand
	if installCmd == "" {
		installCmd = "@echo \"No install dependencies command defined\""
	}

	buildCmd := info.BuildCommand
	if buildCmd == "" {
		buildCmd = "@echo \"No build command defined\""
	}

	testCmd := info.TestCommand
	if testCmd == "" {
		testCmd = "@echo \"No test command defined\""
	}

	runCmd := strings.Join(info.StartCommand, " ")
	if runCmd == "" {
		runCmd = "@echo \"No start command defined\""
	}

	lang := info.Language
	if lang == "" {
		lang = "base"
	}
	if len(lang) > 0 {
		lang = strings.ToUpper(lang[:1]) + lang[1:]
	}

	sb.WriteString(".PHONY: help install_base install_deps build test run build_docker run_docker\n\n")
	sb.WriteString("help:\n")
	sb.WriteString(fmt.Sprintf("\t@echo \"  make install_base  - Install %s\"\n", lang))
	sb.WriteString("\t@echo \"  make install_deps  - Install dependencies\"\n")
	sb.WriteString("\t@echo \"  make build         - Build the application\"\n")
	sb.WriteString("\t@echo \"  make test          - Run tests locally\"\n")
	sb.WriteString("\t@echo \"  make run           - Run the application\"\n")
	sb.WriteString("\t@echo \"  make build_docker  - Build Docker images\"\n")
	sb.WriteString("\t@echo \"  make run_docker    - Run Docker images\"\n\n")

	sb.WriteString("install_base:\n\t@echo \"Please install " + lang + " manually\"\n\n")
	sb.WriteString("install_deps:\n\t" + installCmd + "\n\n")
	sb.WriteString("build:\n\t" + buildCmd + "\n\n")
	sb.WriteString("test:\n\t" + testCmd + "\n\n")
	sb.WriteString("run:\n\t" + runCmd + "\n\n")
	sb.WriteString("build_docker:\n")
	sb.WriteString("\tdocker build -t app-debian -f debian.Dockerfile .\n")
	sb.WriteString("\tdocker build -t app-alpine -f alpine.Dockerfile .\n")
	sb.WriteString("\tdocker build -t app-distroless -f distroless.Dockerfile .\n\n")
	sb.WriteString("run_docker:\n\tdocker run --rm -it app-debian\n")

	return sb.String()
}

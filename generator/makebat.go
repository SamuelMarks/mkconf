// Package generator provides functionality for generating configuration files.
package generator

import (
	"fmt"
	"mkconf/scanner"
	"strings"
)

// GenerateMakeBat prepares a Windows make.bat script from the scanned project metadata.
//
// Parameters:
//   - info: The scanned *scanner.ProjectInfo containing language and lifecycle commands.
func GenerateMakeBat(info *scanner.ProjectInfo) string {
	var sb strings.Builder
	rc := resolveProjectCommands(info, "echo", "")

	lines := []string{
		"@echo off",
		`IF "%1"=="" GOTO help`,
		"GOTO %1",
		"",
		":help",
		fmt.Sprintf("echo   make.bat install_base  - Install %s", rc.lang),
		"echo   make.bat install_deps  - Install dependencies",
		"echo   make.bat build         - Build the application",
		"echo   make.bat test          - Run tests locally",
		"echo   make.bat run           - Run the application",
		"echo   make.bat build_docker  - Build Docker images",
		"echo   make.bat run_docker    - Run Docker images",
		"GOTO :EOF",
		"",
		":install_base",
		"echo Please install " + rc.lang + " manually",
		"GOTO :EOF",
		"",
		":install_deps",
		rc.installCmd,
		"GOTO :EOF",
		"",
		":build",
		rc.buildCmd,
		"GOTO :EOF",
		"",
		":test",
		rc.testCmd,
		"GOTO :EOF",
		"",
		":run",
		rc.runCmd,
		"GOTO :EOF",
		"",
		":build_docker",
		"docker build -t app-debian -f debian.Dockerfile .",
		"docker build -t app-alpine -f alpine.Dockerfile .",
		"docker build -t app-distroless -f distroless.Dockerfile .",
		"GOTO :EOF",
		"",
		":run_docker",
		"docker run --rm -it app-debian",
		"GOTO :EOF",
	}

	for _, line := range lines {
		sb.WriteString(line)
		sb.WriteString(newline)
	}

	return sb.String()
}

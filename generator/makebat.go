// Package generator provides functionality for generating configuration files.
package generator

import (
	"fmt"
	"mkconf/scanner"
	"strings"
)

// GenerateMakeBat prepares a make.bat file.
func GenerateMakeBat(info *scanner.ProjectInfo) string {
	var sb strings.Builder

	installCmd := info.InstallCommand
	if installCmd == "" {
		installCmd = "echo No install dependencies command defined"
	}

	buildCmd := info.BuildCommand
	if buildCmd == "" {
		buildCmd = "echo No build command defined"
	}

	testCmd := info.TestCommand
	if testCmd == "" {
		testCmd = "echo No test command defined"
	}

	runCmd := strings.Join(info.StartCommand, " ")
	if runCmd == "" {
		runCmd = "echo No start command defined"
	}

	lang := info.Language
	if lang == "" {
		lang = "base"
	}
	if len(lang) > 0 {
		lang = strings.ToUpper(lang[:1]) + lang[1:]
	}

	sb.WriteString("@echo off\n")
	sb.WriteString("IF \"%1\"==\"\" GOTO help\n")
	sb.WriteString("GOTO %1\n\n")

	sb.WriteString(":help\n")
	sb.WriteString(fmt.Sprintf("echo   make.bat install_base  - Install %s\n", lang))
	sb.WriteString("echo   make.bat install_deps  - Install dependencies\n")
	sb.WriteString("echo   make.bat build         - Build the application\n")
	sb.WriteString("echo   make.bat test          - Run tests locally\n")
	sb.WriteString("echo   make.bat run           - Run the application\n")
	sb.WriteString("echo   make.bat build_docker  - Build Docker images\n")
	sb.WriteString("echo   make.bat run_docker    - Run Docker images\n")
	sb.WriteString("GOTO :EOF\n\n")

	sb.WriteString(":install_base\n")
	sb.WriteString("echo Please install " + lang + " manually\n")
	sb.WriteString("GOTO :EOF\n\n")

	sb.WriteString(":install_deps\n")
	sb.WriteString(installCmd + "\n")
	sb.WriteString("GOTO :EOF\n\n")

	sb.WriteString(":build\n")
	sb.WriteString(buildCmd + "\n")
	sb.WriteString("GOTO :EOF\n\n")

	sb.WriteString(":test\n")
	sb.WriteString(testCmd + "\n")
	sb.WriteString("GOTO :EOF\n\n")

	sb.WriteString(":run\n")
	sb.WriteString(runCmd + "\n")
	sb.WriteString("GOTO :EOF\n\n")

	sb.WriteString(":build_docker\n")
	sb.WriteString("docker build -t app-debian -f debian.Dockerfile .\n")
	sb.WriteString("docker build -t app-alpine -f alpine.Dockerfile .\n")
	sb.WriteString("docker build -t app-distroless -f distroless.Dockerfile .\n")
	sb.WriteString("GOTO :EOF\n\n")

	sb.WriteString(":run_docker\n")
	sb.WriteString("docker run --rm -it app-debian\n")
	sb.WriteString("GOTO :EOF\n")

	return sb.String()
}

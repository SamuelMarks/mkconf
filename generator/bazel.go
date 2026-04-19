// Package generator provides functionality for generating configuration files.
package generator

import (
	"fmt"
	"mkconf/scanner"
	"strings"
)

// GenerateBazelBuild prepares a basic Bazel BUILD file.
func GenerateBazelBuild(info *scanner.ProjectInfo) string {
	var sb strings.Builder

	installCmd := info.InstallCommand
	if installCmd == "" {
		installCmd = "echo \"No install dependencies command defined\""
	}

	buildCmd := info.BuildCommand
	if buildCmd == "" {
		buildCmd = "echo \"No build command defined\""
	}

	testCmd := info.TestCommand
	if testCmd == "" {
		testCmd = "echo \"No test command defined\""
	}

	runCmd := strings.Join(info.StartCommand, " ")
	if runCmd == "" {
		runCmd = "echo \"No start command defined\""
	}

	lang := info.Language
	if lang == "" {
		lang = "base"
	}
	if len(lang) > 0 {
		lang = strings.ToUpper(lang[:1]) + lang[1:]
	}

	sb.WriteString("# Bazel BUILD file\n\n")

	sb.WriteString("genrule(\n")
	sb.WriteString("    name = \"install_base\",\n")
	sb.WriteString("    outs = [\"install_base.out\"],\n")
	sb.WriteString(fmt.Sprintf("    cmd = \"echo 'Please install %s manually' > $@\",\n", lang))
	sb.WriteString(")\n\n")

	sb.WriteString("genrule(\n")
	sb.WriteString("    name = \"install_deps\",\n")
	sb.WriteString("    outs = [\"install_deps.out\"],\n")
	sb.WriteString(fmt.Sprintf("    cmd = \"%s > $@\",\n", installCmd))
	sb.WriteString(")\n\n")

	sb.WriteString("genrule(\n")
	sb.WriteString("    name = \"build\",\n")
	sb.WriteString("    outs = [\"build.out\"],\n")
	sb.WriteString(fmt.Sprintf("    cmd = \"%s > $@\",\n", buildCmd))
	sb.WriteString(")\n\n")

	sb.WriteString("genrule(\n")
	sb.WriteString("    name = \"test\",\n")
	sb.WriteString("    outs = [\"test.out\"],\n")
	sb.WriteString(fmt.Sprintf("    cmd = \"%s > $@\",\n", testCmd))
	sb.WriteString(")\n\n")

	sb.WriteString("genrule(\n")
	sb.WriteString("    name = \"run\",\n")
	sb.WriteString("    outs = [\"run.out\"],\n")
	sb.WriteString(fmt.Sprintf("    cmd = \"%s > $@\",\n", runCmd))
	sb.WriteString(")\n\n")

	sb.WriteString("genrule(\n")
	sb.WriteString("    name = \"build_docker\",\n")
	sb.WriteString("    outs = [\"build_docker.out\"],\n")
	sb.WriteString("    cmd = \"docker build -t app-debian -f debian.Dockerfile . && docker build -t app-alpine -f alpine.Dockerfile . && docker build -t app-distroless -f distroless.Dockerfile . > $@\",\n")
	sb.WriteString(")\n\n")

	sb.WriteString("genrule(\n")
	sb.WriteString("    name = \"run_docker\",\n")
	sb.WriteString("    outs = [\"run_docker.out\"],\n")
	sb.WriteString("    cmd = \"docker run --rm -it app-debian > $@\",\n")
	sb.WriteString(")\n")

	return sb.String()
}

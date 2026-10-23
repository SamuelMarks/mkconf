// Package generator provides functionality for generating configuration files.
package generator

import (
	"fmt"
	"mkconf/scanner"
	"strings"
)

// escapeBazelCmd escapes double quotes inside a shell command for embedding in a Bazel genrule cmd string.
//
// Parameters:
//   - cmd: The raw shell command string to escape.
func escapeBazelCmd(cmd string) string {
	return strings.ReplaceAll(cmd, `"`, string(rune(92))+`"`)
}

// GenerateBazelBuild prepares a basic Bazel BUILD file from the scanned project metadata.
//
// Parameters:
//   - info: The scanned *scanner.ProjectInfo containing language and lifecycle commands.
func GenerateBazelBuild(info *scanner.ProjectInfo) string {
	var sb strings.Builder
	rc := resolveProjectCommands(info, "echo", "'")

	sb.WriteString("# Bazel BUILD file" + newline + newline)

	writeGenrule := func(name string, out string, cmdBody string) {
		sb.WriteString("genrule(" + newline)
		sb.WriteString(fmt.Sprintf(`    name = "%s",`, name) + newline)
		sb.WriteString(`    srcs = glob(["**/*"]),` + newline)
		sb.WriteString(fmt.Sprintf(`    outs = ["%s"],`, out) + newline)
		sb.WriteString(fmt.Sprintf(`    cmd = "%s",`, cmdBody) + newline)
		sb.WriteString(")" + newline)
	}

	writeGenrule("install_base", "install_base.out", fmt.Sprintf("echo 'Please install %s manually' > $@", rc.lang))
	sb.WriteString(newline)
	writeGenrule("install_deps", "install_deps.out", fmt.Sprintf("%s > $@", escapeBazelCmd(rc.installCmd)))
	sb.WriteString(newline)
	writeGenrule("build", "build.out", fmt.Sprintf("%s > $@", escapeBazelCmd(rc.buildCmd)))
	sb.WriteString(newline)
	writeGenrule("test", "test.out", fmt.Sprintf("%s > $@", escapeBazelCmd(rc.testCmd)))
	sb.WriteString(newline)
	writeGenrule("run", "run.out", fmt.Sprintf("%s > $@", escapeBazelCmd(rc.runCmd)))
	sb.WriteString(newline)
	writeGenrule("build_docker", "build_docker.out", "(docker build -t app-debian -f debian.Dockerfile . && docker build -t app-alpine -f alpine.Dockerfile . && docker build -t app-distroless -f distroless.Dockerfile .) > $@")
	sb.WriteString(newline)
	writeGenrule("run_docker", "run_docker.out", "docker run --rm app-debian > $@")

	return sb.String()
}

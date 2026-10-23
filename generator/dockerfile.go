// Package generator provides functionality for generating configuration files.
package generator

import (
	"encoding/json"
	"mkconf/config"
	"mkconf/scanner"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// resolveStartCommand returns the effective start command slice for info and langDef.
//
// Parameters:
//   - info: The scanned *scanner.ProjectInfo.
//   - langDef: The language definition from config.
func resolveStartCommand(info *scanner.ProjectInfo, langDef *config.LanguageDef) []string {
	if len(info.StartCommand) > 0 {
		return info.StartCommand
	}
	return langDef.StartCommand
}

// isSingleBinaryArtifact reports whether the compiled project outputs a single binary at /app/app.
//
// Parameters:
//   - langDef: The language definition from config.
//   - startCmd: The resolved start command slice.
func isSingleBinaryArtifact(langDef *config.LanguageDef, startCmd []string) bool {
	return langDef.IsCompiled && len(startCmd) == 1 && startCmd[0] == "/app/app"
}

// GenerateDockerfile creates a Dockerfile based on the project's detected language
// and the requested base image type (e.g., "debian", "alpine", "distroless").
//
// Parameters:
//   - info: The scanned *scanner.ProjectInfo containing language and command metadata.
//   - baseType: The target base image variant name ("debian", "alpine", or "distroless").
func GenerateDockerfile(info *scanner.ProjectInfo, baseType string) string {
	root := &parser.Node{Value: "root"}

	var safeInfo scanner.ProjectInfo
	if info != nil {
		safeInfo = *info
	}

	langDef := config.GetLanguage(safeInfo.Language)
	if langDef == nil {
		addInstruction(root, "FROM", "ubuntu:22.04")
		addInstruction(root, "WORKDIR", "/app")
		addInstruction(root, "COPY", ".", ".")
		return GenerateString(root)
	}

	startCmd := resolveStartCommand(&safeInfo, langDef)
	cmdJSON, _ := json.Marshal(startCmd)

	base := config.BaseImage(baseType)
	isDistroless := base == config.BaseDistroless
	isAlpine := base == config.BaseAlpine

	var baseImage, builderImage string

	switch base {
	case config.BaseDebian:
		baseImage = langDef.Docker.DebianImage
	case config.BaseAlpine:
		baseImage = langDef.Docker.AlpineImage
	case config.BaseDistroless:
		baseImage = langDef.Docker.DistrolessImage
		builderImage = langDef.Docker.DistrolessBuilderImage
	}

	if baseImage == "" {
		baseImage = "ubuntu:22.04"
	}

	if builderImage != "" && isDistroless {
		addInstruction(root, "FROM", builderImage, "AS", "builder")
		addInstruction(root, "WORKDIR", "/app")

		if strings.Contains(builderImage, "alpine") && len(langDef.Docker.AlpinePackages) > 0 {
			addInstruction(root, "RUN", "apk add --no-cache "+strings.Join(langDef.Docker.AlpinePackages, " "))
		}

		addInstruction(root, "COPY", ".", ".")

		if safeInfo.InstallCommand != "" {
			addInstruction(root, "RUN", safeInfo.InstallCommand)
		}
		if safeInfo.BuildCommand != "" {
			addInstruction(root, "RUN", safeInfo.BuildCommand)
		}

		addInstruction(root, "FROM", baseImage)
		addInstruction(root, "WORKDIR", "/app")

		if isSingleBinaryArtifact(langDef, startCmd) {
			addInstruction(root, "COPY", "--from=builder", "/app/app", "/app/app")
		} else {
			addInstruction(root, "COPY", "--from=builder", "/app", "/app")
		}

		if langDef.IsCompiled {
			addInstruction(root, "ENTRYPOINT", string(cmdJSON))
		} else {
			addInstruction(root, "CMD", string(cmdJSON))
		}
	} else {
		addInstruction(root, "FROM", baseImage)
		addInstruction(root, "WORKDIR", "/app")

		if isAlpine && len(langDef.Docker.AlpinePackages) > 0 {
			addInstruction(root, "RUN", "apk add --no-cache "+strings.Join(langDef.Docker.AlpinePackages, " "))
		}

		addInstruction(root, "COPY", ".", ".")

		if safeInfo.InstallCommand != "" {
			addInstruction(root, "RUN", safeInfo.InstallCommand)
		}
		if safeInfo.BuildCommand != "" {
			addInstruction(root, "RUN", safeInfo.BuildCommand)
		}

		if langDef.IsCompiled {
			addInstruction(root, "ENTRYPOINT", string(cmdJSON))
		} else {
			addInstruction(root, "CMD", string(cmdJSON))
		}
	}

	return GenerateString(root)
}

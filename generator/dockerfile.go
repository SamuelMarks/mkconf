// Package generator provides functionality for generating configuration files like Dockerfiles.
package generator

import (
	"encoding/json"
	"mkconf/config"
	"mkconf/scanner"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// GenerateDockerfile creates a Dockerfile based on the project's detected language
// and the requested base image type (e.g., "debian", "alpine", "distroless").
// Returns the stringified representation of the generated Dockerfile AST.
func GenerateDockerfile(info *scanner.ProjectInfo, baseType string) string {
	root := &parser.Node{Value: "root"}

	langDef := config.GetLanguage(info.Language)
	if langDef == nil {
		// Fallback if language not found
		addInstruction(root, "FROM", "ubuntu:22.04")
		addInstruction(root, "WORKDIR", "/app")
		addInstruction(root, "COPY", ".", ".")
		return GenerateString(root)
	}

	cmdJson, _ := json.Marshal(info.StartCommand)

	isDistroless := baseType == "distroless"
	isAlpine := baseType == "alpine"

	var baseImage, builderImage string

	switch baseType {
	case "debian":
		baseImage = langDef.Docker.DebianImage
	case "alpine":
		baseImage = langDef.Docker.AlpineImage
	case "distroless":
		baseImage = langDef.Docker.DistrolessImage
		builderImage = langDef.Docker.DistrolessBuilderImage
	}

	if baseImage == "" {
		baseImage = "ubuntu:22.04"
	}

	if builderImage != "" && isDistroless {
		addInstruction(root, "FROM", builderImage, "AS", "builder")
		addInstruction(root, "WORKDIR", "/app")
		addInstruction(root, "COPY", ".", ".")

		if info.InstallCommand != "" {
			addInstruction(root, "RUN", info.InstallCommand)
		}
		if info.BuildCommand != "" {
			addInstruction(root, "RUN", info.BuildCommand)
		}

		addInstruction(root, "FROM", baseImage)
		addInstruction(root, "WORKDIR", "/app")

		if langDef.IsCompiled {
			addInstruction(root, "COPY", "--from=builder", "/app/app", "/app/app")
			addInstruction(root, "ENTRYPOINT", string(cmdJson))
		} else {
			addInstruction(root, "COPY", "--from=builder", "/app", "/app")
			addInstruction(root, "CMD", string(cmdJson))
		}
	} else {
		addInstruction(root, "FROM", baseImage)
		addInstruction(root, "WORKDIR", "/app")

		if isAlpine && len(langDef.Docker.AlpinePackages) > 0 {
			addInstruction(root, "RUN", "apk add --no-cache "+strings.Join(langDef.Docker.AlpinePackages, " "))
		}

		addInstruction(root, "COPY", ".", ".")

		if info.InstallCommand != "" {
			addInstruction(root, "RUN", info.InstallCommand)
		}
		if info.BuildCommand != "" {
			addInstruction(root, "RUN", info.BuildCommand)
		}

		if langDef.IsCompiled {
			addInstruction(root, "ENTRYPOINT", string(cmdJson))
		} else {
			addInstruction(root, "CMD", string(cmdJson))
		}
	}

	return GenerateString(root)
}

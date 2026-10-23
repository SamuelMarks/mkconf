// Package generator provides functionality for generating configuration files.
package generator

import (
	"mkconf/scanner"

	"gopkg.in/yaml.v3"
)

// BuildConfig represents the Docker build context and Dockerfile configuration in a Compose service.
type BuildConfig struct {
	// Context specifies the path to the build context directory.
	Context string `yaml:"context"`
	// Dockerfile specifies the filename of the generated Dockerfile to build.
	Dockerfile string `yaml:"dockerfile"`
}

// DockerCompose represents a docker-compose.yml file structure.
type DockerCompose struct {
	// Version specifies the Docker Compose file format version.
	Version string `yaml:"version"`
	// Services maps service names to their Service definitions.
	Services map[string]Service `yaml:"services"`
}

// Service represents a container service within a docker-compose.yml file.
type Service struct {
	// Build specifies the build context and Dockerfile for the service.
	Build BuildConfig `yaml:"build"`
	// Ports lists port mappings in HOST:CONTAINER format.
	Ports []string `yaml:"ports,omitempty"`
	// Environment lists environment variables passed to the service container.
	Environment []string `yaml:"environment,omitempty"`
	// Command overrides the default container start command when present.
	Command []string `yaml:"command,omitempty"`
}

// GenerateDockerCompose prepares a docker-compose.yml file referencing the generated Dockerfiles.
//
// Parameters:
//   - info: The scanned *scanner.ProjectInfo containing detected start commands.
func GenerateDockerCompose(info *scanner.ProjectInfo) string {
	svc := Service{
		Build: BuildConfig{
			Context:    ".",
			Dockerfile: "debian.Dockerfile",
		},
	}
	if info != nil && len(info.StartCommand) > 0 {
		cmdCopy := make([]string, len(info.StartCommand))
		copy(cmdCopy, info.StartCommand)
		svc.Command = cmdCopy
	}

	dc := DockerCompose{
		Version: "3.8",
		Services: map[string]Service{
			"app": svc,
		},
	}

	out, _ := yaml.Marshal(dc)
	return string(out)
}

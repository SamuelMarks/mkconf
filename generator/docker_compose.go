// Package generator provides functionality for generating configuration files.
package generator

import (
	"gopkg.in/yaml.v3"
	"mkconf/scanner"
)

// DockerCompose represents a docker-compose.yml file structure.
type DockerCompose struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
}

// Service represents a service within a docker-compose.yml file.
type Service struct {
	Build       string   `yaml:"build"`
	Ports       []string `yaml:"ports,omitempty"`
	Environment []string `yaml:"environment,omitempty"`
	Command     []string `yaml:"command,omitempty"`
}

// GenerateDockerCompose prepares a basic docker-compose.yml file.
func GenerateDockerCompose(info *scanner.ProjectInfo) string {
	dc := DockerCompose{
		Version: "3.8",
		Services: map[string]Service{
			"app": {
				Build: ".",
			},
		},
	}

	bytes, _ := yaml.Marshal(dc)
	return string(bytes)
}

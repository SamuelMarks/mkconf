// Package sdk provides a programmatic interface to mkconf's scanning, generation, and container building features.
package sdk

import (
	"mkconf/builder"
	"mkconf/generator"
	"mkconf/scanner"
)

// Project wraps the scanning, configuration generation, and build functionality of mkconf for a specific repository.
type Project struct {
	// Path is the filesystem directory path of the scanned repository.
	Path string
	// Info holds the detected language and lifecycle command metadata for the repository.
	Info *scanner.ProjectInfo
}

// New creates a new Project instance by scanning the provided repository path.
//
// Parameters:
//   - path: The filesystem directory path of the repository to scan.
func New(path string) (*Project, error) {
	info, err := scanner.Scan(path)
	if err != nil {
		return nil, err
	}
	return &Project{
		Path: path,
		Info: info,
	}, nil
}

// RunTest executes the project's detected test command if one is configured.
func (p *Project) RunTest() error {
	if p != nil && p.Info != nil && p.Info.TestCommand != "" {
		return builder.RunTest(p.Path, p.Info.TestCommand)
	}
	return nil
}

// GenerateDockerfile creates a Dockerfile based on the project's detected language and requested base image.
//
// Parameters:
//   - base: The target base image variant ("debian", "alpine", or "distroless").
func (p *Project) GenerateDockerfile(base string) string {
	return generator.GenerateDockerfile(p.Info, base)
}

// BuildImage builds a Docker image for the project using the given Dockerfile content and image tag.
//
// Parameters:
//   - dockerfileContent: The complete Dockerfile text to build.
//   - imageName: The Docker image tag name to construct.
func (p *Project) BuildImage(dockerfileContent string, imageName string) error {
	return builder.BuildImage(p.Path, dockerfileContent, imageName)
}

// GenerateDockerCompose prepares a docker-compose.yml configuration string for the project.
func (p *Project) GenerateDockerCompose() string {
	return generator.GenerateDockerCompose(p.Info)
}

// GenerateMakefile prepares a GNU Makefile configuration string for the project.
func (p *Project) GenerateMakefile() string {
	return generator.GenerateMakefile(p.Info)
}

// GenerateMakeBat prepares a Windows make.bat batch script string for the project.
func (p *Project) GenerateMakeBat() string {
	return generator.GenerateMakeBat(p.Info)
}

// GenerateBazelBuild prepares a Bazel BUILD configuration string for the project.
func (p *Project) GenerateBazelBuild() string {
	return generator.GenerateBazelBuild(p.Info)
}

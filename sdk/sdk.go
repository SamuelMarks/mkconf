// Package sdk provides a programmatic interface to mkconf's scanning and generation features.
package sdk

import (
	"mkconf/builder"
	"mkconf/generator"
	"mkconf/scanner"
)

// Project wraps the functionality of mkconf for a specific repository.
type Project struct {
	Path string
	Info *scanner.ProjectInfo
}

// New creates a new Project instance by scanning the provided path.
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

// RunTest executes the test command if one was detected.
func (p *Project) RunTest() error {
	if p.Info.TestCommand != "" {
		return builder.RunTest(p.Path, p.Info.TestCommand)
	}
	return nil
}

// GenerateDockerfile creates a Dockerfile based on the project's detected language.
func (p *Project) GenerateDockerfile(base string) string {
	return generator.GenerateDockerfile(p.Info, base)
}

// BuildImage builds a Docker image for the project using the given Dockerfile content.
func (p *Project) BuildImage(dockerfileContent string, imageName string) error {
	return builder.BuildImage(p.Path, dockerfileContent, imageName)
}

// GenerateDockerCompose prepares a docker-compose.yml file.
func (p *Project) GenerateDockerCompose() string {
	return generator.GenerateDockerCompose(p.Info)
}

// GenerateMakefile prepares a Makefile.
func (p *Project) GenerateMakefile() string {
	return generator.GenerateMakefile(p.Info)
}

// GenerateMakeBat prepares a make.bat.
func (p *Project) GenerateMakeBat() string {
	return generator.GenerateMakeBat(p.Info)
}

// GenerateBazelBuild prepares a Bazel BUILD file.
func (p *Project) GenerateBazelBuild() string {
	return generator.GenerateBazelBuild(p.Info)
}

// Package scanner provides functionality to analyze a repository directory
// and detect the programming language, build commands, and start commands.
package scanner

import (
	"mkconf/config"
	"os"
	"path/filepath"
	"strings"
)

// ProjectInfo contains the detected metadata about a project.
// It includes the language, and various commands needed to install dependencies,
// build, test, and run the application.
type ProjectInfo struct {
	// Language detected in the project
	Language string
	// Dependencies lists project dependencies (unused currently)
	Dependencies []string
	// BuildCommand is the command used to compile the project
	BuildCommand string
	// TestCommand is the command used to run tests
	TestCommand string
	// RunCommand is an alternative run command (unused currently)
	RunCommand string
	// InstallCommand is the command to install dependencies
	InstallCommand string
	// StartCommand is the command array used in ENTRYPOINT or CMD
	StartCommand []string
}

// Scan analyzes the given directory and returns a ProjectInfo struct.
// It checks for known manifest files and falls back to scanning file extensions.
func Scan(dir string) (*ProjectInfo, error) {
	info := &ProjectInfo{}

	// First pass: primary heuristics
	for _, langDef := range config.Languages {
		for _, heuristic := range langDef.Heuristics.Primary {
			if fileExists(filepath.Join(dir, heuristic)) {
				info.Language = langDef.Name
				info.BuildCommand = langDef.BuildCommand
				info.TestCommand = langDef.TestCommand
				info.StartCommand = langDef.StartCommand
				info.InstallCommand = langDef.InstallCommand
				return info, nil
			}
		}
	}

	// Second pass: fallback heuristics
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			filename := filepath.Base(path)
			for _, langDef := range config.Languages {
				if config.MatchesExtension(filename, langDef.Heuristics.Fallback) {
					info.Language = langDef.Name

					// Use fallback commands if they exist, otherwise use primary
					if langDef.FallbackBuildCommand != "" {
						info.BuildCommand = langDef.FallbackBuildCommand
					} else {
						info.BuildCommand = langDef.BuildCommand
					}

					if langDef.FallbackTestCommand != "" {
						info.TestCommand = langDef.FallbackTestCommand
					} else {
						info.TestCommand = langDef.TestCommand
					}

					info.InstallCommand = langDef.InstallCommand

					if len(langDef.FallbackStartCommand) > 0 {
						// Replace {{FILE}} template with the actual matched filename
						info.StartCommand = make([]string, len(langDef.FallbackStartCommand))
						for i, cmdPart := range langDef.FallbackStartCommand {
							info.StartCommand[i] = strings.ReplaceAll(cmdPart, "{{FILE}}", filename)
						}
					} else {
						info.StartCommand = langDef.StartCommand
					}

					return filepath.SkipDir
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if info.Language == "" {
		return nil, os.ErrNotExist
	}

	return info, nil
}

// fileExists checks if a file exists and is not a directory.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

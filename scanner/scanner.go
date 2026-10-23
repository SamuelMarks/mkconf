// Package scanner provides functionality to analyze a repository directory
// and detect the programming language, build commands, and start commands.
package scanner

import (
	"io/fs"
	"mkconf/config"
	"os"
	"path/filepath"
	"strings"
)

// ProjectInfo contains the detected metadata and lifecycle commands about a project.
type ProjectInfo struct {
	// Language is the canonical name of the programming language detected in the project.
	Language string
	// BuildCommand is the shell command used to compile the project.
	BuildCommand string
	// TestCommand is the shell command used to run the project's test suite.
	TestCommand string
	// InstallCommand is the shell command to install project dependencies.
	InstallCommand string
	// StartCommand is the command and arguments slice used in ENTRYPOINT or CMD.
	StartCommand []string
	// DetectedFile is the base filename that triggered language detection.
	DetectedFile string
	// IsFallback reports whether detection occurred via fallback extension scanning.
	IsFallback bool
}

// isIgnoredDir reports whether a directory name should be skipped during fallback traversal.
//
// Parameters:
//   - name: The base name of the directory being inspected.
func isIgnoredDir(name string) bool {
	if name != "." && name != ".." && strings.HasPrefix(name, ".") {
		return true
	}
	switch name {
	case "vendor", "node_modules", "target", "dist", "build":
		return true
	default:
		return false
	}
}

// copyStartCommand returns a newly allocated copy of the given command slice.
//
// Parameters:
//   - cmd: The source command slice to clone.
func copyStartCommand(cmd []string) []string {
	out := make([]string, len(cmd))
	copy(out, cmd)
	return out
}

// applyPrimaryMatch constructs a ProjectInfo from a matched primary heuristic language definition.
//
// Parameters:
//   - langDef: The matched LanguageDef configuration.
//   - filename: The manifest filename that matched the primary heuristic.
func applyPrimaryMatch(langDef config.LanguageDef, filename string) *ProjectInfo {
	return &ProjectInfo{
		Language:       langDef.Name,
		BuildCommand:   langDef.BuildCommand,
		TestCommand:    langDef.TestCommand,
		StartCommand:   copyStartCommand(langDef.StartCommand),
		InstallCommand: langDef.InstallCommand,
		DetectedFile:   filename,
		IsFallback:     false,
	}
}

// applyFallbackCommands populates info from langDef using fallback overrides and {{FILE}} substitution.
//
// Parameters:
//   - info: The target ProjectInfo struct to populate.
//   - langDef: The matched LanguageDef configuration.
//   - filename: The matched source filename used for {{FILE}} substitution.
func applyFallbackCommands(info *ProjectInfo, langDef config.LanguageDef, filename string) {
	info.Language = langDef.Name
	info.DetectedFile = filename
	info.IsFallback = true
	info.InstallCommand = langDef.InstallCommand

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

	if len(langDef.FallbackStartCommand) > 0 {
		info.StartCommand = make([]string, len(langDef.FallbackStartCommand))
		for i, cmdPart := range langDef.FallbackStartCommand {
			info.StartCommand[i] = strings.ReplaceAll(cmdPart, "{{FILE}}", filename)
		}
	} else {
		info.StartCommand = copyStartCommand(langDef.StartCommand)
	}
}

// matchPrimaryEntries checks root directory entries against primary heuristics, optionally deferring generic Makefile matches.
//
// Parameters:
//   - entries: The root directory entries read via os.ReadDir.
//   - includeGenericMakefile: Whether to include generic Makefile heuristics in this pass.
func matchPrimaryEntries(entries []os.DirEntry, includeGenericMakefile bool) *ProjectInfo {
	for _, langDef := range config.Languages {
		if !includeGenericMakefile && langDef.Name == "c" {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if config.MatchesExtension(entry.Name(), langDef.Heuristics.Primary) {
				return applyPrimaryMatch(langDef, entry.Name())
			}
		}
	}
	return nil
}

// Scan analyzes the given directory and returns a populated ProjectInfo struct.
// It checks root directory entries against language-specific manifest patterns first,
// then generic Makefile manifests, and falls back to recursively scanning source file extensions.
//
// Parameters:
//   - dir: The root filesystem directory path of the repository to scan.
func Scan(dir string) (*ProjectInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, config.NewError(config.ErrDirectoryNotFound, "failed to read target directory", err)
	}

	// Pass 1a: check language-specific primary manifests before generic Makefile
	if info := matchPrimaryEntries(entries, false); info != nil {
		return info, nil
	}
	// Pass 1b: check generic Makefile primary heuristic
	if info := matchPrimaryEntries(entries, true); info != nil {
		return info, nil
	}

	info := &ProjectInfo{}

	// Second pass: fallback heuristics via filepath.WalkDir with fs.SkipAll on first match
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			if path != dir && isIgnoredDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		filename := d.Name()
		for _, langDef := range config.Languages {
			if config.MatchesExtension(filename, langDef.Heuristics.Fallback) {
				applyFallbackCommands(info, langDef, filename)
				return fs.SkipAll
			}
		}
		return nil
	})

	if err != nil {
		return nil, config.NewError(config.ErrDirectoryNotFound, "error walking project directory", err)
	}

	if info.Language == "" {
		return nil, config.NewError(config.ErrLanguageNotDetected, "no supported language detected in directory", os.ErrNotExist)
	}

	return info, nil
}

// Package config provides configuration, strong domain types, error definitions,
// and language heuristics definitions for mkconf.
package config

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrKind represents a strongly-typed category of domain error in mkconf.
type ErrKind string

const (
	// ErrDirectoryNotFound indicates the target scan directory does not exist or is not a directory.
	ErrDirectoryNotFound ErrKind = "directory_not_found"
	// ErrLanguageNotDetected indicates no programming language matched the directory contents.
	ErrLanguageNotDetected ErrKind = "language_not_detected"
	// ErrUnknownLanguage indicates a requested language name is not registered.
	ErrUnknownLanguage ErrKind = "unknown_language"
	// ErrInvalidBaseImage indicates an unsupported Docker base image variant was requested.
	ErrInvalidBaseImage ErrKind = "invalid_base_image"
	// ErrConfigParse indicates language configuration JSON failed to parse or validate.
	ErrConfigParse ErrKind = "config_parse_error"
	// ErrCommandExecution indicates an external subprocess command failed.
	ErrCommandExecution ErrKind = "command_execution_error"
	// ErrFileWrite indicates writing a generated file or directory failed.
	ErrFileWrite ErrKind = "file_write_error"
)

// Error represents a structured domain error with a typed category, message, and underlying cause.
type Error struct {
	// Kind categorizes the domain error.
	Kind ErrKind
	// Message provides human-readable context for the error.
	Message string
	// Err holds the wrapped underlying error, if any.
	Err error
}

// NewError constructs a new domain Error with the given kind, message, and underlying cause.
//
// Parameters:
//   - kind: The ErrKind category of the error.
//   - message: Contextual description of the failure.
//   - err: The underlying cause error (may be nil).
func NewError(kind ErrKind, message string, err error) *Error {
	return &Error{
		Kind:    kind,
		Message: message,
		Err:     err,
	}
}

// Error returns the formatted string representation of the domain Error.
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Kind, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Kind, e.Message)
}

// Unwrap returns the underlying wrapped error for compatibility with errors.Unwrap.
func (e *Error) Unwrap() error {
	return e.Err
}

// Is reports whether target matches this domain Error by ErrKind or underlying os.ErrNotExist.
//
// Parameters:
//   - target: The target error to compare against.
func (e *Error) Is(target error) bool {
	if t, ok := target.(*Error); ok {
		return e.Kind == t.Kind
	}
	if errors.Is(target, os.ErrNotExist) && (e.Kind == ErrDirectoryNotFound || e.Kind == ErrLanguageNotDetected) {
		return true
	}
	return false
}

// BaseImage represents a strongly-typed Docker base image variant.
type BaseImage string

const (
	// BaseDebian represents the standard Debian/Ubuntu-based Docker image variant.
	BaseDebian BaseImage = "debian"
	// BaseAlpine represents the lightweight Alpine Linux Docker image variant.
	BaseAlpine BaseImage = "alpine"
	// BaseDistroless represents the multi-stage Distroless runtime Docker image variant.
	BaseDistroless BaseImage = "distroless"
)

// IsValid reports whether b is one of the supported BaseImage variants.
func (b BaseImage) IsValid() bool {
	switch b {
	case BaseDebian, BaseAlpine, BaseDistroless:
		return true
	default:
		return false
	}
}

// AllBaseImages returns a new slice containing all supported BaseImage variants in canonical order.
func AllBaseImages() []BaseImage {
	return []BaseImage{BaseDebian, BaseAlpine, BaseDistroless}
}

// LanguageName represents a strongly-typed programming language identifier.
type LanguageName string

// OutputFormat represents a strongly-typed configuration file output format.
type OutputFormat string

const (
	// FormatDockerfile represents the Dockerfile output format.
	FormatDockerfile OutputFormat = "Dockerfile"
	// FormatDockerCompose represents the docker-compose.yml output format.
	FormatDockerCompose OutputFormat = "docker-compose.yml"
	// FormatMakefile represents the GNU Makefile output format.
	FormatMakefile OutputFormat = "Makefile"
	// FormatMakeBat represents the Windows make.bat batch script output format.
	FormatMakeBat OutputFormat = "make.bat"
	// FormatBazel represents the Bazel BUILD file output format.
	FormatBazel OutputFormat = "BUILD"
)

// IsValid reports whether f is a recognized OutputFormat constant.
func (f OutputFormat) IsValid() bool {
	switch f {
	case FormatDockerfile, FormatDockerCompose, FormatMakefile, FormatMakeBat, FormatBazel:
		return true
	default:
		return false
	}
}

// languagesJSON holds the embedded raw bytes of languages.json.
//
//go:embed languages.json
var languagesJSON []byte

// HeuristicsDef defines the file pattern heuristics used to identify a programming language.
type HeuristicsDef struct {
	// Primary lists manifest filenames or globs checked in the project root directory.
	Primary []string `json:"primary"`
	// Fallback lists source file extension globs checked recursively across the project tree.
	Fallback []string `json:"fallback"`
}

// DockerDef defines the Docker container image and package configuration for a language.
type DockerDef struct {
	// DebianImage specifies the Debian or Ubuntu base image tag.
	DebianImage string `json:"debian_image"`
	// AlpineImage specifies the Alpine Linux base image tag.
	AlpineImage string `json:"alpine_image"`
	// DistrolessBuilderImage specifies the build-stage image tag for multi-stage Distroless builds.
	DistrolessBuilderImage string `json:"distroless_builder_image"`
	// DistrolessImage specifies the minimal runtime stage image tag for Distroless builds.
	DistrolessImage string `json:"distroless_image"`
	// AlpinePackages lists additional apk packages required when building on Alpine.
	AlpinePackages []string `json:"alpine_packages"`
}

// LanguageDef defines the properties, lifecycle commands, and Docker images for a programming language.
type LanguageDef struct {
	// Name is the canonical lowercase identifier for the language.
	Name string `json:"name"`
	// Heuristics defines primary manifest and fallback extension patterns for detection.
	Heuristics HeuristicsDef `json:"heuristics"`
	// InstallCommand is the shell command used to install project dependencies.
	InstallCommand string `json:"install_command"`
	// BuildCommand is the shell command used to compile or build the project.
	BuildCommand string `json:"build_command"`
	// TestCommand is the shell command used to execute the project's test suite.
	TestCommand string `json:"test_command"`
	// StartCommand is the command and arguments array used for ENTRYPOINT or CMD.
	StartCommand []string `json:"start_command"`
	// IsCompiled indicates whether the language compiles to a standalone artifact.
	IsCompiled bool `json:"is_compiled"`
	// Docker contains container base image and package settings for the language.
	Docker DockerDef `json:"docker"`

	// FallbackBuildCommand is an optional build command used when detected via fallback heuristics.
	FallbackBuildCommand string `json:"fallback_build_command"`
	// FallbackStartCommand is an optional start command template used when detected via fallback heuristics.
	FallbackStartCommand []string `json:"fallback_start_command"`
	// FallbackTestCommand is an optional test command used when detected via fallback heuristics.
	FallbackTestCommand string `json:"fallback_test_command"`
}

// Clone returns a deep copy of the LanguageDef so callers cannot mutate shared slice state.
func (l LanguageDef) Clone() LanguageDef {
	cp := l
	cp.Heuristics.Primary = cloneStrings(l.Heuristics.Primary)
	cp.Heuristics.Fallback = cloneStrings(l.Heuristics.Fallback)
	cp.StartCommand = cloneStrings(l.StartCommand)
	cp.FallbackStartCommand = cloneStrings(l.FallbackStartCommand)
	cp.Docker.AlpinePackages = cloneStrings(l.Docker.AlpinePackages)
	return cp
}

// cloneStrings creates an independent copy of a string slice.
//
// Parameters:
//   - src: The source string slice to duplicate.
func cloneStrings(src []string) []string {
	if src == nil {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

// Registry stores an encapsulated collection of language definitions.
type Registry struct {
	// languages holds the ordered list of registered LanguageDef entries.
	languages []LanguageDef
}

// NewRegistry parses raw JSON language definitions and returns an initialized Registry.
//
// Parameters:
//   - data: Raw JSON bytes representing an array of LanguageDef objects.
func NewRegistry(data []byte) (*Registry, error) {
	var defs []LanguageDef
	if err := json.Unmarshal(data, &defs); err != nil {
		return nil, NewError(ErrConfigParse, "failed to unmarshal languages JSON", err)
	}
	return &Registry{languages: defs}, nil
}

// All returns a deep-copied slice of all LanguageDef entries in the Registry.
func (r *Registry) All() []LanguageDef {
	out := make([]LanguageDef, len(r.languages))
	for i, d := range r.languages {
		out[i] = d.Clone()
	}
	return out
}

// Get returns a deep copy of the LanguageDef matching name, or nil if not found.
//
// Parameters:
//   - name: The lowercase language name to look up.
func (r *Registry) Get(name string) *LanguageDef {
	for i := range r.languages {
		if r.languages[i].Name == name {
			cp := r.languages[i].Clone()
			return &cp
		}
	}
	return nil
}

// Languages holds the loaded definitions for all supported languages.
var Languages []LanguageDef

// LoadLanguages parses the given JSON data into the package-level Languages slice.
//
// Parameters:
//   - data: Raw JSON bytes representing an array of LanguageDef objects.
func LoadLanguages(data []byte) error {
	reg, err := NewRegistry(data)
	if err != nil {
		return err
	}
	Languages = reg.All()
	return nil
}

// initialize loads the embedded languagesJSON data and panics if the embedded asset is corrupted.
func initialize() {
	if err := LoadLanguages(languagesJSON); err != nil {
		panic("failed to parse languages.json: " + err.Error())
	}
}

// GetLanguage returns a deep-copied LanguageDef pointer for the given language name, or nil if not found.
//
// Parameters:
//   - name: The lowercase language identifier to look up.
func GetLanguage(name string) *LanguageDef {
	for i := range Languages {
		if Languages[i].Name == name {
			cp := Languages[i].Clone()
			return &cp
		}
	}
	return nil
}

// MatchesExtension checks if the given filename matches any of the provided glob or exact heuristics.
//
// Parameters:
//   - file: The base filename to test (e.g., "main.go", "Makefile", "app.fsproj").
//   - heuristics: The slice of glob patterns or filenames to match against.
func MatchesExtension(file string, heuristics []string) bool {
	fileLower := strings.ToLower(file)
	for _, h := range heuristics {
		hLower := strings.ToLower(h)
		if matched, err := filepath.Match(hLower, fileLower); err == nil && matched {
			return true
		}
	}
	return false
}

// init initializes the package-level Languages registry from the embedded languages.json asset.
func init() { initialize() }

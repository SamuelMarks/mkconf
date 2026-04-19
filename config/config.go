// Package config provides configuration and language definitions for mkconf.
package config

import (
	_ "embed"
	"encoding/json"
	"strings"
)

//go:embed languages.json
var languagesJSON []byte

// HeuristicsDef defines the heuristics to identify a language.
type HeuristicsDef struct {
	Primary  []string `json:"primary"`
	Fallback []string `json:"fallback"`
}

// DockerDef defines the Docker-related configuration for a language.
type DockerDef struct {
	DebianImage            string   `json:"debian_image"`
	AlpineImage            string   `json:"alpine_image"`
	DistrolessBuilderImage string   `json:"distroless_builder_image"`
	DistrolessImage        string   `json:"distroless_image"`
	AlpinePackages         []string `json:"alpine_packages"`
}

// LanguageDef defines the properties and commands for a programming language.
type LanguageDef struct {
	Name           string        `json:"name"`
	Heuristics     HeuristicsDef `json:"heuristics"`
	InstallCommand string        `json:"install_command"`
	BuildCommand   string        `json:"build_command"`
	TestCommand    string        `json:"test_command"`
	StartCommand   []string      `json:"start_command"`
	IsCompiled     bool          `json:"is_compiled"`
	Docker         DockerDef     `json:"docker"`

	// Fallback specific commands (optional)
	FallbackBuildCommand string   `json:"fallback_build_command"`
	FallbackStartCommand []string `json:"fallback_start_command"`
	FallbackTestCommand  string   `json:"fallback_test_command"`
}

// Languages holds the definitions for all supported languages.
var Languages []LanguageDef

// LoadLanguages parses the given JSON data into the Languages variable.
func LoadLanguages(data []byte) error {
	return json.Unmarshal(data, &Languages)
}

func initialize() {
	if err := LoadLanguages(languagesJSON); err != nil {
		panic("failed to parse languages.json: " + err.Error())
	}
}

// GetLanguage returns the LanguageDef for the given language name, or nil if not found.
func GetLanguage(name string) *LanguageDef {
	for i := range Languages {
		if Languages[i].Name == name {
			return &Languages[i]
		}
	}
	return nil
}

// MatchesExtension checks if the given file matches any of the provided heuristics.
func MatchesExtension(file string, heuristics []string) bool {
	fileLower := strings.ToLower(file)
	for _, h := range heuristics {
		hLower := strings.ToLower(h)
		if strings.HasPrefix(hLower, "*") {
			if strings.HasSuffix(fileLower, hLower[1:]) {
				return true
			}
		} else {
			if fileLower == hLower {
				return true
			}
		}
	}
	return false
}
func init() { initialize() }

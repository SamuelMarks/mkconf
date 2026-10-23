package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDomainErrors(t *testing.T) {
	cause := errors.New("disk full")
	errWithCause := NewError(ErrFileWrite, "cannot write Dockerfile", cause)
	if !strings.Contains(errWithCause.Error(), "file_write_error: cannot write Dockerfile: disk full") {
		t.Errorf("unexpected error string: %s", errWithCause.Error())
	}
	if !errors.Is(errWithCause, cause) {
		t.Error("expected errors.Is to match wrapped cause")
	}
	if !errors.Is(errWithCause, NewError(ErrFileWrite, "other", nil)) {
		t.Error("expected errors.Is to match same ErrKind")
	}
	if errors.Is(errWithCause, NewError(ErrConfigParse, "other", nil)) {
		t.Error("did not expect errors.Is to match different ErrKind")
	}

	errNoCause := NewError(ErrLanguageNotDetected, "no language matched", nil)
	if errNoCause.Error() != "language_not_detected: no language matched" {
		t.Errorf("unexpected error string without cause: %s", errNoCause.Error())
	}
	if !errors.Is(errNoCause, os.ErrNotExist) {
		t.Error("expected ErrLanguageNotDetected to match os.ErrNotExist for backward compatibility")
	}
	dirErr := NewError(ErrDirectoryNotFound, "dir missing", nil)
	if !errors.Is(dirErr, os.ErrNotExist) {
		t.Error("expected ErrDirectoryNotFound to match os.ErrNotExist")
	}
}

func TestBaseImageAndOutputFormatTypes(t *testing.T) {
	all := AllBaseImages()
	if len(all) != 3 {
		t.Fatalf("expected 3 base images, got %d", len(all))
	}
	for _, b := range all {
		if !b.IsValid() {
			t.Errorf("expected base image %q to be valid", b)
		}
	}
	if BaseImage("invalid").IsValid() {
		t.Error("expected invalid base image to return false")
	}

	validFormats := []OutputFormat{
		FormatDockerfile,
		FormatDockerCompose,
		FormatMakefile,
		FormatMakeBat,
		FormatBazel,
	}
	for _, f := range validFormats {
		if !f.IsValid() {
			t.Errorf("expected format %q to be valid", f)
		}
	}
	if OutputFormat("unknown.txt").IsValid() {
		t.Error("expected invalid output format to return false")
	}
}

func TestGetLanguageAndImmutability(t *testing.T) {
	lang := GetLanguage("go")
	if lang == nil {
		t.Fatal("Expected to find go language definition")
	}
	// Mutate returned copy and verify global definition is unaffected
	lang.Name = "mutated"
	lang.StartCommand[0] = "mutated-cmd"
	langAgain := GetLanguage("go")
	if langAgain.Name != "go" || langAgain.StartCommand[0] == "mutated-cmd" {
		t.Error("GetLanguage must return an isolated deep copy")
	}

	cLang := GetLanguage("c")
	if cLang == nil || len(cLang.Docker.AlpinePackages) == 0 {
		t.Fatal("Expected C language with AlpinePackages")
	}
	cLang.Docker.AlpinePackages[0] = "mutated-pkg"
	if GetLanguage("c").Docker.AlpinePackages[0] == "mutated-pkg" {
		t.Error("GetLanguage must deep-copy AlpinePackages")
	}

	if GetLanguage("Unknown") != nil {
		t.Error("Expected nil for Unknown language")
	}
}

func TestRegistryGet(t *testing.T) {
	reg, err := NewRegistry(languagesJSON)
	if err != nil {
		t.Fatalf("NewRegistry failed: %v", err)
	}
	if got := reg.Get("rust"); got == nil || got.Name != "rust" {
		t.Errorf("expected rust from Registry.Get, got %v", got)
	}
	if got := reg.Get("nonexistent"); got != nil {
		t.Errorf("expected nil for nonexistent language, got %v", got)
	}
}

func TestLanguagesJSONIntegrityAndReadmeSync(t *testing.T) {
	if len(Languages) != 40 {
		t.Errorf("expected exactly 40 unique languages, got %d", len(Languages))
	}

	seen := make(map[string]bool)
	for _, l := range Languages {
		if l.Name == "" {
			t.Error("found language entry with empty name")
		}
		if seen[l.Name] {
			t.Errorf("duplicate language entry in languages.json: %s", l.Name)
		}
		seen[l.Name] = true

		if len(l.StartCommand) == 0 {
			t.Errorf("language %s is missing start_command", l.Name)
		}
		if len(l.Heuristics.Primary) == 0 && len(l.Heuristics.Fallback) == 0 {
			t.Errorf("language %s has neither primary nor fallback heuristics", l.Name)
		}
		if l.Docker.DebianImage == "" || l.Docker.AlpineImage == "" || l.Docker.DistrolessBuilderImage == "" || l.Docker.DistrolessImage == "" {
			t.Errorf("language %s has incomplete docker image definitions", l.Name)
		}
	}

	readmeBytes, err := os.ReadFile(filepath.Join("..", "README.md"))
	if err != nil {
		t.Fatalf("failed to read README.md: %v", err)
	}
	readmeLower := strings.ToLower(string(readmeBytes))
	for name := range seen {
		checkName := name
		if name == "nodejs" {
			checkName = "node.js"
		}
		if !strings.Contains(readmeLower, "["+checkName+"]") {
			t.Errorf("language %q is missing from README.md Supported Languages table", name)
		}
	}
}

func TestMatchesExtension(t *testing.T) {
	if !MatchesExtension("main.go", []string{"*.go"}) {
		t.Error("Expected main.go to match *.go")
	}
	if !MatchesExtension("MyProject.fsproj", []string{"*.fsproj"}) {
		t.Error("Expected MyProject.fsproj to match *.fsproj")
	}
	if !MatchesExtension("Makefile", []string{"Makefile"}) {
		t.Error("Expected Makefile to match Makefile")
	}
	if MatchesExtension("main.go", []string{"*.py", "[invalid"}) {
		t.Error("Did not expect main.go to match *.py or invalid glob")
	}
}

func TestLoadLanguagesError(t *testing.T) {
	err := LoadLanguages([]byte("invalid json"))
	if err == nil {
		t.Error("Expected error for invalid json")
	}
}

func TestInitPanic(t *testing.T) {
	old := languagesJSON
	languagesJSON = []byte("invalid json")
	defer func() {
		languagesJSON = old
		_ = LoadLanguages(languagesJSON)
	}()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected initialize to panic on invalid embedded JSON")
		}
	}()

	initialize()
}

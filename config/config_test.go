package config

import (
	"testing"
)

func TestGetLanguage(t *testing.T) {
	lang := GetLanguage("go")
	if lang == nil {
		t.Errorf("Expected to find go language definition")
	}
	lang = GetLanguage("Unknown")
	if lang != nil {
		t.Errorf("Expected nil for Unknown language")
	}
}

func TestMatchesExtension(t *testing.T) {
	if !MatchesExtension("main.go", []string{"*.go"}) {
		t.Errorf("Expected main.go to match *.go")
	}
	if !MatchesExtension("Makefile", []string{"Makefile"}) {
		t.Errorf("Expected Makefile to match Makefile")
	}
	if MatchesExtension("main.go", []string{"*.py"}) {
		t.Errorf("Did not expect main.go to match *.py")
	}
}

func TestLoadLanguagesError(t *testing.T) {
	err := LoadLanguages([]byte("invalid json"))
	if err == nil {
		t.Errorf("Expected error for invalid json")
	}
}

func TestInitPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	// Temporarily mess up the JSON
	old := languagesJSON
	languagesJSON = []byte("invalid json")
	defer func() { languagesJSON = old }()

	// Call our initialization function
	initialize()
}

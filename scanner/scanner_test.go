package scanner

import (
	"errors"
	"mkconf/config"
	"os"
	"path/filepath"
	"testing"
)

func TestScan(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name       string
		files      []string
		expected   string
		isFallback bool
		errKind    config.ErrKind
	}{
		{"Go via go.mod", []string{"go.mod"}, "go", false, ""},
		{"Rust via Cargo.toml", []string{"Cargo.toml"}, "rust", false, ""},
		{"Python via requirements.txt", []string{"requirements.txt"}, "python", false, ""},
		{"Ruby via Gemfile", []string{"Gemfile"}, "ruby", false, ""},
		{"C++ via CMakeLists.txt", []string{"CMakeLists.txt"}, "c++", false, ""},
		{"C via Makefile", []string{"Makefile"}, "c", false, ""},
		{"F# via wildcard *.fsproj", []string{"App.fsproj"}, "f#", false, ""},
		{"Nim via wildcard *.nimble", []string{"pkg.nimble"}, "nim", false, ""},
		{"Go fallback", []string{"main.go"}, "go", true, ""},
		{"Rust fallback", []string{"main.rs"}, "rust", true, ""},
		{"Python fallback", []string{"script.py"}, "python", true, ""},
		{"Ruby fallback", []string{"script.rb"}, "ruby", true, ""},
		{"C fallback", []string{"main.c"}, "c", true, ""},
		{"C++ fallback", []string{"main.cpp"}, "c++", true, ""},
		{"C++ fallback cc", []string{"main.cc"}, "c++", true, ""},
		{"No language", []string{"README.md"}, "", false, config.ErrLanguageNotDetected},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := filepath.Join(tmpDir, tt.name)
			_ = os.MkdirAll(dir, 0755)
			_ = os.MkdirAll(filepath.Join(dir, "subdir"), 0755)

			for _, f := range tt.files {
				err := os.WriteFile(filepath.Join(dir, f), []byte(""), 0644)
				if err != nil {
					t.Fatalf("write error: %v", err)
				}
			}

			info, err := Scan(dir)
			if tt.errKind != "" {
				if err == nil {
					t.Fatalf("expected error %s, got nil", tt.errKind)
				}
				if !errors.Is(err, config.NewError(tt.errKind, "", nil)) {
					t.Fatalf("expected error kind %s, got %v", tt.errKind, err)
				}
				if !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("expected error to match os.ErrNotExist, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if info.Language != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, info.Language)
			}
			if info.IsFallback != tt.isFallback {
				t.Errorf("expected IsFallback=%v, got %v", tt.isFallback, info.IsFallback)
			}
		})
	}

	t.Run("Path Does Not Exist", func(t *testing.T) {
		_, err := Scan(filepath.Join(tmpDir, "nonexistent_dir"))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, config.NewError(config.ErrDirectoryNotFound, "", nil)) {
			t.Errorf("expected ErrDirectoryNotFound, got %v", err)
		}
	})
}

func TestScanSkipsIgnoredDirsAndStopsOnFirstSubdirMatch(t *testing.T) {
	tmpDir := t.TempDir()

	// Files inside .git and node_modules must be ignored;
	// name valid dirs so WalkDir visits .git and node_modules first, then y_sub, and skips z_sub.
	gitDir := filepath.Join(tmpDir, ".git")
	nodeModulesDir := filepath.Join(tmpDir, "node_modules")
	subY := filepath.Join(tmpDir, "y_sub")
	subZ := filepath.Join(tmpDir, "z_sub")

	_ = os.MkdirAll(gitDir, 0755)
	_ = os.MkdirAll(nodeModulesDir, 0755)
	_ = os.MkdirAll(subY, 0755)
	_ = os.MkdirAll(subZ, 0755)

	_ = os.WriteFile(filepath.Join(gitDir, "hook.sh"), []byte(""), 0644)
	_ = os.WriteFile(filepath.Join(nodeModulesDir, "index.js"), []byte(""), 0644)
	_ = os.WriteFile(filepath.Join(subY, "first.py"), []byte(""), 0644)
	_ = os.WriteFile(filepath.Join(subZ, "second.rs"), []byte(""), 0644)

	info, err := Scan(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Language != "python" {
		t.Errorf("expected python from first non-ignored subdir y_sub, got %s", info.Language)
	}
	if len(info.StartCommand) != 2 || info.StartCommand[1] != "first.py" {
		t.Errorf("expected {{FILE}} substitution to yield first.py, got %v", info.StartCommand)
	}
}

func TestScanWalkErrorInSubdirectory(t *testing.T) {
	tmpDir := t.TempDir()

	errSubDir := filepath.Join(tmpDir, "unreadable_sub")
	_ = os.Mkdir(errSubDir, 0000)
	t.Cleanup(func() {
		_ = os.Chmod(errSubDir, 0755)
	})

	_, err := Scan(tmpDir)
	if err == nil {
		t.Fatal("expected error from filepath.WalkDir on unreadable subdirectory, got nil")
	}
}

package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScan(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		files    []string
		expected string
		err      error
	}{
		{"Go via go.mod", []string{"go.mod"}, "go", nil},
		{"Rust via Cargo.toml", []string{"Cargo.toml"}, "rust", nil},
		{"Python via requirements.txt", []string{"requirements.txt"}, "python", nil},
		{"Ruby via Gemfile", []string{"Gemfile"}, "ruby", nil},
		{"C++ via CMakeLists.txt", []string{"CMakeLists.txt"}, "c++", nil},
		{"C via Makefile", []string{"Makefile"}, "c", nil},
		{"Go fallback", []string{"main.go"}, "go", nil},
		{"Rust fallback", []string{"main.rs"}, "rust", nil},
		{"Python fallback", []string{"script.py"}, "python", nil},
		{"Ruby fallback", []string{"script.rb"}, "ruby", nil},
		{"C fallback", []string{"main.c"}, "c", nil},
		{"C++ fallback", []string{"main.cpp"}, "c++", nil},
		{"C++ fallback cc", []string{"main.cc"}, "c++", nil},
		{"No language", []string{"README.md"}, "", os.ErrNotExist},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := filepath.Join(tmpDir, tt.name)
			_ = os.MkdirAll(dir, 0755)

			for _, f := range tt.files {
				err := os.WriteFile(filepath.Join(dir, f), []byte(""), 0644)
				if err != nil {
					t.Fatalf("write error: %v", err)
				}
			}

			info, err := Scan(dir)
			if tt.err != nil {
				if err == nil || err.Error() != tt.err.Error() {
					t.Fatalf("expected error %v, got %v", tt.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if info.Language != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, info.Language)
			}
		})
	}

	t.Run("Path Does Not Exist", func(t *testing.T) {
		_, err := Scan(filepath.Join(tmpDir, "nonexistent_dir"))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()

	fpath := filepath.Join(tmpDir, "file.txt")
	_ = os.WriteFile(fpath, []byte("test"), 0644)

	if !fileExists(fpath) {
		t.Error("expected true for existing file")
	}

	if fileExists(filepath.Join(tmpDir, "nonexistent.txt")) {
		t.Error("expected false for nonexistent file")
	}

	if fileExists(tmpDir) {
		t.Error("expected false for directory")
	}
}

func TestScanWalkError(t *testing.T) {
	tmpDir := t.TempDir()

	// Create an unreadable directory to trigger filepath.Walk error
	errDir := filepath.Join(tmpDir, "unreadable")
	_ = os.Mkdir(errDir, 0000)
	defer func() { _ = os.Chmod(errDir, 0755) }() // clean up

	_, err := Scan(errDir)
	if err == nil {
		t.Fatal("expected error from filepath.Walk, got nil")
	}
}

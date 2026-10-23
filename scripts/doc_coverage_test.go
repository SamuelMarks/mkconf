package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunCoverageOnRepoRoot(t *testing.T) {
	cov, err := run("..")
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if cov != 100.0 {
		t.Errorf("Expected exact 100.0%% doc coverage across repository, got %.2f%%", cov)
	}
}

func TestRunCoverageUndocumentedConstructs(t *testing.T) {
	tmpDir := t.TempDir()

	hiddenDir := filepath.Join(tmpDir, ".hidden")
	vendorDir := filepath.Join(tmpDir, "vendor")
	_ = os.MkdirAll(hiddenDir, 0755)
	_ = os.MkdirAll(vendorDir, 0755)
	ignoredSrc := `package ignored
func Undoc() {}
`
	_ = os.WriteFile(filepath.Join(hiddenDir, "ignored.go"), []byte(ignoredSrc), 0644)
	_ = os.WriteFile(filepath.Join(vendorDir, "ignored.go"), []byte(ignoredSrc), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("non-go file"), 0644)

	src := `package sample

type Embedded struct{}

type UndocStruct struct {
	Embedded
	FieldA int
}

var UndocVar = 1

// DocFunc has doc for documentedArg only.
func DocFunc(_ int, documentedArg string, missingParam int) {}

func UndocFunc() {}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "sample.go"), []byte(src), 0644); err != nil {
		t.Fatalf("write error: %v", err)
	}

	cov, err := run(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cov != 20.0 {
		t.Errorf("Expected 20.0%% coverage for partially documented sample, got %.1f%%", cov)
	}
}

func TestRunCoverageErrorsAndMain(t *testing.T) {
	origRoot := rootDir
	origExit := osExit
	t.Cleanup(func() {
		rootDir = origRoot
		osExit = origExit
	})

	badDir := t.TempDir()
	badFile := filepath.Join(badDir, "bad.go")
	badSrc := `package bad
func main() {
`
	if err := os.WriteFile(badFile, []byte(badSrc), 0644); err != nil {
		t.Fatalf("write error: %v", err)
	}
	if _, err := run(badDir); err == nil {
		t.Error("Expected parse error on invalid Go syntax")
	}

	emptyDir := t.TempDir()
	cov, err := run(emptyDir)
	if err != nil {
		t.Errorf("Expected no error for empty dir, got %v", err)
	}
	if cov != 100.0 {
		t.Errorf("Expected 100.0 coverage for empty dir, got %f", cov)
	}

	if _, err := run(filepath.Join(emptyDir, "does_not_exist")); err == nil {
		t.Error("Expected walk error for nonexistent path")
	}

	rootDir = emptyDir
	osExitCalled := false
	osExit = func(_ int) { osExitCalled = true }
	main()
	if osExitCalled {
		t.Error("main exited on success")
	}

	rootDir = filepath.Join(emptyDir, "does_not_exist")
	osExitCalled = false
	main()
	if !osExitCalled {
		t.Error("main did not exit on failure")
	}
}

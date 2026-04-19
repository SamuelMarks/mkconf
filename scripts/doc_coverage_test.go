package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunCoverage(t *testing.T) {
	cov, err := run("..")
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if cov < 0 || cov > 100 {
		t.Errorf("Invalid coverage: %f", cov)
	}

	err = os.MkdirAll("testdata_not_ignored", 0755)
	if err != nil {
		t.Fatalf("mkdir error: %v", err)
	}
	badFile2 := filepath.Join("testdata_not_ignored", "bad.go")
	err = os.WriteFile(badFile2, []byte("package bad\nfunc main() {\n"), 0644)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	_, err = run("testdata_not_ignored")
	if err == nil {
		t.Errorf("Expected parse error")
	}
	os.RemoveAll("testdata_not_ignored")

	err = os.MkdirAll("testdata_empty", 0755)
	if err != nil {
		t.Fatalf("mkdir error: %v", err)
	}
	if err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}
	defer os.RemoveAll("testdata_empty")
	cov, err = run("testdata_empty")
	if err != nil {
		t.Errorf("Expected no error")
	}
	if cov != 100.0 {
		t.Errorf("Expected 100.0 coverage for empty dir")
	}

	_, err = run("/path/does/not/exist/surely")
	if err == nil {
		t.Errorf("Expected walk error")
	}

	// Test main success
	rootDir = "testdata_empty"
	osExitCalled := false
	osExit = func(code int) { osExitCalled = true }
	main()
	if osExitCalled {
		t.Errorf("main exited on success")
	}

	// Test main failure
	rootDir = "/path/does/not/exist/surely"
	osExitCalled = false
	main()
	if !osExitCalled {
		t.Errorf("main did not exit on failure")
	}
}

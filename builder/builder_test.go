package builder

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func fakeExecCommandError(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcessError", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS_ERR=1"}
	return cmd
}

func TestHelperProcess(_ *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}

func TestHelperProcessError(_ *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS_ERR") != "1" {
		return
	}
	os.Exit(1)
}

func TestBuildImage(t *testing.T) {
	tmpDir := t.TempDir()
	origExec := ExecCommand
	origOut := Stdout
	origErr := Stderr
	origCreateTemp := createTempFile
	t.Cleanup(func() {
		ExecCommand = origExec
		Stdout = origOut
		Stderr = origErr
		createTempFile = origCreateTemp
	})

	var buf bytes.Buffer
	Stdout = &buf
	Stderr = &buf

	// 1. Test success
	ExecCommand = fakeExecCommand
	err := BuildImage(tmpDir, "FROM ubuntu", "test-image")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "Building image test-image...") || !strings.Contains(buf.String(), "Successfully built test-image") {
		t.Errorf("expected build output messages, got %q", buf.String())
	}

	// 2. Test createTempFile failure
	err = BuildImage("/nonexistent/dir/for/dockerfile", "FROM ubuntu", "test-image")
	if err == nil || !strings.Contains(err.Error(), "failed to create temporary Dockerfile") {
		t.Errorf("expected error due to temp file creation failure, got %v", err)
	}

	// 3. Test writeTempFile failure (closed file before write)
	createTempFile = func(dir, pattern string) (*os.File, error) {
		f, err := os.CreateTemp(dir, pattern)
		if err != nil {
			return nil, err
		}
		_ = f.Close()
		return f, nil
	}
	err = BuildImage(tmpDir, "FROM ubuntu", "test-image")
	if err == nil || !strings.Contains(err.Error(), "failed to write temporary Dockerfile") {
		t.Errorf("expected error due to write failure on closed file, got %v", err)
	}
	createTempFile = origCreateTemp

	// 4. Test command failure
	ExecCommand = fakeExecCommandError
	err = BuildImage(tmpDir, "FROM ubuntu", "test-image")
	if err == nil || !strings.Contains(err.Error(), "docker build failed for test-image") {
		t.Errorf("expected error due to cmd failure, got %v", err)
	}
}

func TestRunTest(t *testing.T) {
	tmpDir := t.TempDir()
	origExec := ExecCommand
	origOut := Stdout
	origErr := Stderr
	t.Cleanup(func() {
		ExecCommand = origExec
		Stdout = origOut
		Stderr = origErr
	})

	var buf bytes.Buffer
	Stdout = &buf
	Stderr = &buf

	// 1. Test success
	ExecCommand = fakeExecCommand
	err := RunTest(tmpDir, "echo test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "Running tests: echo test...") || !strings.Contains(buf.String(), "Tests passed successfully.") {
		t.Errorf("expected test output messages, got %q", buf.String())
	}

	// 2. Test command failure
	ExecCommand = fakeExecCommandError
	err = RunTest(tmpDir, "echo test")
	if err == nil || !strings.Contains(err.Error(), "tests failed") {
		t.Errorf("expected error due to cmd failure, got %v", err)
	}
	if errors.Unwrap(err) == nil {
		t.Errorf("expected wrapped error from RunTest")
	}
}

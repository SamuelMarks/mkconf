package builder

import (
	"os"
	"os/exec"
	"testing"
)

// Helper to mock exec.Command
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

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}

func TestHelperProcessError(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS_ERR") != "1" {
		return
	}
	os.Exit(1)
}

func TestBuildImage(t *testing.T) {
	tmpDir := t.TempDir()

	// Test success
	ExecCommand = fakeExecCommand
	err := BuildImage(tmpDir, "FROM ubuntu", "test-image")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Test write failure
	err = BuildImage("/nonexistent/dir", "FROM ubuntu", "test-image")
	if err == nil {
		t.Error("expected error due to write failure")
	}

	// Test command failure
	ExecCommand = fakeExecCommandError
	err = BuildImage(tmpDir, "FROM ubuntu", "test-image")
	if err == nil {
		t.Error("expected error due to cmd failure")
	}
}

func TestRunTest(t *testing.T) {
	tmpDir := t.TempDir()

	// Test success
	ExecCommand = fakeExecCommand
	err := RunTest(tmpDir, "echo test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Test command failure
	ExecCommand = fakeExecCommandError
	err = RunTest(tmpDir, "echo test")
	if err == nil {
		t.Error("expected error due to cmd failure")
	}
}

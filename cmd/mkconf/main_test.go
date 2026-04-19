package main

import (
	"bytes"
	"mkconf/builder"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Mocks
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

func resetFlags() {
	outputDir = ""
	noTest = false
	dryRun = false
	emitDockerfile = false
	emitBazel = false
	emitMakefile = false
	buildImages = false
}

func TestMainSuccess(t *testing.T) {
	resetFlags()
	tmpDir := t.TempDir()
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(""), 0644)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	builder.ExecCommand = fakeExecCommand

	RootCmd.SetArgs([]string{tmpDir})
	b := bytes.NewBufferString("")
	RootCmd.SetOut(b)
	RootCmd.SetErr(b)

	err = RootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMainFailureConditions(t *testing.T) {
	resetFlags()
	tmpDir := t.TempDir()

	// 1. Scan Error
	RootCmd.SetArgs([]string{filepath.Join(tmpDir, "nonexistent")})
	err := RootCmd.Execute()
	if err == nil {
		t.Fatal("expected error from scan")
	}

	// 2. Build failures, Test failures
	err = os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(""), 0644)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}
	builder.ExecCommand = fakeExecCommandError
	RootCmd.SetArgs([]string{tmpDir})
	err = RootCmd.Execute()
	// RunE returns nil for these warnings/errors per logic
	if err != nil {
		t.Fatalf("expected nil from Execute on warnings, got %v", err)
	}
}

func TestRunPath(t *testing.T) {
	resetFlags()
	exited := false
	osExit = func(code int) {
		exited = true
	}
	defer func() { osExit = os.Exit }()

	RootCmd.SetArgs([]string{"a", "b"}) // Too many args should fail
	Run()
	if !exited {
		t.Error("expected osExit to be called")
	}
}

func TestMainExecution(t *testing.T) {
	resetFlags()
	exited := false
	osExit = func(code int) {
		exited = true
	}
	defer func() { osExit = os.Exit }()

	RootCmd.SetArgs([]string{"a", "b"}) // Too many args should fail
	main()
	if !exited {
		t.Error("expected main to call osExit via Run")
	}
}

func TestMainBranches(t *testing.T) {
	resetFlags()
	tmpDir := t.TempDir()
	// trigger python to hit InstallCommand
	_ = os.WriteFile(filepath.Join(tmpDir, "requirements.txt"), []byte(""), 0644)

	builder.ExecCommand = fakeExecCommand
	RootCmd.SetArgs([]string{tmpDir})
	b := bytes.NewBufferString("")
	RootCmd.SetOut(b)
	RootCmd.SetErr(b)
	_ = RootCmd.Execute()
}

func TestOutputDirectory(t *testing.T) {
	resetFlags()
	tmpDir := t.TempDir()
	outDir := filepath.Join(tmpDir, "out")
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(""), 0644)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	builder.ExecCommand = fakeExecCommand

	RootCmd.SetArgs([]string{"--output", outDir, tmpDir})
	b := bytes.NewBufferString("")
	RootCmd.SetOut(b)
	RootCmd.SetErr(b)

	err = RootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify files were created
	expectedFiles := []string{
		"debian.Dockerfile",
		"alpine.Dockerfile",
		"distroless.Dockerfile",
		"docker-compose.yml",
		"Makefile",
		"make.bat",
		"BUILD",
	}
	for _, file := range expectedFiles {
		if _, err := os.Stat(filepath.Join(outDir, file)); os.IsNotExist(err) {
			t.Errorf("expected file %s to be created in output dir", file)
		}
	}

	// Output directory creation failure test
	resetFlags()
	badOut := filepath.Join(tmpDir, "bad_out")
	_ = os.WriteFile(badOut, []byte("file, not dir"), 0644)
	RootCmd.SetArgs([]string{"--output", badOut, tmpDir})
	err = RootCmd.Execute()
	if err == nil {
		t.Fatal("expected error from MkdirAll on existing file")
	}

	// Dockerfile WriteFile failure test by setting read-only file
	resetFlags()
	readonlyFile := filepath.Join(tmpDir, "readonly_out")
	_ = os.MkdirAll(readonlyFile, 0755)

	// create the debian.Dockerfile as a directory to force WriteFile failure
	err = os.Mkdir(filepath.Join(readonlyFile, "debian.Dockerfile"), 0755)
	if err != nil {
		t.Fatalf("mkdir error: %v", err)
	}

	RootCmd.SetArgs([]string{"--output", readonlyFile, tmpDir})
	err = RootCmd.Execute()
	if err == nil {
		t.Fatal("expected error from WriteFile on dir")
	}
}

func TestFormatWriteFailure(t *testing.T) {
	resetFlags()
	tmpDir := t.TempDir()
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(""), 0644)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}
	builder.ExecCommand = fakeExecCommand

	// Cause WriteFile to fail for docker-compose.yml by making it a directory
	badFormatFile := filepath.Join(tmpDir, "docker-compose.yml")
	err = os.Mkdir(badFormatFile, 0755)
	if err != nil {
		t.Fatalf("mkdir error: %v", err)
	}

	RootCmd.SetArgs([]string{tmpDir}) // output to tmpDir (which contains the badFormatFile dir)
	b := bytes.NewBufferString("")
	RootCmd.SetOut(b)
	RootCmd.SetErr(b)

	err = RootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when writing to docker-compose.yml dir")
	}
}

func TestNoTestFlag(t *testing.T) {
	resetFlags()
	tmpDir := t.TempDir()
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(""), 0644)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	builder.ExecCommand = fakeExecCommand
	RootCmd.SetArgs([]string{"--no-test", tmpDir})
	b := bytes.NewBufferString("")
	RootCmd.SetOut(b)
	RootCmd.SetErr(b)

	err = RootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDryRunFlag(t *testing.T) {
	resetFlags()
	tmpDir := t.TempDir()
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(""), 0644)
	if err != nil {
		t.Fatalf("write error: %v", err)
	}

	builder.ExecCommand = fakeExecCommand
	RootCmd.SetArgs([]string{"--dry-run", tmpDir})
	b := bytes.NewBufferString("")
	RootCmd.SetOut(b)
	RootCmd.SetErr(b)

	err = RootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify files were not created
	expectedFiles := []string{
		"debian.Dockerfile",
		"alpine.Dockerfile",
		"distroless.Dockerfile",
		"docker-compose.yml",
		"Makefile",
		"make.bat",
		"BUILD",
	}
	for _, file := range expectedFiles {
		if _, err := os.Stat(filepath.Join(tmpDir, file)); !os.IsNotExist(err) {
			t.Errorf("expected file %s to NOT be created in dry-run mode", file)
		}
	}
}

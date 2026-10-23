package main

import (
	"bytes"
	"mkconf/builder"
	"os"
	"os/exec"
	"path/filepath"
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

func resetFlags(t *testing.T) *bytes.Buffer {
	t.Helper()
	origExec := builder.ExecCommand
	origBuilderOut := builder.Stdout
	origBuilderErr := builder.Stderr

	outputDir = ""
	noTest = false
	dryRun = false
	emitDockerfile = false
	emitBazel = false
	emitMakefile = false
	buildImages = false

	var cmdOut bytes.Buffer
	var builderOut bytes.Buffer
	RootCmd.SetOut(&cmdOut)
	RootCmd.SetErr(&cmdOut)
	builder.Stdout = &builderOut
	builder.Stderr = &builderOut

	t.Cleanup(func() {
		builder.ExecCommand = origExec
		builder.Stdout = origBuilderOut
		builder.Stderr = origBuilderErr
	})

	return &cmdOut
}

func TestMainSuccessAndDefaultArg(t *testing.T) {
	cmdOut := resetFlags(t)
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(""), 0644); err != nil {
		t.Fatalf("write error: %v", err)
	}

	builder.ExecCommand = fakeExecCommand

	// Test with explicit arg
	RootCmd.SetArgs([]string{tmpDir})
	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(cmdOut.String(), "Detected language: go") {
		t.Errorf("expected output to contain detected language go, got: %s", cmdOut.String())
	}

	// Test with 0 args (defaults to ".") inside tmpDir
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd error: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir error: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origWd)
	})

	cmdOut.Reset()
	RootCmd.SetArgs([]string{})
	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error on 0-arg execution: %v", err)
	}
	if !strings.Contains(cmdOut.String(), "Scanning repository: .") {
		t.Errorf("expected default scan path '.', got: %s", cmdOut.String())
	}
}

func TestMainFailureConditions(t *testing.T) {
	cmdOut := resetFlags(t)
	tmpDir := t.TempDir()

	// 1. Scan Error
	RootCmd.SetArgs([]string{filepath.Join(tmpDir, "nonexistent")})
	if err := RootCmd.Execute(); err == nil {
		t.Fatal("expected error from scan")
	}

	// 2. Test failure and BuildImage failure warnings
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(""), 0644); err != nil {
		t.Fatalf("write error: %v", err)
	}
	builder.ExecCommand = fakeExecCommandError
	cmdOut.Reset()
	RootCmd.SetArgs([]string{"--build", tmpDir})
	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("expected nil from Execute on warnings, got %v", err)
	}
	if !strings.Contains(cmdOut.String(), "Warning: Tests failed:") {
		t.Errorf("expected test failure warning in output, got: %s", cmdOut.String())
	}
	if !strings.Contains(cmdOut.String(), "Failed to build app-debian:") {
		t.Errorf("expected image build failure warning in output, got: %s", cmdOut.String())
	}

	// 3. BuildImage success path
	builder.ExecCommand = fakeExecCommand
	cmdOut.Reset()
	RootCmd.SetArgs([]string{"--no-test", "--build", tmpDir})
	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("expected nil on successful build, got %v", err)
	}
	if strings.Contains(cmdOut.String(), "Failed to build") {
		t.Errorf("did not expect build failure when mock succeeds, got: %s", cmdOut.String())
	}
}

func TestEmitFlagsSelective(t *testing.T) {
	cmdOut := resetFlags(t)
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(""), 0644); err != nil {
		t.Fatalf("write error: %v", err)
	}
	builder.ExecCommand = fakeExecCommand

	// Only emit Bazel BUILD file (skips Dockerfile and Makefile)
	outBazelOnly := filepath.Join(tmpDir, "bazel_only")
	RootCmd.SetArgs([]string{"--no-test", "--emit-bazel-build-file", "--output", outBazelOnly, tmpDir})
	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outBazelOnly, "BUILD")); err != nil {
		t.Errorf("expected BUILD to be created: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outBazelOnly, "Makefile")); !os.IsNotExist(err) {
		t.Error("did not expect Makefile when only --emit-bazel-build-file is set")
	}
	if _, err := os.Stat(filepath.Join(outBazelOnly, "debian.Dockerfile")); !os.IsNotExist(err) {
		t.Error("did not expect debian.Dockerfile when only --emit-bazel-build-file is set")
	}
	if !strings.Contains(cmdOut.String(), "Saved "+filepath.Join(outBazelOnly, "BUILD")) {
		t.Errorf("expected Saved BUILD message, got: %s", cmdOut.String())
	}
}

func TestRunPath(t *testing.T) {
	resetFlags(t)
	exited := false
	osExit = func(_ int) {
		exited = true
	}
	defer func() { osExit = os.Exit }()

	RootCmd.SetArgs([]string{"a", "b"})
	Run()
	if !exited {
		t.Error("expected osExit to be called")
	}
}

func TestMainExecution(t *testing.T) {
	resetFlags(t)
	exited := false
	osExit = func(_ int) {
		exited = true
	}
	defer func() { osExit = os.Exit }()

	RootCmd.SetArgs([]string{"a", "b"})
	main()
	if !exited {
		t.Error("expected main to call osExit via Run")
	}
}

func TestMainBranches(t *testing.T) {
	cmdOut := resetFlags(t)
	tmpDir := t.TempDir()
	// Trigger python to hit InstallCommand output branch
	if err := os.WriteFile(filepath.Join(tmpDir, "requirements.txt"), []byte(""), 0644); err != nil {
		t.Fatalf("write error: %v", err)
	}

	builder.ExecCommand = fakeExecCommand
	RootCmd.SetArgs([]string{tmpDir})
	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(cmdOut.String(), "Install Command: pip install -r requirements.txt") {
		t.Errorf("expected python Install Command in output, got: %s", cmdOut.String())
	}
}

func TestOutputDirectoryAndWriteErrors(t *testing.T) {
	resetFlags(t)
	tmpDir := t.TempDir()
	outDir := filepath.Join(tmpDir, "out")
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(""), 0644); err != nil {
		t.Fatalf("write error: %v", err)
	}

	builder.ExecCommand = fakeExecCommand

	RootCmd.SetArgs([]string{"--output", outDir, tmpDir})
	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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

	// 1. Output directory creation failure test
	resetFlags(t)
	badOut := filepath.Join(tmpDir, "bad_out")
	_ = os.WriteFile(badOut, []byte("file, not dir"), 0644)
	RootCmd.SetArgs([]string{"--output", badOut, tmpDir})
	if err := RootCmd.Execute(); err == nil {
		t.Fatal("expected error from MkdirAll on existing file")
	}

	// 2. Dockerfile WriteFile failure test
	resetFlags(t)
	readonlyFile := filepath.Join(tmpDir, "readonly_out")
	_ = os.MkdirAll(filepath.Join(readonlyFile, "debian.Dockerfile"), 0755)
	RootCmd.SetArgs([]string{"--output", readonlyFile, tmpDir})
	if err := RootCmd.Execute(); err == nil {
		t.Fatal("expected error from WriteFile on debian.Dockerfile dir")
	}

	// 3. docker-compose.yml WriteFile failure test
	resetFlags(t)
	dcFailDir := filepath.Join(tmpDir, "dc_fail_out")
	_ = os.MkdirAll(filepath.Join(dcFailDir, "docker-compose.yml"), 0755)
	RootCmd.SetArgs([]string{"--output", dcFailDir, tmpDir})
	if err := RootCmd.Execute(); err == nil {
		t.Fatal("expected error when writing to docker-compose.yml dir")
	}

	// 4. formats loop (Makefile) WriteFile failure test
	resetFlags(t)
	makeFailDir := filepath.Join(tmpDir, "make_fail_out")
	_ = os.MkdirAll(filepath.Join(makeFailDir, "Makefile"), 0755)
	RootCmd.SetArgs([]string{"--output", makeFailDir, tmpDir})
	if err := RootCmd.Execute(); err == nil {
		t.Fatal("expected error when writing to Makefile dir in formats loop")
	}
}

func TestNoTestAndDryRunFlags(t *testing.T) {
	cmdOut := resetFlags(t)
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(""), 0644); err != nil {
		t.Fatalf("write error: %v", err)
	}

	builder.ExecCommand = fakeExecCommand
	RootCmd.SetArgs([]string{"--no-test", "--dry-run", "--build", tmpDir})

	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := cmdOut.String()
	if !strings.Contains(output, "Skipping tests due to --no-test flag") {
		t.Errorf("expected --no-test message, got: %s", output)
	}
	if !strings.Contains(output, "Would build app-debian") {
		t.Errorf("expected dry-run Would build message, got: %s", output)
	}
	if !strings.Contains(output, "Would save "+filepath.Join(tmpDir, "Makefile")) {
		t.Errorf("expected dry-run Would save Makefile message, got: %s", output)
	}

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

package sdk

import (
	"bytes"
	"mkconf/builder"
	"mkconf/scanner"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func fakeExecSuccess(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestSDKHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_SDK_HELPER_PROCESS=1"}
	return cmd
}

func fakeExecFailure(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestSDKHelperProcessError", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_SDK_HELPER_PROCESS_ERR=1"}
	return cmd
}

func TestSDKHelperProcess(_ *testing.T) {
	if os.Getenv("GO_WANT_SDK_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}

func TestSDKHelperProcessError(_ *testing.T) {
	if os.Getenv("GO_WANT_SDK_HELPER_PROCESS_ERR") != "1" {
		return
	}
	os.Exit(1)
}

func TestSDKGenerators(t *testing.T) {
	project := &Project{
		Path: ".",
		Info: &scanner.ProjectInfo{
			Language:     "go",
			StartCommand: []string{"/app/app"},
		},
	}

	if out := project.GenerateDockerCompose(); !strings.Contains(out, "version:") || !strings.Contains(out, "debian.Dockerfile") {
		t.Errorf("GenerateDockerCompose failed, got: %s", out)
	}
	if out := project.GenerateMakefile(); !strings.Contains(out, "help:") || !strings.Contains(out, "build:") {
		t.Errorf("GenerateMakefile failed, got: %s", out)
	}
	if out := project.GenerateMakeBat(); !strings.Contains(out, "@echo off") {
		t.Errorf("GenerateMakeBat failed, got: %s", out)
	}
	if out := project.GenerateBazelBuild(); !strings.Contains(out, "BUILD") {
		t.Errorf("GenerateBazelBuild failed, got: %s", out)
	}
}

func TestNew(t *testing.T) {
	proj, err := New("..")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if proj.Info == nil || proj.Info.Language != "go" {
		t.Errorf("Expected Info with go language, got %+v", proj.Info)
	}

	out := proj.GenerateDockerfile("debian")
	if !strings.Contains(out, "FROM golang:") {
		t.Errorf("Expected golang Dockerfile output, got: %s", out)
	}
}

func TestRunTestEmpty(t *testing.T) {
	project := &Project{
		Path: ".",
		Info: &scanner.ProjectInfo{},
	}
	if err := project.RunTest(); err != nil {
		t.Errorf("Expected no error for empty test command, got %v", err)
	}
}

func TestRunTestMocked(t *testing.T) {
	origExec := builder.ExecCommand
	origOut := builder.Stdout
	origErr := builder.Stderr
	t.Cleanup(func() {
		builder.ExecCommand = origExec
		builder.Stdout = origOut
		builder.Stderr = origErr
	})

	var buf bytes.Buffer
	builder.Stdout = &buf
	builder.Stderr = &buf

	project := &Project{
		Path: t.TempDir(),
		Info: &scanner.ProjectInfo{
			TestCommand: "go test ./...",
		},
	}

	builder.ExecCommand = fakeExecSuccess
	if err := project.RunTest(); err != nil {
		t.Errorf("Expected no error for mocked test run, got %v", err)
	}

	builder.ExecCommand = fakeExecFailure
	if err := project.RunTest(); err == nil {
		t.Error("Expected error when mocked test command fails")
	}
}

func TestBuildImageMocked(t *testing.T) {
	origExec := builder.ExecCommand
	origOut := builder.Stdout
	origErr := builder.Stderr
	t.Cleanup(func() {
		builder.ExecCommand = origExec
		builder.Stdout = origOut
		builder.Stderr = origErr
	})

	var buf bytes.Buffer
	builder.Stdout = &buf
	builder.Stderr = &buf

	project := &Project{
		Path: t.TempDir(),
		Info: &scanner.ProjectInfo{},
	}

	builder.ExecCommand = fakeExecSuccess
	if err := project.BuildImage("FROM scratch", "app-debian"); err != nil {
		t.Errorf("Expected BuildImage to succeed with mock, got %v", err)
	}

	builder.ExecCommand = fakeExecFailure
	if err := project.BuildImage("FROM scratch", "app-debian"); err == nil {
		t.Error("Expected BuildImage to return error when mock fails")
	}
}

func TestNewError(t *testing.T) {
	_, err := New("/path/that/does/not/exist")
	if err == nil {
		t.Error("Expected error for missing directory")
	}
}

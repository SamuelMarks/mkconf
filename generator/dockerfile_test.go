package generator

import (
	"mkconf/scanner"
	"strings"
	"testing"
)

func TestGenerateDockerfile(t *testing.T) {
	// Normal case with known language
	info := &scanner.ProjectInfo{
		Language:       "go",
		InstallCommand: "go mod download",
		BuildCommand:   "go build -o app",
		TestCommand:    "go test ./...",
		StartCommand:   []string{"/app/app"},
	}

	debian := GenerateDockerfile(info, "debian")
	if !strings.Contains(debian, "FROM golang:") {
		t.Errorf("Expected FROM golang in debian Dockerfile")
	}

	alpine := GenerateDockerfile(info, "alpine")
	if !strings.Contains(alpine, "FROM golang:") {
		t.Errorf("Expected FROM golang in alpine Dockerfile")
	}

	distroless := GenerateDockerfile(info, "distroless")
	if !strings.Contains(distroless, "FROM golang:") || !strings.Contains(distroless, "AS builder") {
		t.Errorf("Expected multi-stage build in distroless Dockerfile")
	}

	// Unknown language
	unknownInfo := &scanner.ProjectInfo{
		Language: "unknown",
	}
	unknownOutput := GenerateDockerfile(unknownInfo, "debian")
	if !strings.Contains(unknownOutput, "FROM ubuntu:") {
		t.Errorf("Expected fallback to ubuntu for unknown language")
	}

	// Uncompiled language distroless
	nodeInfo := &scanner.ProjectInfo{
		Language: "nodejs",
		StartCommand: []string{"npm", "start"},
		InstallCommand: "npm install",
		BuildCommand: "npm run build",
	}
	distrolessNode := GenerateDockerfile(nodeInfo, "distroless")
	if !strings.Contains(distrolessNode, "CMD") {
		t.Errorf("Expected CMD for uncompiled distroless")
	}

	// Uncompiled alpine
	alpineNode := GenerateDockerfile(nodeInfo, "alpine")
	if !strings.Contains(alpineNode, "CMD") {
		t.Errorf("Expected CMD for uncompiled alpine")
	}
}

func TestGenerateDockerfileEmptyBase(t *testing.T) {
	// Base image becomes empty when baseType is not recognized
	info := &scanner.ProjectInfo{
		Language:       "go",
	}
	output := GenerateDockerfile(info, "unknown_base")
	if !strings.Contains(output, "FROM ubuntu:22.04") {
		t.Errorf("Expected ubuntu:22.04 when base image is empty")
	}
}

func TestGenerateDockerfileAlpinePackages(t *testing.T) {
	info := &scanner.ProjectInfo{
		Language:       "c",
		StartCommand:   []string{"/app/app"},
	}
	output := GenerateDockerfile(info, "alpine")
	if !strings.Contains(output, "apk add --no-cache") {
		t.Errorf("Expected apk add in alpine Dockerfile for C")
	}
}

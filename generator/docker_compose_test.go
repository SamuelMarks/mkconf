package generator

import (
	"mkconf/scanner"
	"strings"
	"testing"
)

func TestDockerCompose(t *testing.T) {
	info := &scanner.ProjectInfo{}
	compose := GenerateDockerCompose(info)

	if !strings.Contains(compose, "version: \"3.8\"") && !strings.Contains(compose, "version: '3.8'") && !strings.Contains(compose, "version: 3.8") {
		t.Errorf("Expected Docker Compose output to contain version, got: %s", compose)
	}
	if !strings.Contains(compose, "services:") {
		t.Errorf("Expected Docker Compose output to contain services, got: %s", compose)
	}
	if !strings.Contains(compose, "app:") {
		t.Errorf("Expected Docker Compose output to contain app, got: %s", compose)
	}
	if !strings.Contains(compose, "build: .") {
		t.Errorf("Expected Docker Compose output to contain build, got: %s", compose)
	}
}

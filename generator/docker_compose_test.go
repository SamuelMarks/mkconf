package generator

import (
	"mkconf/scanner"
	"strings"
	"testing"
)

func TestDockerCompose(t *testing.T) {
	info := &scanner.ProjectInfo{
		Language:     "go",
		StartCommand: []string{"/app/app"},
	}
	compose := GenerateDockerCompose(info)

	if !strings.Contains(compose, `version: "3.8"`) && !strings.Contains(compose, "version: '3.8'") && !strings.Contains(compose, "version: 3.8") {
		t.Errorf("Expected Docker Compose output to contain version, got: %s", compose)
	}
	if !strings.Contains(compose, "services:") {
		t.Errorf("Expected Docker Compose output to contain services, got: %s", compose)
	}
	if !strings.Contains(compose, "app:") {
		t.Errorf("Expected Docker Compose output to contain app, got: %s", compose)
	}
	if !strings.Contains(compose, "context: .") || !strings.Contains(compose, "dockerfile: debian.Dockerfile") {
		t.Errorf("Expected Docker Compose output to reference context and debian.Dockerfile, got: %s", compose)
	}
	if !strings.Contains(compose, "/app/app") {
		t.Errorf("Expected Docker Compose output to include start command from info, got: %s", compose)
	}

	emptyCompose := GenerateDockerCompose(nil)
	if strings.Contains(emptyCompose, "command:") {
		t.Errorf("Did not expect command field when info is nil, got: %s", emptyCompose)
	}
}

package generator

import (
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"testing"
)

func TestGenerateString(t *testing.T) {
	if GenerateString(nil) != "" {
		t.Error("expected empty string for nil node")
	}

	root := &parser.Node{Value: "root"}
	addInstruction(root, "FROM", "ubuntu")
	addInstruction(root, "RUN", "echo", "hello")

	expected := "FROM ubuntu\nRUN echo hello\n"
	if out := GenerateString(root); out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}

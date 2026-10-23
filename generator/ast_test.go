package generator

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func TestGenerateString(t *testing.T) {
	if GenerateString(nil) != "" {
		t.Error("expected empty string for nil node")
	}

	root := &parser.Node{Value: "root"}
	addInstruction(root, "FROM", "ubuntu")
	addInstruction(root, "RUN", "echo", "hello")

	expected := "FROM ubuntu" + newline + "RUN echo hello" + newline
	if out := GenerateString(root); out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}

	rootMulti := &parser.Node{Value: "root"}
	addInstruction(rootMulti, "FROM", "ubuntu", "AS", "builder")
	addInstruction(rootMulti, "RUN", "echo", "hello")
	addInstruction(rootMulti, "FROM", "scratch")
	addInstruction(rootMulti, "COPY", "--from=builder", "/app", "/app")

	expectedMulti := "FROM ubuntu AS builder" + newline + "RUN echo hello" + newline + newline + "FROM scratch" + newline + "COPY --from=builder /app /app" + newline
	if out := GenerateString(rootMulti); out != expectedMulti {
		t.Errorf("expected %q, got %q", expectedMulti, out)
	}

	res := resolveProjectCommands(nil, "echo", "'")
	if res.lang != "Base" {
		t.Errorf("expected Base for nil info, got %s", res.lang)
	}
}

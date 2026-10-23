// Package generator provides functionality for generating configuration files.
package generator

import (
	"mkconf/scanner"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// newline holds a single line feed character for formatting multi-line outputs.
const newline = `
`

// resolvedCommands holds normalized lifecycle commands and display language name for generators.
type resolvedCommands struct {
	// lang is the title-cased language name (defaults to "Base").
	lang string
	// installCmd is the resolved dependency installation command or echo fallback.
	installCmd string
	// buildCmd is the resolved build command or echo fallback.
	buildCmd string
	// testCmd is the resolved test command or echo fallback.
	testCmd string
	// runCmd is the resolved start command string or echo fallback.
	runCmd string
}

// resolveProjectCommands normalizes a ProjectInfo into fallback-safe commands for build scripts.
//
// Parameters:
//   - info: The scanned ProjectInfo metadata (may be nil).
//   - echoPrefix: The echo command prefix (e.g., "@echo", "echo").
//   - quoteStyle: The quote character to wrap fallback messages.
func resolveProjectCommands(info *scanner.ProjectInfo, echoPrefix string, quoteStyle string) resolvedCommands {
	var safe scanner.ProjectInfo
	if info != nil {
		safe = *info
	}

	formatFallback := func(msg string) string {
		return echoPrefix + " " + quoteStyle + msg + quoteStyle
	}

	installCmd := safe.InstallCommand
	if installCmd == "" {
		installCmd = formatFallback("No install dependencies command defined")
	}

	buildCmd := safe.BuildCommand
	if buildCmd == "" {
		buildCmd = formatFallback("No build command defined")
	}

	testCmd := safe.TestCommand
	if testCmd == "" {
		testCmd = formatFallback("No test command defined")
	}

	runCmd := strings.Join(safe.StartCommand, " ")
	if runCmd == "" {
		runCmd = formatFallback("No start command defined")
	}

	lang := safe.Language
	if lang == "" {
		lang = "base"
	}
	lang = strings.ToUpper(lang[:1]) + lang[1:]

	return resolvedCommands{
		lang:       lang,
		installCmd: installCmd,
		buildCmd:   buildCmd,
		testCmd:    testCmd,
		runCmd:     runCmd,
	}
}

// GenerateString recursively formats a parser.Node and its children into a
// string representation of a Dockerfile.
//
// Parameters:
//   - node: The root *parser.Node AST representing the Dockerfile instructions.
func GenerateString(node *parser.Node) string {
	if node == nil {
		return ""
	}
	var sb strings.Builder

	for i, child := range node.Children {
		if i > 0 && child.Value == "FROM" {
			sb.WriteString(newline)
		}
		sb.WriteString(child.Value)
		curr := child.Next
		for curr != nil {
			sb.WriteString(" ")
			sb.WriteString(curr.Value)
			curr = curr.Next
		}
		sb.WriteString(newline)
	}
	return sb.String()
}

// addInstruction appends a new Dockerfile instruction node and its argument nodes to root.
//
// Parameters:
//   - root: The parent *parser.Node to append the instruction child to.
//   - instruction: The uppercase Dockerfile directive (e.g., "FROM", "RUN", "COPY").
//   - args: Variadic string arguments following the instruction.
func addInstruction(root *parser.Node, instruction string, args ...string) {
	node := &parser.Node{
		Value: instruction,
	}

	var prev *parser.Node
	for _, arg := range args {
		argNode := &parser.Node{Value: arg}
		if prev == nil {
			node.Next = argNode
		} else {
			prev.Next = argNode
		}
		prev = argNode
	}

	root.Children = append(root.Children, node)
}

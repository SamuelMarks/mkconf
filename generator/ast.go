// Package generator provides utilities for constructing and formatting Dockerfiles.
package generator

import (
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"strings"
)

// GenerateString recursively formats a parser.Node and its children into a
// string representation of a Dockerfile.
func GenerateString(node *parser.Node) string {
	if node == nil {
		return ""
	}
	var sb strings.Builder

	for i, child := range node.Children {
		if i > 0 && child.Value == "FROM" {
			sb.WriteString("\n")
		}
		sb.WriteString(child.Value)
		curr := child.Next
		for curr != nil {
			sb.WriteString(" ")
			sb.WriteString(curr.Value)
			curr = curr.Next
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// addInstruction is a helper that appends a new instruction and its arguments
// to the root node's children.
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

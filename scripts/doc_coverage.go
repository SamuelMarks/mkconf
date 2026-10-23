// Package main provides an AST-based analyzer to calculate documentation coverage
// across files, packages, types, struct fields, constants, variables, functions, and arguments.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// rootDir specifies the default root directory scanned by main.
var rootDir = "."

// isSkippedDir reports whether a directory name should be ignored during documentation scanning.
//
// Parameters:
//   - name: The base directory name to evaluate.
func isSkippedDir(name string) bool {
	if name != "." && name != ".." && strings.HasPrefix(name, ".") {
		return true
	}
	return name == "vendor" || name == "testdata"
}

// inspectFuncDecl counts documentation coverage for a function or method and all of its named parameters.
//
// Parameters:
//   - d: The *ast.FuncDecl AST node to inspect.
//
// Returns the total number of documentable items and the number of documented items.
func inspectFuncDecl(d *ast.FuncDecl) (int, int) {
	total := 1
	documented := 0
	docText := ""
	if d.Doc != nil && strings.TrimSpace(d.Doc.Text()) != "" {
		documented++
		docText = d.Doc.Text()
	}

	if d.Type.Params != nil {
		for _, field := range d.Type.Params.List {
			for _, name := range field.Names {
				if name.Name == "_" {
					continue
				}
				total++
				if strings.Contains(docText, name.Name) {
					documented++
				}
			}
		}
	}

	return total, documented
}

// inspectGenDecl counts documentation coverage for types, struct fields, constants, and variables.
//
// Parameters:
//   - d: The *ast.GenDecl AST node to inspect.
//
// Returns the total number of documentable items and the number of documented items.
func inspectGenDecl(d *ast.GenDecl) (int, int) {
	var total, documented int
	for _, spec := range d.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			total++
			if d.Doc != nil || s.Doc != nil {
				documented++
			}
			if st, ok := s.Type.(*ast.StructType); ok && st.Fields != nil {
				for _, field := range st.Fields.List {
					count := len(field.Names)
					if count == 0 {
						count = 1
					}
					total += count
					if field.Doc != nil || field.Comment != nil {
						documented += count
					}
				}
			}
		case *ast.ValueSpec:
			for range s.Names {
				total++
				if d.Doc != nil || s.Doc != nil || s.Comment != nil {
					documented++
				}
			}
		}
	}
	return total, documented
}

// run walks the directory tree rooted at root and returns the documentation coverage percentage.
//
// Parameters:
//   - root: The root directory path to recursively analyze.
func run(root string) (float64, error) {
	var total, documented int

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if isSkippedDir(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".go") || strings.HasSuffix(info.Name(), "_test.go") {
			return nil
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		total++
		if f.Doc != nil {
			documented++
		}

		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				tCount, dCount := inspectFuncDecl(d)
				total += tCount
				documented += dCount
			case *ast.GenDecl:
				tCount, dCount := inspectGenDecl(d)
				total += tCount
				documented += dCount
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	if total == 0 {
		return 100.0, nil
	}

	return float64(documented) / float64(total) * 100, nil
}

// osExit points to os.Exit to allow overriding termination behavior in unit tests.
var osExit = os.Exit

// main executes the documentation coverage analyzer and prints the result to stdout.
func main() {
	cov, err := run(rootDir)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		osExit(1)
		return
	}
	_, _ = fmt.Fprintln(os.Stdout, strconv.FormatFloat(cov, 'f', 1, 64))
}

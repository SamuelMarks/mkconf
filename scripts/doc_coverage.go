// Package main provides a script to calculate documentation coverage for Go source files.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

var rootDir = "."

func run(root string) (float64, error) {
	var total, documented int

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if (info.Name() != "." && info.Name() != ".." && strings.HasPrefix(info.Name(), ".")) || info.Name() == "vendor" || info.Name() == "scripts" || info.Name() == "testdata" {
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
				if d.Name.IsExported() {
					total++
					if d.Doc != nil {
						documented++
					}
				}
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						if s.Name.IsExported() {
							total++
							if d.Doc != nil || s.Doc != nil {
								documented++
							}
						}
					case *ast.ValueSpec:
						for _, name := range s.Names {
							if name.IsExported() {
								total++
								if d.Doc != nil || s.Doc != nil {
									documented++
								}
							}
						}
					}
				}
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

// override exit for testing
var osExit = os.Exit

func main() {
	cov, err := run(rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		osExit(1)
		return
	}
	fmt.Printf("%.1f\n", cov)
}

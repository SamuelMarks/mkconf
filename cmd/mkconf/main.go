// Package main is the entry point for the mkconf command line application.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"mkconf/sdk"
)

var osExit = os.Exit
var outputDir string

// noTest determines if the test suite should be skipped.
var noTest bool

// dryRun determines if file writes and image builds should be skipped.
var dryRun bool

var emitDockerfile bool
var emitBazel bool
var emitMakefile bool
var buildImages bool

func init() {
	RootCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for generated files (defaults to repo_path)")
	RootCmd.Flags().BoolVar(&noTest, "no-test", false, "Skip running the project's test suite")
	RootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Do not write files or build images, only print output")

	RootCmd.Flags().BoolVar(&emitDockerfile, "emit-dockerfile", false, "Emit Dockerfile(s)")
	RootCmd.Flags().BoolVar(&emitBazel, "emit-bazel-build-file", false, "Emit Bazel BUILD file")
	RootCmd.Flags().BoolVar(&emitMakefile, "emit-makefile", false, "Emit Makefile")
	RootCmd.Flags().BoolVar(&buildImages, "build", false, "Build Docker images after generating Dockerfiles (requires Docker)")
}

func shouldEmit(flag bool, emitAll bool) bool {
	return emitAll || flag
}

// RootCmd is the main entry point for the application commands
var RootCmd = &cobra.Command{
	Use:   "mkconf [repo_path]",
	Short: "Generates build and environment configuration files",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		fmt.Printf("Scanning repository: %s\n", path)

		project, err := sdk.New(path)
		if err != nil {
			return fmt.Errorf("error scanning repository: %v", err)
		}

		info := project.Info
		fmt.Printf("Detected language: %s\n", info.Language)
		if info.InstallCommand != "" {
			fmt.Printf("Install Command: %s\n", info.InstallCommand)
		}
		if info.BuildCommand != "" {
			fmt.Printf("Build Command: %s\n", info.BuildCommand)
		}
		if info.TestCommand != "" {
			fmt.Printf("Test Command: %s\n", info.TestCommand)
		}

		if !noTest {
			if err := project.RunTest(); err != nil {
				fmt.Printf("Warning: Tests failed: %v\n", err)
			}
		} else {
			fmt.Println("Skipping tests due to --no-test flag")
		}

		targetDir := outputDir
		if targetDir == "" {
			targetDir = path
		}

		if !dryRun {
			if err := os.MkdirAll(targetDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %v", err)
			}
		}

		emitAll := !emitDockerfile && !emitBazel && !emitMakefile

		baseImages := []string{"debian", "alpine", "distroless"}

		if shouldEmit(emitDockerfile, emitAll) {
			for _, base := range baseImages {
				dockerfile := project.GenerateDockerfile(base)
				filePath := filepath.Join(targetDir, fmt.Sprintf("%s.Dockerfile", base))

				imageName := fmt.Sprintf("%s-%s", "app", base)

				if dryRun {
					fmt.Printf("Would save %s\n", filePath)
					if buildImages {
						fmt.Printf("Would build %s\n", imageName)
					}
				} else {
					if err := os.WriteFile(filePath, []byte(dockerfile), 0644); err != nil {
						return fmt.Errorf("failed to write %s: %v", filePath, err)
					}
					fmt.Printf("Saved %s\n", filePath)

					if buildImages {
						err = project.BuildImage(dockerfile, imageName)
						if err != nil {
							fmt.Printf("Failed to build %s: %v\n", imageName, err)
						}
					}
				}
			}

			// Docker compose kind of belongs with Dockerfile
			dcPath := filepath.Join(targetDir, "docker-compose.yml")
			if dryRun {
				fmt.Printf("Would save %s\n", dcPath)
			} else {
				if err := os.WriteFile(dcPath, []byte(project.GenerateDockerCompose()), 0644); err != nil {
					return fmt.Errorf("failed to write %s: %v", dcPath, err)
				}
				fmt.Printf("Saved %s\n", dcPath)
			}
		}

		formats := []struct {
			filename string
			content  string
			emit     bool
		}{
			{"Makefile", project.GenerateMakefile(), shouldEmit(emitMakefile, emitAll)},
			{"make.bat", project.GenerateMakeBat(), shouldEmit(emitMakefile, emitAll)},
			{"BUILD", project.GenerateBazelBuild(), shouldEmit(emitBazel, emitAll)},
		}

		for _, f := range formats {
			if !f.emit {
				continue
			}
			filePath := filepath.Join(targetDir, f.filename)
			if dryRun {
				fmt.Printf("Would save %s\n", filePath)
			} else {
				if err := os.WriteFile(filePath, []byte(f.content), 0644); err != nil {
					return fmt.Errorf("failed to write %s: %v", filePath, err)
				}
				fmt.Printf("Saved %s\n", filePath)
			}
		}

		return nil
	},
}

// Run executes the RootCmd and exits appropriately on error.
func Run() {
	if err := RootCmd.Execute(); err != nil {
		osExit(1)
	}
}

func main() {
	Run()
}

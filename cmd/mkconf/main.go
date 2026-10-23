// Package main is the entry point for the mkconf command line application.
package main

import (
	"fmt"
	"mkconf/config"
	"mkconf/sdk"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// osExit points to os.Exit to allow testing exit status codes without terminating the test runner.
var osExit = os.Exit

// outputDir stores the user-specified output directory for generated files.
var outputDir string

// noTest determines if the project's test suite should be skipped.
var noTest bool

// dryRun determines if file writes and image builds should be skipped.
var dryRun bool

// emitDockerfile determines if Dockerfiles and docker-compose.yml should be emitted.
var emitDockerfile bool

// emitBazel determines if the Bazel BUILD file should be emitted.
var emitBazel bool

// emitMakefile determines if Makefile and make.bat should be emitted.
var emitMakefile bool

// buildImages determines if Docker images should be built after generating Dockerfiles.
var buildImages bool

// formatSpec describes a generated configuration file and whether it should be emitted.
type formatSpec struct {
	// filename is the output filename within targetDir.
	filename string
	// content is the generated text content of the file.
	content string
	// emit indicates whether this format is selected for output.
	emit bool
}

// init registers command-line flags on RootCmd.
func init() {
	RootCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for generated files (defaults to repo_path)")
	RootCmd.Flags().BoolVar(&noTest, "no-test", false, "Skip running the project's test suite")
	RootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Do not write files or build images, only print output")

	RootCmd.Flags().BoolVar(&emitDockerfile, "emit-dockerfile", false, "Emit Dockerfile(s)")
	RootCmd.Flags().BoolVar(&emitBazel, "emit-bazel-build-file", false, "Emit Bazel BUILD file")
	RootCmd.Flags().BoolVar(&emitMakefile, "emit-makefile", false, "Emit Makefile")
	RootCmd.Flags().BoolVar(&buildImages, "build", false, "Build Docker images after generating Dockerfiles (requires Docker)")
}

// shouldEmit reports whether a specific format flag is active or all formats are being emitted.
//
// Parameters:
//   - flag: Whether the specific format flag was explicitly enabled.
//   - emitAll: Whether no individual --emit-* flags were specified.
func shouldEmit(flag bool, emitAll bool) bool {
	return emitAll || flag
}

// executeRoot handles the execution of the mkconf CLI command.
//
// Parameters:
//   - cmd: The executing *cobra.Command used for formatted output.
//   - args: Positional command-line arguments containing the optional repository path.
func executeRoot(cmd *cobra.Command, args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}
	cmd.Println(fmt.Sprintf("Scanning repository: %s", path))

	project, err := sdk.New(path)
	if err != nil {
		return fmt.Errorf("error scanning repository: %w", err)
	}

	info := project.Info
	cmd.Println(fmt.Sprintf("Detected language: %s", info.Language))
	if info.InstallCommand != "" {
		cmd.Println(fmt.Sprintf("Install Command: %s", info.InstallCommand))
	}
	if info.BuildCommand != "" {
		cmd.Println(fmt.Sprintf("Build Command: %s", info.BuildCommand))
	}
	if info.TestCommand != "" {
		cmd.Println(fmt.Sprintf("Test Command: %s", info.TestCommand))
	}

	if !noTest {
		if err := project.RunTest(); err != nil {
			cmd.Println(fmt.Sprintf("Warning: Tests failed: %v", err))
		}
	} else {
		cmd.Println("Skipping tests due to --no-test flag")
	}

	targetDir := outputDir
	if targetDir == "" {
		targetDir = path
	}

	if !dryRun {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	emitAll := !emitDockerfile && !emitBazel && !emitMakefile

	if shouldEmit(emitDockerfile, emitAll) {
		for _, base := range config.AllBaseImages() {
			baseStr := string(base)
			dockerfile := project.GenerateDockerfile(baseStr)
			filePath := filepath.Join(targetDir, fmt.Sprintf("%s.Dockerfile", baseStr))
			imageName := fmt.Sprintf("app-%s", baseStr)

			if dryRun {
				cmd.Println(fmt.Sprintf("Would save %s", filePath))
				if buildImages {
					cmd.Println(fmt.Sprintf("Would build %s", imageName))
				}
			} else {
				if err := os.WriteFile(filePath, []byte(dockerfile), 0644); err != nil {
					return fmt.Errorf("failed to write %s: %w", filePath, err)
				}
				cmd.Println(fmt.Sprintf("Saved %s", filePath))

				if buildImages {
					if err := project.BuildImage(dockerfile, imageName); err != nil {
						cmd.Println(fmt.Sprintf("Failed to build %s: %v", imageName, err))
					}
				}
			}
		}

		dcPath := filepath.Join(targetDir, string(config.FormatDockerCompose))
		if dryRun {
			cmd.Println(fmt.Sprintf("Would save %s", dcPath))
		} else {
			if err := os.WriteFile(dcPath, []byte(project.GenerateDockerCompose()), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %w", dcPath, err)
			}
			cmd.Println(fmt.Sprintf("Saved %s", dcPath))
		}
	}

	formats := []formatSpec{
		{string(config.FormatMakefile), project.GenerateMakefile(), shouldEmit(emitMakefile, emitAll)},
		{string(config.FormatMakeBat), project.GenerateMakeBat(), shouldEmit(emitMakefile, emitAll)},
		{string(config.FormatBazel), project.GenerateBazelBuild(), shouldEmit(emitBazel, emitAll)},
	}

	for _, f := range formats {
		if !f.emit {
			continue
		}
		filePath := filepath.Join(targetDir, f.filename)
		if dryRun {
			cmd.Println(fmt.Sprintf("Would save %s", filePath))
		} else {
			if err := os.WriteFile(filePath, []byte(f.content), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %w", filePath, err)
			}
			cmd.Println(fmt.Sprintf("Saved %s", filePath))
		}
	}

	return nil
}

// RootCmd is the main Cobra command entry point for the mkconf CLI application.
var RootCmd = &cobra.Command{
	Use:   "mkconf [repo_path]",
	Short: "Generates build and environment configuration files",
	Args:  cobra.MaximumNArgs(1),
	RunE:  executeRoot,
}

// Run executes RootCmd and terminates with status code 1 on error.
func Run() {
	if err := RootCmd.Execute(); err != nil {
		osExit(1)
	}
}

// main invokes Run as the binary entrypoint.
func main() {
	Run()
}

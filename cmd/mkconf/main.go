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

func init() {
	RootCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory for generated files (defaults to repo_path)")
}

// RootCmd is the main entry point for the application commands
var RootCmd = &cobra.Command{
	Use:   "mkconf [repo_path]",
	Short: "Generates build and environment configuration files",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
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

		if err := project.RunTest(); err != nil {
			fmt.Printf("Warning: Tests failed: %v\n", err)
		}

		targetDir := outputDir
		if targetDir == "" {
			targetDir = path
		}

		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}

		baseImages := []string{"debian", "alpine", "distroless"}

		for _, base := range baseImages {
			dockerfile := project.GenerateDockerfile(base)
			filePath := filepath.Join(targetDir, fmt.Sprintf("%s.Dockerfile", base))
			if err := os.WriteFile(filePath, []byte(dockerfile), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", filePath, err)
			}
			fmt.Printf("Saved %s\n", filePath)

			imageName := fmt.Sprintf("%s-%s", "app", base)
			err = project.BuildImage(dockerfile, imageName)
			if err != nil {
				fmt.Printf("Failed to build %s: %v\n", imageName, err)
			}
		}

		formats := []struct {
			filename string
			content  string
		}{
			{"docker-compose.yml", project.GenerateDockerCompose()},
			{"Makefile", project.GenerateMakefile()},
			{"make.bat", project.GenerateMakeBat()},
			{"BUILD", project.GenerateBazelBuild()},
		}

		for _, f := range formats {
			filePath := filepath.Join(targetDir, f.filename)
			if err := os.WriteFile(filePath, []byte(f.content), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", filePath, err)
			}
			fmt.Printf("Saved %s\n", filePath)
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

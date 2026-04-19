// Package builder provides functionality to interact with the Docker daemon and shell
// to execute image builds and run tests.
package builder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ExecCommand is a variable pointing to exec.Command to allow mocking in tests.
var ExecCommand = exec.Command

// BuildImage writes the given Dockerfile content to a temporary file,
// then runs the docker build command to construct the image.
func BuildImage(dir string, dockerfileContent string, imageName string) error {
	tmpFile := filepath.Join(dir, fmt.Sprintf("%s.Dockerfile", imageName))
	err := os.WriteFile(tmpFile, []byte(dockerfileContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write Dockerfile: %w", err)
	}
	defer os.Remove(tmpFile)

	fmt.Printf("Building image %s...\n", imageName)
	cmd := ExecCommand("docker", "build", "-t", imageName, "-f", tmpFile, dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker build failed for %s: %w", imageName, err)
	}
	fmt.Printf("Successfully built %s\n", imageName)
	return nil
}

// RunTest executes the specified test command within the given directory
// by invoking a shell.
func RunTest(dir string, testCmd string) error {
	fmt.Printf("Running tests: %s...\n", testCmd)
	cmd := ExecCommand("sh", "-c", testCmd)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}
	fmt.Println("Tests passed successfully.")
	return nil
}

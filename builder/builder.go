// Package builder provides functionality to interact with the Docker daemon and shell
// to execute image builds and run project test suites.
package builder

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// CommandFunc defines the signature for creating an external *exec.Cmd subprocess.
// It accepts command (the executable name) and args (the command-line arguments).
type CommandFunc func(command string, args ...string) *exec.Cmd

// ExecCommand points to exec.Command to allow mocking subprocess creation in tests.
var ExecCommand CommandFunc = exec.Command

// Stdout is the default output writer for builder subprocesses and status messages.
var Stdout io.Writer = os.Stdout

// Stderr is the default error writer for builder subprocesses.
var Stderr io.Writer = os.Stderr

// createTempFile points to os.CreateTemp to allow deterministic testing of temporary file creation.
var createTempFile = os.CreateTemp

// writeTempFile writes data to a temporary file and closes it.
// It accepts f (the open temporary *os.File) and content (the Dockerfile string content).
func writeTempFile(f *os.File, content string) error {
	_, writeErr := f.WriteString(content)
	closeErr := f.Close()
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}

// BuildImage writes the given Dockerfile content to a temporary file inside dir,
// then runs the docker build command to construct the container image.
//
// Parameters:
//   - dir: The project directory path used as the Docker build context.
//   - dockerfileContent: The complete Dockerfile text to build.
//   - imageName: The tag name to assign to the built Docker image.
//
// Returns an error if creating/writing the temporary Dockerfile or running docker build fails.
func BuildImage(dir string, dockerfileContent string, imageName string) error {
	tmpFile, err := createTempFile(dir, fmt.Sprintf("%s-*.Dockerfile", imageName))
	if err != nil {
		return fmt.Errorf("failed to create temporary Dockerfile: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	if err := writeTempFile(tmpFile, dockerfileContent); err != nil {
		return fmt.Errorf("failed to write temporary Dockerfile: %w", err)
	}

	_, _ = fmt.Fprintln(Stdout, "Building image "+imageName+"...")
	cmd := ExecCommand("docker", "build", "-t", imageName, "-f", tmpPath, dir)
	cmd.Stdout = Stdout
	cmd.Stderr = Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker build failed for %s: %w", imageName, err)
	}
	_, _ = fmt.Fprintln(Stdout, "Successfully built "+imageName)
	return nil
}

// RunTest executes the specified test command within the given directory by invoking a POSIX shell.
//
// Parameters:
//   - dir: The working directory in which to run the test command.
//   - testCmd: The shell command string to execute.
//
// Returns an error if the test command exits with a non-zero status.
func RunTest(dir string, testCmd string) error {
	_, _ = fmt.Fprintln(Stdout, "Running tests: "+testCmd+"...")
	cmd := ExecCommand("sh", "-c", testCmd)
	cmd.Dir = dir
	cmd.Stdout = Stdout
	cmd.Stderr = Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}
	_, _ = fmt.Fprintln(Stdout, "Tests passed successfully.")
	return nil
}

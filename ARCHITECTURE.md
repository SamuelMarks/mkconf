# Architecture

This document describes the high-level architecture of `mkconf`.

## Core Components

The tool is segmented into distinct core domains orchestrated by the **SDK** (`sdk/sdk.go`), which in turn is consumed by the CLI entrypoint (`cmd/mkconf/main.go`):

### 1. SDK (`sdk/sdk.go`)
Provides a high-level wrapper API containing the exact same feature set as the CLI. Consumers can import the `mkconf/sdk` package to programmatically scan directories, generate configs, and build containers without using CLI logic.

### 2. Scanner (`scanner/scanner.go`)
Performs a heuristic scan of a provided directory path.
- **Manifest Detection**: It prioritizes scanning for manifest files (`go.mod`, `Cargo.toml`, `requirements.txt`, etc.).
- **Fallback Extension Scanning**: If a manifest is unavailable, it recursively traverses the root tree (via `filepath.Walk`) searching for primary extensions (`.rs`, `.go`, `.rb`, etc.) and bails eagerly on the first match.
- **Output**: Returns a `ProjectInfo` containing mapped lifecycle commands (`InstallCommand`, `BuildCommand`, `TestCommand`, `StartCommand`).

### 3. Generator (`generator/ast.go`, `generator/bazel.go`, `generator/docker_compose.go`, `generator/dockerfile.go`, `generator/makebat.go`, `generator/makefile.go`)
Constructs target files depending on user requests.
- **Dockerfiles**: Constructs standard, alpine, and multi-stage distroless `Dockerfile`s depending on user requests. Uses `github.com/moby/buildkit/frontend/dockerfile/parser.Node` to construct instructions.
- **New Formats**: Generates Docker Compose (`docker-compose.yml`), Makefiles (`Makefile` and `make.bat`), and Bazel (`BUILD`) files.

### 4. Builder (`builder/builder.go`)
Manages execution logic in external shells and daemons.
- **Testing**: Invokes the scanner's resolved `TestCommand` in an attached PTY/stdout prior to final generation.
- **Image Construction**: Dumps output strings to temporary `app-.*.Dockerfile*` files and constructs the image directly utilizing the host's `docker` socket CLI tool.

# Usage

`mkconf` leverages standard system interfaces to generate configuration files for your projects.

## Installation

Ensure you have Go 1.22+ installed.

```bash
git clone https://github.com/samuel/mkconf.git
cd mkconf
go build -o mkconf ./cmd/mkconf
```

## Running the Application

Provide the path to your source repository as the first argument:

```bash
./mkconf /path/to/my-repo
```

### Flow

1. The scanner will identify the target project structure.
2. The resolved test commands will be executed natively locally (e.g., `go test ./...` or `cargo test`).
3. By default, `docker-compose.yml`, `Makefile`, `make.bat`, Bazel `BUILD` files, and `Dockerfile`s will be generated in the target directory.
4. You can specify which formats to emit using the `--emit-*` flags.
5. If the `--build` flag is provided (and a Docker daemon is running), the tool will spawn three individual `docker build` processes:
   - `app-debian`
   - `app-alpine`
   - `app-distroless`

## CLI Reference

### `mkconf`

Generates build and environment configuration files.

**Usage:**
```bash
mkconf [repo_path] [flags]
```

**Flags:**
- `-o`, `--output string`: Output directory for generated files (defaults to repo_path)
- `--no-test`: Skip running the project's test suite
- `--dry-run`: Do not write files or build images, only print output
- `--emit-dockerfile`: Emit Dockerfile(s) and docker-compose.yml
- `--emit-bazel-build-file`: Emit Bazel BUILD file
- `--emit-makefile`: Emit Makefile and make.bat
- `--build`: Build Docker images after generating Dockerfiles (requires Docker)
- `-h`, `--help`: help for mkconf

## Supported Languages

`mkconf` supports 40 programming languages and frameworks. For a comprehensive list, including default build commands and base images, please refer to the **Supported Languages** section in the [README.md](README.md).

Some examples include:
- **Go**: Uses `golang:1.22` and statically compiles for `gcr.io/distroless/static`.
- **Rust**: Uses `rust:1` and multi-stages into `gcr.io/distroless/cc`.
- **Python**: Uses `python:3.11` and relies on `gcr.io/distroless/python3`.
- **Ruby**: Uses `ruby:3.3` and relies on `cgr.dev/chainguard/ruby`.
- **C/C++**: Relies on `gcc:13`, dynamically linking against `gcr.io/distroless/cc` runtime.

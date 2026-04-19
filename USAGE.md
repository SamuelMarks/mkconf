# Usage

`mkconf` leverages standard system interfaces and requires a running Docker daemon.

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
3. Three Dockerfiles will be generated and dumped to standard output.
4. The tool will spawn three individual `docker build` processes:
   - `app-debian`
   - `app-alpine`
   - `app-distroless`
5. Previews of `docker-compose.yml`, `Makefile`, `make.bat`, and Bazel `BUILD` files will be printed to standard output.

## Supported Languages

`mkconf` supports 40 programming languages and frameworks. For a comprehensive list, including default build commands and base images, please refer to the **Supported Languages** section in the [README.md](README.md).

Some examples include:
- **Go**: Uses `golang:1.22` and statically compiles for `gcr.io/distroless/static`.
- **Rust**: Uses `rust:1` and multi-stages into `gcr.io/distroless/cc`.
- **Python**: Uses `python:3.11` and relies on `gcr.io/distroless/python3`.
- **Ruby**: Uses `ruby:3.3` and relies on `cgr.dev/chainguard/ruby`.
- **C/C++**: Relies on `gcc:13`, dynamically linking against `gcr.io/distroless/cc` runtime.

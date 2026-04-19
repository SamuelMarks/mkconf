setup:
	@echo "Setting up pre-commit hooks..."
	git config core.hooksPath .githooks

.PHONY: help install_base install_deps build test run build_docker run_docker

help:
	@echo "  make install_base  - Install Go"
	@echo "  make install_deps  - Install dependencies"
	@echo "  make build         - Build the application"
	@echo "  make test          - Run tests locally"
	@echo "  make run           - Run the application"
	@echo "  make build_docker  - Build Docker images"
	@echo "  make run_docker    - Run Docker images"

install_base:
	@echo "Please install Go manually"

install_deps:
	@echo "No install dependencies command defined"

build:
	go build -o app

test:
	go test ./...

run:
	/app/app

build_docker:
	docker build -t app-debian -f debian.Dockerfile .
	docker build -t app-alpine -f alpine.Dockerfile .
	docker build -t app-distroless -f distroless.Dockerfile .

run_docker:
	docker run --rm -it app-debian

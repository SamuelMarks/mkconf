package generator

import (
	"mkconf/config"
	"mkconf/scanner"
	"strings"
	"testing"
)

func TestGenerateDockerfile(t *testing.T) {
	info := &scanner.ProjectInfo{
		Language:       "go",
		InstallCommand: "go mod download",
		BuildCommand:   "go build -o app",
		TestCommand:    "go test ./...",
		StartCommand:   []string{"/app/app"},
	}

	debian := GenerateDockerfile(info, "debian")
	if !strings.Contains(debian, "FROM golang:") {
		t.Error("Expected FROM golang in debian Dockerfile")
	}

	alpine := GenerateDockerfile(info, "alpine")
	if !strings.Contains(alpine, "FROM golang:") {
		t.Error("Expected FROM golang in alpine Dockerfile")
	}

	distroless := GenerateDockerfile(info, "distroless")
	if !strings.Contains(distroless, "FROM golang:") || !strings.Contains(distroless, "AS builder") {
		t.Error("Expected multi-stage build in distroless Dockerfile")
	}
	if !strings.Contains(distroless, "COPY --from=builder /app/app /app/app") {
		t.Error("Expected single binary COPY /app/app for Go distroless")
	}

	clojureInfo := &scanner.ProjectInfo{
		Language: "clojure",
	}
	distrolessClojure := GenerateDockerfile(clojureInfo, "distroless")
	if !strings.Contains(distrolessClojure, "COPY --from=builder /app /app") {
		t.Errorf("Expected full /app COPY for non-/app/app compiled language, got: %s", distrolessClojure)
	}
	if !strings.Contains(distrolessClojure, `ENTRYPOINT ["java","-jar","target/app.jar"]`) {
		t.Errorf("Expected default language StartCommand when info.StartCommand is empty, got: %s", distrolessClojure)
	}

	nilOutput := GenerateDockerfile(nil, "debian")
	if !strings.Contains(nilOutput, "FROM ubuntu:22.04") {
		t.Error("Expected fallback to ubuntu:22.04 for nil info")
	}

	unknownInfo := &scanner.ProjectInfo{
		Language: "unknown",
	}
	unknownOutput := GenerateDockerfile(unknownInfo, "debian")
	if !strings.Contains(unknownOutput, "FROM ubuntu:") {
		t.Error("Expected fallback to ubuntu for unknown language")
	}

	nodeInfo := &scanner.ProjectInfo{
		Language:       "nodejs",
		StartCommand:   []string{"npm", "start"},
		InstallCommand: "npm install",
		BuildCommand:   "npm run build",
	}
	distrolessNode := GenerateDockerfile(nodeInfo, "distroless")
	if !strings.Contains(distrolessNode, "CMD") {
		t.Error("Expected CMD for uncompiled distroless")
	}

	alpineNode := GenerateDockerfile(nodeInfo, "alpine")
	if !strings.Contains(alpineNode, "CMD") {
		t.Error("Expected CMD for uncompiled alpine")
	}
}

func TestGenerateDockerfileEmptyBase(t *testing.T) {
	info := &scanner.ProjectInfo{
		Language: "go",
	}
	output := GenerateDockerfile(info, "unknown_base")
	if !strings.Contains(output, "FROM ubuntu:22.04") {
		t.Error("Expected ubuntu:22.04 when base image is empty")
	}
}

func TestGenerateDockerfileAlpinePackages(t *testing.T) {
	info := &scanner.ProjectInfo{
		Language:     "c",
		StartCommand: []string{"/app/app"},
	}
	output := GenerateDockerfile(info, "alpine")
	if !strings.Contains(output, "apk add --no-cache") {
		t.Error("Expected apk add in alpine Dockerfile for C")
	}

	origLanguages := config.Languages
	t.Cleanup(func() {
		config.Languages = origLanguages
	})
	config.Languages = append(config.Languages, config.LanguageDef{
		Name:         "custom_alpine_builder",
		StartCommand: []string{"/app/app"},
		IsCompiled:   true,
		Docker: config.DockerDef{
			DebianImage:            "debian:bookworm-slim",
			AlpineImage:            "alpine:3.19",
			DistrolessBuilderImage: "alpine:3.19",
			DistrolessImage:        "gcr.io/distroless/static-debian12",
			AlpinePackages:         []string{"musl-dev"},
		},
	})
	distrolessAlpineBuilder := GenerateDockerfile(&scanner.ProjectInfo{Language: "custom_alpine_builder"}, "distroless")
	if !strings.Contains(distrolessAlpineBuilder, "RUN apk add --no-cache musl-dev") {
		t.Errorf("Expected apk add in alpine-based distroless builder stage, got: %s", distrolessAlpineBuilder)
	}
}

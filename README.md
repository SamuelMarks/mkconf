`mkconf`
========

[![Test Coverage](https://img.shields.io/badge/Test_Coverage-100.0%25-brightgreen.svg)]()
[![Doc Coverage](https://img.shields.io/badge/Doc_Coverage-100.0%25-brightgreen.svg)]()
[![CI](https://github.com/SamuelMarks/mkconf/actions/workflows/ci.yml/badge.svg)](https://github.com/SamuelMarks/mkconf/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-Apache--2.0%20OR%20MIT-blue.svg)](https://opensource.org/licenses/Apache-2.0)

`mkconf` is a powerful, automated CLI utility and Go SDK that seamlessly detects the language and framework of your repository to generate comprehensive build, deployment, and containerization configurations. Designed to be highly useful both with and without Docker, it can construct three types of optimized Docker images ([Debian](https://www.debian.org/)-based, [Alpine](https://alpinelinux.org/)-based, and [Distroless](https://github.com/GoogleContainerTools/distroless)) alongside Docker Compose environments, or work natively by scaffolding Makefiles and/or Bazel `BUILD` files. Built on an AST-powered generator, `mkconf` ensures reliability by automatically invoking your test suite and verifying builds before outputting anything.

## Features

- **Multi-language Detection**: Automatically heuristically detects Awk, Bash, Bun, C, C#, C++, Clojure, Crystal, D, Dart, Deno, Elixir, Erlang, F#, Fortran, Gleam, Go, Groovy, Haskell, Haxe, Java, Julia, Kotlin, Lua, Nim, Node.js, OCaml, Perl, PHP, PowerShell, Python, R, Racket, Ruby, Rust, Scala, Swift, Tcl, V, and Zig.
- **Three Target Formats**: Generates [Alpine](https://alpinelinux.org/), standard [Debian](https://www.debian.org/), and multi-stage [Distroless](https://github.com/GoogleContainerTools/distroless) builds.
- **AST Based Generation**: Powered by Moby's official `buildkit` Dockerfile parser.
- **Automated Verification**: Automatically invokes your application's test suite before generation and runs `docker build` immediately to verify builds.
- **Make**: Generates `Makefile` and `make.bat`.
- **Bazel**: Generates Bazel `BUILD` files.
- **Docker Compose**: Generates `docker-compose.yml`.


## Supported Languages

| Language | Primary Heuristic | Fallback Heuristic | Default Build Command | Base Images ([Debian](https://www.debian.org/) / [Alpine](https://alpinelinux.org/) / [Distroless](https://github.com/GoogleContainerTools/distroless)) |
|---|---|---|---|---|
| **[Awk](https://en.wikipedia.org/wiki/AWK)** | - | `.awk` files | N/A | `debian:bookworm-slim`, `alpine:3.19` / `base-debian12` |
| **[Bash](https://en.wikipedia.org/wiki/Bash_(Unix_shell))** | - | `.sh`, `.bash` files | N/A | `bash:5.2`, `bash:5.2-alpine` / `base-debian12` |
| **[Bun](https://bun.sh/)** | `bunfig.toml` | - | `bun install` | `oven/bun:debian`, `oven/bun:alpine` / `oven/bun:distroless` |
| **[C](https://en.wikipedia.org/wiki/C_(programming_language))** | `Makefile` | `.c` files | `make` | `gcc:13`, `alpine` / `cc-debian12` |
| **[C#](https://en.wikipedia.org/wiki/C_Sharp_(programming_language))** | - | `.cs` files | `dotnet build` | `dotnet/sdk:8.0` / `dotnet/aspnet:8.0` |
| **[C++](https://en.wikipedia.org/wiki/C%2B%2B)** | `CMakeLists.txt` | `.cpp`, `.cc` files | `cmake . && make` | `gcc:13`, `alpine` / `cc-debian12` |
| **[Clojure](https://en.wikipedia.org/wiki/Clojure)** | `project.clj` | `.clj`, `.cljs`, `.cljc` files | `lein uberjar` | `clojure:temurin-21-lein-jammy`, `clojure:temurin-21-lein-alpine` / `java21-debian12` |
| **[Crystal](https://en.wikipedia.org/wiki/Crystal_(programming_language))** | `shard.yml` | `.cr` files | `crystal build --release src/app.cr` | `crystallang/crystal:1.12`, `crystallang/crystal:1.12-alpine` / `cc-debian12` |
| **[D](https://en.wikipedia.org/wiki/D_(programming_language))** | `dub.json`, `dub.sdl` | `.d` files | `dub build --build=release` | `dlang2/dmd-ubuntu:2.107.1`, `alpine:3.19` / `cc-debian12` |
| **[Dart](https://en.wikipedia.org/wiki/Dart_(programming_language))** | `pubspec.yaml` | `.dart` files | `dart pub get` | `dart:stable` / `base-debian12` |
| **[Deno](https://deno.com/)** | `deno.json`, `deno.jsonc` | - | N/A | `denoland/deno:latest`, `denoland/deno:alpine` / `denoland/deno:distroless` |
| **[Elixir](https://en.wikipedia.org/wiki/Elixir_(programming_language))** | `mix.exs` | `.ex`, `.exs` files | `mix compile` | `elixir:1.16` / `chainguard/elixir` |
| **[Erlang](https://en.wikipedia.org/wiki/Erlang_(programming_language))** | `rebar.config` | `.erl` files | `rebar3 compile` | `erlang:26`, `erlang:26-alpine` / `base-debian12` |
| **[F#](https://en.wikipedia.org/wiki/F_Sharp_(programming_language))** | `*.fsproj` | `.fs` files | `dotnet build` | `mcr.microsoft.com/dotnet/sdk:8.0`, `mcr.microsoft.com/dotnet/sdk:8.0-alpine` / `mcr.microsoft.com/dotnet/aspnet:8.0` |
| **[Fortran](https://en.wikipedia.org/wiki/Fortran)** | - | `.f90`, `.f`, `.f03` files | `gfortran -o app *.f90` | `gcc:13-bookworm`, `alpine:3.19` / `cc-debian12` |
| **[Gleam](https://en.wikipedia.org/wiki/Gleam_(programming_language))** | `gleam.toml` | `.gleam` files | `gleam build` | `ghcr.io/gleam-lang/gleam:v1.0.0`, `ghcr.io/gleam-lang/gleam:v1.0.0-erlang-alpine` / `cc-debian12` |
| **[Go](https://en.wikipedia.org/wiki/Go_(programming_language))** | `go.mod` | `.go` files | `go build -o app` | `golang:1.22` / `static-debian12` |
| **[Groovy](https://en.wikipedia.org/wiki/Apache_Groovy)** | - | `.groovy` files | N/A | `groovy:4-jdk21`, `groovy:4-jdk21-alpine` / `java21-debian12` |
| **[Haskell](https://en.wikipedia.org/wiki/Haskell)** | `stack.yaml` | `.hs` files | `stack build` | `haskell:9`, `haskell:9-alpine` / `cc-debian12` |
| **[Haxe](https://en.wikipedia.org/wiki/Haxe)** | `build.hxml` | `.hx` files | `haxe build.hxml` | `haxe:4.3.2`, `haxe:4.3.2-alpine` / `nodejs20-debian12` |
| **[Java](https://en.wikipedia.org/wiki/Java_(programming_language))** | `pom.xml`, `build.gradle` | `.java` files | `mvn clean package` / `gradle build` | `eclipse-temurin:21-jdk` / `java21-debian12` |
| **[Julia](https://en.wikipedia.org/wiki/Julia_(programming_language))** | `Project.toml` | `.jl` files | N/A | `julia:1.10-bookworm`, `julia:1.10-alpine` / `chainguard/julia` |
| **[Kotlin](https://en.wikipedia.org/wiki/Kotlin_(programming_language))** | - | `.kt`, `.kts` files | `kotlinc *.kt -include-runtime -d app.jar` | `eclipse-temurin:21-jdk` / `java21-debian12` |
| **[Lua](https://en.wikipedia.org/wiki/Lua_(programming_language))** | - | `.lua` files | N/A | `lua:5.4` / `chainguard/lua` |
| **[Nim](https://en.wikipedia.org/wiki/Nim_(programming_language))** | `*.nimble` | `.nim`, `.nimble` files | `nimble build -y -d:release` | `nimlang/nim:2.0.4`, `nimlang/nim:2.0.4-alpine` / `cc-debian12` |
| **[Node.js](https://en.wikipedia.org/wiki/Node.js)** | `package.json` | `.js`, `.ts` files | `npm install` | `node:20` / `nodejs20-debian12` |
| **[OCaml](https://en.wikipedia.org/wiki/OCaml)** | `dune-project` | `.ml` files | `dune build` | `ocaml/opam:debian`, `ocaml/opam:alpine` / `cc-debian12` |
| **[Perl](https://en.wikipedia.org/wiki/Perl)** | - | `.pl`, `.pm` files | N/A | `perl:5.38` / `base-debian12` |
| **[PHP](https://en.wikipedia.org/wiki/PHP)** | `composer.json` | `.php` files | `composer install` | `php:8.2-cli` / `chainguard/php` |
| **[PowerShell](https://en.wikipedia.org/wiki/PowerShell)** | - | `.ps1`, `.psm1`, `.psd1` files | N/A | `mcr.microsoft.com/powershell:lts-debian-11`, `mcr.microsoft.com/powershell:lts-alpine-3.17` / `base-debian12` |
| **[Python](https://en.wikipedia.org/wiki/Python_(programming_language))** | `requirements.txt` | `.py` files | `pip install -r requirements.txt` | `python:3.11` / `python3-debian12` |
| **[R](https://en.wikipedia.org/wiki/R_(programming_language))** | - | `.r` files | N/A | `r-base:4.3.3`, `alpine:3.19` / `base-debian12` |
| **[Racket](https://racket-lang.org/)** | `info.rkt` | `.rkt` files | N/A | `racket/racket:8.11-full`, `alpine:3.19` / `base-debian12` |
| **[Ruby](https://en.wikipedia.org/wiki/Ruby_(programming_language))** | `Gemfile` | `.rb` files | `bundle install` | `ruby:3.3` / `chainguard/ruby` |
| **[Rust](https://en.wikipedia.org/wiki/Rust_(programming_language))** | `Cargo.toml` | `.rs` files | `cargo build --release` | `rust:1` / `cc-debian12` |
| **[Scala](https://en.wikipedia.org/wiki/Scala_(programming_language))** | `build.sbt` | `.scala` files | `sbt compile` | `sbtscala/scala-sbt:eclipse-temurin-jammy-21.0.2_13_1.9.9_3.4.1`, `sbtscala/scala-sbt:eclipse-temurin-alpine-21.0.2_13_1.9.9_3.4.1` / `java21-debian12` |
| **[Swift](https://en.wikipedia.org/wiki/Swift_(programming_language))** | - | `.swift` files | `swift build -c release` | `swift:5.10` / `cc-debian12` |
| **[Tcl](https://www.tcl.tk/)** | - | `.tcl` files | N/A | `debian:bookworm-slim`, `alpine:3.19` / `base-debian12` |
| **[V](https://en.wikipedia.org/wiki/V_(programming_language))** | `v.mod` | `.v` files | `v .` | `thevlang/vlang:alpine`, `alpine:3.19` / `static-debian12` |
| **[Zig](https://en.wikipedia.org/wiki/Zig_(programming_language))** | `build.zig` | `.zig` files | `zig build -Doptimize=ReleaseSafe` | `ziglang/zig:0.12.0`, `alpine:3.21` / `cc-debian12` |



## Getting Started

See [USAGE.md](USAGE.md) for installation and invocation commands.
See [ARCHITECTURE.md](ARCHITECTURE.md) to understand the component boundaries and implementation details.

## Development

We enforce 100% test coverage and 100% doc comment coverage. A pre-commit hook automatically updates the coverage shields.

---

## License

Licensed under either of

- Apache License, Version 2.0 ([LICENSE-APACHE](LICENSE-APACHE) or <https://www.apache.org/licenses/LICENSE-2.0>)
- MIT license ([LICENSE-MIT](LICENSE-MIT) or <https://opensource.org/licenses/MIT>)

at your option.

### Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in the work by you, as defined in the Apache-2.0 license, shall be
dual licensed as above, without any additional terms or conditions.

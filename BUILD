# Bazel BUILD file

genrule(
    name = "install_base",
    outs = ["install_base.out"],
    cmd = "echo 'Please install Go manually' > $@",
)

genrule(
    name = "install_deps",
    outs = ["install_deps.out"],
    cmd = "echo "No install dependencies command defined" > $@",
)

genrule(
    name = "build",
    outs = ["build.out"],
    cmd = "go build -o app > $@",
)

genrule(
    name = "test",
    outs = ["test.out"],
    cmd = "go test ./... > $@",
)

genrule(
    name = "run",
    outs = ["run.out"],
    cmd = "/app/app > $@",
)

genrule(
    name = "build_docker",
    outs = ["build_docker.out"],
    cmd = "docker build -t app-debian -f debian.Dockerfile . && docker build -t app-alpine -f alpine.Dockerfile . && docker build -t app-distroless -f distroless.Dockerfile . > $@",
)

genrule(
    name = "run_docker",
    outs = ["run_docker.out"],
    cmd = "docker run --rm -it app-debian > $@",
)

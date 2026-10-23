# Bazel BUILD file

genrule(
    name = "install_base",
    srcs = glob(["**/*"]),
    outs = ["install_base.out"],
    cmd = "echo 'Please install Go manually' > $@",
)

genrule(
    name = "install_deps",
    srcs = glob(["**/*"]),
    outs = ["install_deps.out"],
    cmd = "echo 'No install dependencies command defined' > $@",
)

genrule(
    name = "build",
    srcs = glob(["**/*"]),
    outs = ["build.out"],
    cmd = "go build -o app > $@",
)

genrule(
    name = "test",
    srcs = glob(["**/*"]),
    outs = ["test.out"],
    cmd = "go test ./... > $@",
)

genrule(
    name = "run",
    srcs = glob(["**/*"]),
    outs = ["run.out"],
    cmd = "/app/app > $@",
)

genrule(
    name = "build_docker",
    srcs = glob(["**/*"]),
    outs = ["build_docker.out"],
    cmd = "(docker build -t app-debian -f debian.Dockerfile . && docker build -t app-alpine -f alpine.Dockerfile . && docker build -t app-distroless -f distroless.Dockerfile .) > $@",
)

genrule(
    name = "run_docker",
    srcs = glob(["**/*"]),
    outs = ["run_docker.out"],
    cmd = "docker run --rm app-debian > $@",
)

@echo off
IF "%1"=="" GOTO help
GOTO %1

:help
echo   make.bat install_base  - Install Go
echo   make.bat install_deps  - Install dependencies
echo   make.bat build         - Build the application
echo   make.bat test          - Run tests locally
echo   make.bat run           - Run the application
echo   make.bat build_docker  - Build Docker images
echo   make.bat run_docker    - Run Docker images
GOTO :EOF

:install_base
echo Please install Go manually
GOTO :EOF

:install_deps
echo No install dependencies command defined
GOTO :EOF

:build
go build -o app
GOTO :EOF

:test
go test ./...
GOTO :EOF

:run
/app/app
GOTO :EOF

:build_docker
docker build -t app-debian -f debian.Dockerfile .
docker build -t app-alpine -f alpine.Dockerfile .
docker build -t app-distroless -f distroless.Dockerfile .
GOTO :EOF

:run_docker
docker run --rm -it app-debian
GOTO :EOF

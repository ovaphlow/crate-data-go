# 设置容器名称和镜像名称
$containerName = "crate-data-build"
$imageName = "docker.io/library/golang:1.24-alpine"

# 确保build-target目录存在
if (-not (Test-Path "build-target")) {
    New-Item -ItemType Directory -Path "build-target"
}

# 停止正在运行的容器
Write-Host "Stopping running container if exists..." -ForegroundColor Yellow
podman stop $containerName 2>$null

# 运行容器来同时构建Linux和Windows二进制文件
Write-Host "Building binaries for Linux and Windows..." -ForegroundColor Green
podman run --rm --name $containerName `
    -v "${PWD}:/app" `
    -e GOPROXY=https://goproxy.cn,direct `
    -e HTTP_PROXY="" `
    -e HTTPS_PROXY="" `
    -e http_proxy="" `
    -e https_proxy="" `
    $imageName /bin/sh -c "sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && apk update && apk add --no-cache mingw-w64-gcc build-base sqlite-dev && cd /app && go mod download && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags='-s -w' -o build-target/crate-api-data cmd/main.go && CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -trimpath -ldflags='-s -w' -o build-target/crate-api-data.exe cmd/main.go"

# 复制配置文件
Write-Host "Copying configuration files..." -ForegroundColor Green
if (Test-Path ".env") {
    Copy-Item -Path ".env" -Destination "build-target/"
}

Write-Host "Build completed!" -ForegroundColor Cyan
Write-Host "Output files are in the ./build-target directory"
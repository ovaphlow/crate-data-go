$ErrorActionPreference = "Stop"

# Set container name and image name
$containerName = "crate-data-dev"
$imageName = "docker.io/library/golang:1.24-alpine"

Write-Host "Starting development container setup..." -ForegroundColor Cyan
Write-Host "Container name: $containerName" -ForegroundColor Yellow
Write-Host "Base image: $imageName" -ForegroundColor Yellow

Write-Host "`nLaunching container with Chinese mirrors and development configuration..." -ForegroundColor Green
# Run development container with host network and auto-remove when stopped
podman run -it --rm `
    --name $containerName `
    --network host `
    -v ${PWD}:/app `
    -e GOPROXY=https://goproxy.cn,direct `
    -e HTTP_PROXY="" `
    -e HTTPS_PROXY="" `
    -e http_proxy="" `
    -e https_proxy="" `
    $imageName /bin/sh -c "echo 'Configuring Alpine mirrors...' && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && echo 'Updating package list...' && apk update && echo 'Changing to application directory...' && cd /app && echo 'Tidying Go dependencies...' && go mod tidy && echo 'Development environment ready!' && /bin/sh"

Write-Host "`nContainer stopped and removed." -ForegroundColor Cyan

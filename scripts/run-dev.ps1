# 定义容器名称和镜像标签
$containerName = "crate-data"
$imageTag = "localhost/crate-data:latest"

# 检查并停止所有现有的同名容器
Write-Host "Checking for existing containers..." -ForegroundColor Yellow
$existingContainers = podman ps -a --filter name=$containerName --format "{{.ID}}"
if ($existingContainers) {
    Write-Host "Found existing containers. Stopping and removing..." -ForegroundColor Yellow
    $existingContainers | ForEach-Object {
        podman stop $_ 2>$null
        podman rm $_ 2>$null
    }
}

# 删除旧的镜像（如果存在）
podman rmi $imageTag 2>$null

# 构建新的镜像
Write-Host "Building new image..." -ForegroundColor Green
podman build --no-cache -t $imageTag -f Dockerfile.dev .

# 运行新的容器
Write-Host "Starting container in background..." -ForegroundColor Green
podman run -d --name $containerName -p 8421:8421 -v ${PWD}:/app $imageTag

# 显示容器状态
Write-Host "Container status:" -ForegroundColor Green
podman ps --filter name=$containerName
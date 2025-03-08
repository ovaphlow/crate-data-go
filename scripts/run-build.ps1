$ErrorActionPreference = "Stop"

# 设置变量
$ImageName = "crate-data-build"
$ContainerName = "crate-data-build"
$Version = "latest"

Write-Host "清理旧的容器和镜像..." -ForegroundColor Yellow

# 检查并删除旧容器
$existingContainer = podman ps -a --filter "name=$ContainerName" --format "{{.Names}}"
if ($existingContainer) {
    Write-Host "删除已存在的构建容器..." -ForegroundColor Yellow
    podman rm $ContainerName
}

# 检查并删除旧镜像
$existingImage = podman images --filter "reference=${ImageName}:${Version}" --format "{{.Repository}}"
if ($existingImage) {
    Write-Host "删除已存在的构建镜像..." -ForegroundColor Yellow
    podman rmi ${ImageName}:${Version} -f
}

Write-Host "开始构建程序..." -ForegroundColor Green

# 构建镜像
podman build -t ${ImageName}:${Version} .

Write-Host "创建临时容器..." -ForegroundColor Green

# 创建容器（不运行）用于复制文件
podman create --name $ContainerName ${ImageName}:${Version}

# 确保build-target目录存在
New-Item -ItemType Directory -Force -Path "build-target"

Write-Host "从容器中复制编译好的程序..." -ForegroundColor Green

# 从容器中复制编译好的程序到本地build-target目录
podman cp "${ContainerName}:/app/build-target/crate-api-data" "build-target/"
podman cp "${ContainerName}:/app/build-target/crate-api-data.exe" "build-target/"
Copy-Item ".env" "build-target/" -ErrorAction SilentlyContinue

Write-Host "清理临时容器和镜像..." -ForegroundColor Yellow

# 清理容器和镜像
podman rm $ContainerName
podman rmi ${ImageName}:${Version}

Write-Host "构建完成！程序已保存到 build-target 目录" -ForegroundColor Green
# Go 开发环境镜像制作步骤

## 1. 构建基础镜像

```bash
# 使用 Dockerfile.dev 构建基础镜像
podman build -t crate-data-env -f Dockerfile.dev .
```

## 2. 运行容器并安装依赖

```bash
# 运行容器（不使用 --rm 参数，以便保存更改）
podman run -it --name crate-data-dev crate-data-env

# 在容器内执行以下命令安装依赖
go mod download
go mod tidy
```

## 3. 保存带依赖的镜像

打开新的终端窗口，执行：

```bash
# 将容器保存为新镜像
podman commit crate-data-dev crate-data-env:with-deps
```

## 4. 开发时使用

开发时可以使用以下两种方式之一运行容器：

### 方式一：使用端口映射

```bash
# 使用端口映射访问应用
podman run -it -v .:/app -p 8421:8421 crate-data-env:with-deps
```

### 方式二：使用主机网络

```bash
# 直接使用主机网络
podman run -it -v .:/app --network host crate-data-env:with-deps
```

## 注意事项

1. 容器内已配置了国内镜像源：
   - Go 代理：`https://goproxy.cn`
   - Alpine 镜像源：USTC 镜像

2. 容器内已安装基础开发工具：
   - git
   - gcc
   - musl-dev

3. 环境变量设置：
   - `CGO_ENABLED=1`
   - 已清除所有代理设置
FROM docker.io/library/golang:1.24-alpine AS builder

# 使用阿里云的 Alpine 镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装编译所需的工具，包括 Windows 交叉编译工具
RUN apk add --no-cache gcc musl-dev mingw-w64-gcc g++

WORKDIR /app
COPY . .

# 设置 Go 代理
ENV GOPROXY=https://goproxy.cn,direct

RUN mkdir -p build-target
RUN go mod download
# 构建 Linux 版本
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o build-target/crate-api-data cmd/main.go
# 构建 Windows 版本（与 meson 配置保持一致）
RUN CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -trimpath -ldflags="-s -w" -o build-target/crate-api-data.exe cmd/main.go

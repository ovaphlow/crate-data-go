# 第一阶段：使用 Go 官方镜像编译
FROM docker.1ms.run/library/golang:1.23-bookworm AS builder
WORKDIR /app
COPY . .
RUN go build -o crate-api-data ./cmd/main.go

# 第二阶段：生成最终镜像
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/crate-api-data /app/crate-api-data
COPY .env* ./
RUN mkdir -p /app/logs
EXPOSE 8421
CMD ["/app/crate-api-data"]

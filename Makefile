OUTPUT_DIR := ./target

# 同步依赖
tidy:
	go mod tidy

# 构建项目
build-linux:
	mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=1 go build -ldflags "-s -w" -trimpath -o $(OUTPUT_DIR)/crate-api-data cmd/main.go
	cp .env $(OUTPUT_DIR)/

# 交叉编译到 Windows
build-windows:
	mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o $(OUTPUT_DIR)/crate-api-data.exe cmd/main.go
	cp .env $(OUTPUT_DIR)/

# 同时编译 Linux 和 Windows
build: clean build-linux build-windows

# 清理构建文件
clean:
	rm -rf $(OUTPUT_DIR)

.PHONY: tidy build build-windows build-linux clean

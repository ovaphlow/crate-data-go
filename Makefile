# 设置库路径
LIB_PATH := ./lib
OUTPUT_DIR := ./target

# 同步依赖
tidy:
	go mod tidy

# 构建项目
build-linux:
	mkdir -p $(OUTPUT_DIR)
	# export LD_LIBRARY_PATH=$(LIB_PATH):$$LD_LIBRARY_PATH && go build -o $(OUTPUT_DIR)/crate-api-data
	go build -o $(OUTPUT_DIR)/crate-api-data cmd/main.go
	# cp $(LIB_PATH)/* $(OUTPUT_DIR)
	cp .env $(OUTPUT_DIR)/
	# cp script/run-linux.sh $(OUTPUT_DIR)/

# 交叉编译到 Windows
build-windows:
	mkdir -p $(OUTPUT_DIR)
	# CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 \
	go build -o $(OUTPUT_DIR)/crate-api-data.exe cmd/main.go
	# cp $(LIB_PATH)/* $(OUTPUT_DIR)
	cp .env $(OUTPUT_DIR)/
	# cp script/run-windows.cmd $(OUTPUT_DIR)/

# 运行Go程序
run-go:
	# export LD_LIBRARY_PATH=$(LIB_PATH):$$LD_LIBRARY_PATH && go run cmd/main.go
	go run cmd/main.go

# 同时编译 Linux 和 Windows
build: clean build-linux build-windows

# 清理构建文件
clean:
	rm -rf $(OUTPUT_DIR)

# 在工程根目录执行utility目录的测试
test-utility:
	# export LD_LIBRARY_PATH=.$(LIB_PATH):$$LD_LIBRARY_PATH && go test ./utility -v
	go test ./utility -v

.PHONY: tidy build build-windows build-linux clean test-utility run-go

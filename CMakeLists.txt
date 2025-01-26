cmake_minimum_required(VERSION 3.10)

# 设置项目名称
project(CrateDataGo)

# 设置 Go 编译器
find_program(GO_EXECUTABLE NAMES go REQUIRED)

# 设置输出目录
set(OUTPUT_DIR "${CMAKE_BINARY_DIR}/target")

# 同步依赖
add_custom_target(tidy
    COMMAND ${GO_EXECUTABLE} mod tidy
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
)

# 构建项目
add_custom_target(build-linux
    COMMAND ${CMAKE_COMMAND} -E make_directory ${OUTPUT_DIR}
    COMMAND ${GO_EXECUTABLE} build -ldflags="-s -w" -trimpath -o ${OUTPUT_DIR}/crate-api-data cmd/main.go
    COMMAND ${CMAKE_COMMAND} -E copy ${CMAKE_SOURCE_DIR}/.env ${OUTPUT_DIR}/
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
)

# 交叉编译到 Windows
add_custom_target(build-windows
    COMMAND ${CMAKE_COMMAND} -E make_directory ${OUTPUT_DIR}
    COMMAND ${GO_EXECUTABLE} build -ldflags="-s -w" -trimpath -o ${OUTPUT_DIR}/crate-api-data.exe cmd/main.go
    COMMAND ${CMAKE_COMMAND} -E copy ${CMAKE_SOURCE_DIR}/.env ${OUTPUT_DIR}/
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    ENV GOOS=windows GOARCH=amd64
)

# 同时编译 Linux 和 Windows
add_custom_target(build ALL
    DEPENDS clean build-linux build-windows
)

# 清理构建文件
add_custom_target(clean
    COMMAND ${CMAKE_COMMAND} -E remove_directory ${OUTPUT_DIR}
)

cmake_minimum_required(VERSION 3.10)

# 设置项目名称
project(CrateDataGo)

# 设置 Go 编译器
find_program(GO_EXECUTABLE NAMES go REQUIRED)

# 设置输出目录
set(OUTPUT_DIR "${CMAKE_BINARY_DIR}/target")

# 设置编译选项
if (WIN32)
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -pthread")
    set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -pthread")
endif()

# 同步依赖
add_custom_target(tidy
    COMMAND ${GO_EXECUTABLE} mod tidy
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
)

# 构建项目
add_custom_target(build-linux
    COMMAND ${CMAKE_COMMAND} -E make_directory ${OUTPUT_DIR}
    COMMAND CGO_ENABLED=1 ${GO_EXECUTABLE} build -ldflags "-s -w" -trimpath -o ${OUTPUT_DIR}/crate-api-data cmd/main.go
    COMMAND ${CMAKE_COMMAND} -E copy ${CMAKE_SOURCE_DIR}/.env ${OUTPUT_DIR}/
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
)

# 交叉编译到 Windows (MySQL)
add_custom_target(build-windows-mysql
    COMMAND ${CMAKE_COMMAND} -E make_directory ${OUTPUT_DIR}
    COMMAND CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ ${GO_EXECUTABLE} build -ldflags "-s -w" -trimpath -o ${OUTPUT_DIR}/crate-api-data-mysql.exe cmd/service-mysql/main.go
    COMMAND ${CMAKE_COMMAND} -E copy ${CMAKE_SOURCE_DIR}/.env ${OUTPUT_DIR}/
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
)

# 同时编译 Linux 和 Windows
add_custom_target(build ALL
    DEPENDS clean_target build-linux build-windows-mysql
)

# 清理构建文件
add_custom_target(clean_target
    COMMAND ${CMAKE_COMMAND} -E remove_directory ${OUTPUT_DIR}
)

# CRATE DATA

## 初始化构建目录

在项目根目录下创建并初始化 `build` 目录：

```shell
mkdir -p build
cd build
cmake ..
```

## 同步依赖

```shell
cmake --build . --target tidy
```

## 构建

### 构建 Linux 版本

```shell
cmake --build . --target build-linux
```

### 交叉编译到 Windows (MySQL)

```shell
cmake --build . --target build-windows-mysql
```

### 同时编译 Linux 和 Windows

```shell
cmake --build . --target build
```

## 清理构建文件

```shell
cmake --build . --target clean_target
```

## 在 Windows 系统上使用 CMake

在 Windows 系统上使用 CMake，您需要安装一个兼容的工具，例如 Git Bash 或 Cygwin。安装完成后，您可以在这些终端中运行 `cmake` 和 `cmake --build` 命令。
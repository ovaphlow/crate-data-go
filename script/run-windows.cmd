@echo off

set LD_LIBRARY_PATH=%OUTPUT_DIR%;%LD_LIBRARY_PATH%
crate-api-data.exe

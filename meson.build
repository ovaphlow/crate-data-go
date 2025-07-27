project('crate-data-go', 'c', 'cpp',
  version: '0.1.0',
  default_options: ['warning_level=3']
)

go = find_program('go', required: true)
output_dir = join_paths(meson.current_build_dir(), 'target')
source_dir = meson.current_source_dir()

# Custom targets
run_command('mkdir', '-p', output_dir, check: true)

# Clean target (renamed to clean_target)
clean_target = custom_target('clean_target',
  output: 'clean_target',
  command: ['rm', '-rf', output_dir],
  build_by_default: false,
  build_always_stale: true
)

# Linux build target
linux_build = custom_target('linux_build',
  output: 'crate-api-data',
  command: [
    'sh', '-c',
    'cd "' + source_dir + '" && GOPROXY=https://goproxy.cn,direct ' + go.path() + ' build -ldflags "-s -w" -trimpath -o "' + output_dir + '/crate-api-data" cmd/main.go && cp .env "' + output_dir + '/"'
  ],
  depends: clean_target,
  build_by_default: false,
  build_always_stale: true
)

# Windows build target
windows_build = custom_target('windows_build',
  output: 'crate-api-data.exe',
  command: [
    'sh', '-c',
    'cd "' + source_dir + '" && GOPROXY=https://goproxy.cn,direct GOOS=windows GOARCH=amd64 ' + go.path() + ' build -ldflags "-s -w" -trimpath -o "' + output_dir + '/crate-api-data.exe" cmd/main.go'
  ],
  depends: clean_target,
  build_by_default: false,
  build_always_stale: true
)

# Build target for both platforms
build = custom_target('build_all',
  output: ['crate-api-data_linux', 'crate-api-data_win.exe'],
  command: ['true'],
  depends: [linux_build, windows_build],
  build_by_default: true,
  build_always_stale: true
)
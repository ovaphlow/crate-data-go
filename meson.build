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
linux_build = custom_target('linux-build',
  output: 'crate-api-data',
  command: [
    'sh', '-c',
    'cd "' + source_dir + '" && CGO_ENABLED=1 ' + go.path() + ' build -ldflags "-s -w" -trimpath -o "' + output_dir + '/crate-api-data" cmd/main.go && cp .env "' + output_dir + '/"'
  ],
  depends: clean_target,
  build_by_default: true,
  build_always_stale: true
)

# Windows build target
windows_build = custom_target('windows-build',
  output: 'crate-api-data.exe',
  command: [
    'sh', '-c',
    'cd "' + source_dir + '" && CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ ' + go.path() + ' build -ldflags "-s -w" -trimpath -o "' + output_dir + '/crate-api-data.exe" cmd/main.go && cp .env "' + output_dir + '/"'
  ],
  depends: clean_target,
  build_by_default: false,
  build_always_stale: true
)
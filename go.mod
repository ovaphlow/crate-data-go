module ovaphlow.com/crate/data

go 1.24.3

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.24
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

// Set GOPROXY for this module
// replace command: go env -w GOPROXY=https://goproxy.cn,direct

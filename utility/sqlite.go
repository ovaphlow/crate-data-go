package utility

import (
	"database/sql"
	"errors"
	"os"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

var SQLite *sql.DB

func InitSQLite() {
	err := godotenv.Load()
	if err != nil {
		ZapLogger.Warn("环境变量未设置 SQLite", zap.Error(err))
	}

	dsn := os.Getenv("SQLITE_DATABASE")
	if dsn == "" {
		ZapLogger.Fatal("环境变量未设置 SQLite", zap.Error(errors.New("环境变量未设置 SQLite")))
	}

	SQLite, err = sql.Open("sqlite3", dsn+"?_journal_mode=WAL&_cache=shared&_synchronous=NORMAL&_temp_store=MEMORY&_auto_vacuum=INCREMENTAL")
	if err != nil {
		ZapLogger.Fatal("Failed to open SQLite connection", zap.Error(err))
	}

	if err := SQLite.Ping(); err != nil {
		ZapLogger.Fatal("连接数据库失败 SQLite", zap.Error(err))
	}

	// 设置连接池
	numCPU := runtime.NumCPU()
	SQLite.SetMaxIdleConns(1)                   // 设置最大空闲连接数
	SQLite.SetMaxOpenConns(numCPU*2 + 1)        // 设置最大打开连接数
	SQLite.SetConnMaxIdleTime(15 * time.Minute) // 设置最大空闲时间

	ZapLogger.Info("连接数据库成功 SQLite")
}

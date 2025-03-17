package utility

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var SQLite *sql.DB

func InitSQLite() {
	err := godotenv.Load()
	if err != nil {
		log.Println("环境变量未设置 SQLite")
		log.Fatal(err.Error())
	}

	dsn := os.Getenv("SQLITE_DATABASE")
	if dsn == "" {
		log.Fatal(errors.New("环境变量未设置 SQLite"))
	}

	SQLite, err = sql.Open("sqlite3", dsn+"?_journal_mode=WAL&_cache=shared&_synchronous=NORMAL&_temp_store=MEMORY&_auto_vacuum=INCREMENTAL")
	if err != nil {
		log.Fatal(err)
	}

	if err := SQLite.Ping(); err != nil {
		log.Println("连接数据库失败 SQLite")
		log.Fatal(err.Error())
	}

	// 设置连接池
	numCPU := runtime.NumCPU()
	SQLite.SetMaxIdleConns(1)                   // 设置最大空闲连接数
	SQLite.SetMaxOpenConns(numCPU*2 + 1)        // 设置最大打开连接数
	SQLite.SetConnMaxIdleTime(15 * time.Minute) // 设置最大空闲时间

	log.Println("连接数据库成功 SQLite")
}

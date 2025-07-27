package utility

import (
	"database/sql"
	"fmt"
	"runtime"
	"time"

	_ "time/tzdata"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

var MySQL *sql.DB

func init() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		ZapLogger.Warn("Failed to load Asia/Shanghai timezone", zap.Error(err))
	} else {
		time.Local = loc
	}
}

func InitMySQL(user, password, host, port, database string) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&charset=utf8mb4&timeout=5s&readTimeout=5s&writeTimeout=5s",
		user,
		password,
		host,
		port,
		database,
	)
	var err error
	MySQL, err = sql.Open("mysql", dsn)
	if err != nil {
		ZapLogger.Fatal(err.Error())
	}
	MySQL.SetConnMaxLifetime(time.Minute * 3)
	cpuCount := runtime.NumCPU()
	MySQL.SetMaxOpenConns(cpuCount*2 + 1)
	// MySQL.SetMaxIdleConns(cpuCount*2 + 1)
	MySQL.SetMaxIdleConns(0)
	MySQL.SetConnMaxLifetime(time.Second * 30)
	if err = MySQL.Ping(); err != nil {
		ZapLogger.Fatal("连接数据库失败 MySQL", zap.Error(err))
	}
	ZapLogger.Info("连接数据库成功 MySQL")
}

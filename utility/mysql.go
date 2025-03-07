package utility

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "time/tzdata" // 引入时区数据
)

var MySQL *sql.DB

func init() {
	// 确保time包加载了时区数据
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("Warning: Failed to load Asia/Shanghai timezone: %v", err)
	} else {
		time.Local = loc
	}
}

func InitMySQL(user, password, host, port, database string) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&time_zone=%%27%%2B08:00%%27",
		user,
		password,
		host,
		port,
		database,
	)
	var err error
	MySQL, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err.Error())
	}
	MySQL.SetConnMaxLifetime(time.Minute * 3)
	cpuCount := runtime.NumCPU()
	MySQL.SetMaxOpenConns(cpuCount*2 + 1)
	// MySQL.SetMaxIdleConns(cpuCount*2 + 1)
	MySQL.SetMaxIdleConns(0)
	MySQL.SetConnMaxLifetime(time.Second * 30)
	if err = MySQL.Ping(); err != nil {
		log.Println("连接数据库失败 MySQL")
		log.Fatal(err.Error())
	}
	log.Println("连接数据库成功 MySQL")
}

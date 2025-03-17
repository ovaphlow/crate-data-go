package main

// 导入必要的包
import (
	"log"
	"net/http"
	"os"

	"ovaphlow.com/crate/data/middleware"
	"ovaphlow.com/crate/data/repository"
	"ovaphlow.com/crate/data/router"
	"ovaphlow.com/crate/data/service"
	"ovaphlow.com/crate/data/utility"

	"github.com/joho/godotenv"
)

func init() {
	// 加载环境变量文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 初始化结构化日志
	utility.InitSlog()

	// 初始化 PostgreSQL 数据库
	postgres_enabled := os.Getenv("POSTGRES_ENABLED")
	if postgres_enabled == "true" || postgres_enabled == "1" {
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		host := os.Getenv("POSTGRES_HOST")
		port := os.Getenv("POSTGRES_PORT")
		database := os.Getenv("POSTGRES_DATABASE")
		utility.InitPostgres(user, password, host, port, database)
	}

	// 初始化 MySQL 数据库
	mysql_enabled := os.Getenv("MYSQL_ENABLED")
	if mysql_enabled == "true" || mysql_enabled == "1" {
		user := os.Getenv("MYSQL_USER")
		password := os.Getenv("MYSQL_PASSWORD")
		host := os.Getenv("MYSQL_HOST")
		port := os.Getenv("MYSQL_PORT")
		database := os.Getenv("MYSQL_DATABASE")
		utility.InitMySQL(user, password, host, port, database)
	}

	// 初始化 SQLite 数据库
	sqlite_enabled := os.Getenv("SQLITE_ENABLED")
	if sqlite_enabled == "true" || sqlite_enabled == "1" {
		utility.InitSQLite()
	}
}

type Middleware func(http.Handler) http.Handler

// applyMiddlewares 应用给定的中间件到 HTTP 处理器。
// 参数:
//   - h: 初始的 http.Handler，后续的中间件将应用于此。
//   - middlewares: 可变参数列表，包含依次应用的中间件函数。
//
// 返回值:
//   - 一个应用了所有中间件的 http.Handler。
func applyMiddlewares(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func main() {
	// 创建一个新的 ServeMux
	mux := http.NewServeMux()

	// 应用多个中间件到 mux
	handler := applyMiddlewares(mux,
		middleware.LogRequest, // 添加请求日志中间件
		middleware.APIVersionMiddleware,
		middleware.CORSMiddleware,
		middleware.SecurityHeadersMiddleware,
	)
	log.Println("中间件已加载")

	// 加载 PostgreSQL 路由
	postgres_enabled := os.Getenv("POSTGRES_ENABLED")
	if postgres_enabled == "true" || postgres_enabled == "1" {
		postgresRepo := repository.NewPostgresRepo(utility.Postgres)
		postgresService := service.NewApplicationService(postgresRepo)
		router.LoadPostgresRouter(mux, "/crate-api-data", postgresService)
	}

	// 加载 MySQL 路由
	mysql_enabled := os.Getenv("MYSQL_ENABLED")
	if mysql_enabled == "true" || mysql_enabled == "1" {
		mysqlRepo := repository.NewMySQLRepo(utility.MySQL)
		mysqlService := service.NewApplicationService(mysqlRepo)
		router.LoadMySQLRouter(mux, "/crate-api-data", mysqlService)
	}

	// 加载 SQLite 路由
	sqlite_enabled := os.Getenv("SQLITE_ENABLED")
	if sqlite_enabled == "true" || sqlite_enabled == "1" {
		sqliteRepo := repository.NewSQLiteRepo(utility.SQLite)
		sqliteService := service.NewApplicationService(sqliteRepo)
		router.LoadSQLiteRouter(mux, "/crate-api-data", sqliteService)
	}

	// 获取端口号并启动 HTTP 服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8421"
	}
	log.Println("0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, handler))
}

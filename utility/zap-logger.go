package utility

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var ZapLogger *zap.Logger

func InitZapLogger() {
	// 创建日志目录
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic(fmt.Sprintf("无法创建日志目录: %v", err))
	}

	// 配置编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	// 1. 控制台输出配置 (info级别，普通格式)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel, // 终端只显示 info 及以上级别
	)

	// 2. 文件输出配置 (warn级别，JSON格式)
	logFileName := fmt.Sprintf("logs/crate.%s.log", time.Now().Format("2006-01-02"))
	fileWriter := &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    100,  // MB
		MaxBackups: 90,   // 保留90个备份文件
		MaxAge:     90,   // 保留90天
		Compress:   true, // 压缩旧文件
	}

	// JSON 编码器配置
	jsonEncoderConfig := zap.NewProductionEncoderConfig()
	jsonEncoderConfig.TimeKey = "timestamp"
	jsonEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(jsonEncoderConfig)
	fileCore := zapcore.NewCore(
		fileEncoder,
		zapcore.AddSync(fileWriter),
		zapcore.WarnLevel, // 文件只记录 warn 及以上级别
	)

	// 合并两个输出核心
	core := zapcore.NewTee(consoleCore, fileCore)

	// 创建 logger
	ZapLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// 设置为全局 logger
	zap.ReplaceGlobals(ZapLogger)
}

// 优雅关闭日志器
func CloseZapLogger() {
	if ZapLogger != nil {
		ZapLogger.Sync()
	}
}

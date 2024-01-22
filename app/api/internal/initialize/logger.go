package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gobang/app/api/global"
	"os"
)

func SetupLogger() {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 设置全局日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.DebugLevel)

	// 创建控制台输出
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleDebugging, atomicLevel)

	// 创建文件输出
	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	fileDebugging := zapcore.AddSync(file)
	fileEncoder := zapcore.NewJSONEncoder(encoderCfg)
	fileCore := zapcore.NewCore(fileEncoder, fileDebugging, atomicLevel)

	// 创建 Logger
	logger1 := zap.New(zapcore.NewTee(consoleCore, fileCore))
	global.Logger = logger1

}

package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func InitLogger(logfile string) {
	Logger = Init(logfile)
}

func Init(logfile string) *zap.Logger {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logfile,
		MaxSize:    128, // megabytes for MB
		MaxBackups: 9999,
		MaxAge:     1,    // days
		Compress:   true, // disabled by default
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	// dev mod
	caller := zap.AddCaller()
	// filename and line
	development := zap.Development()
	// initial app name
	filed := zap.Fields(zap.String("serviceName", "gateway"))
	logger := zap.New(core, caller, development, filed)
	defer logger.Sync()
	return logger
}

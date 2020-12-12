package logger

import (
	"os"

	"go.uber.org/zap"
)

var zapLogger *zap.SugaredLogger

func init() {
	l, _ := zap.NewProduction()
	if os.Getenv("ENV") == "development" {
		l, _ = zap.NewDevelopment()
	}
	zapLogger = l.Sugar()
}

func sync() {
	if err := zapLogger.Sync(); err != nil {
		zapLogger.Errorw("logger sync error", "error", err.Error())
	}
}

func Debug(msg string, keysAndValues ...interface{}) {
	defer sync()

	zapLogger.Debugw(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	defer sync()

	zapLogger.Errorw(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	defer sync()

	zapLogger.Warnw(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...interface{}) {
	defer sync()

	zapLogger.Infow(msg, keysAndValues...)
}

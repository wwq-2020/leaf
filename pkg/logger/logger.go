package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/wwq1988/leaf/pkg/conf"
)

var std *zap.Logger

func init() {
	std = New()
}

func New() *zap.Logger {
	var level zapcore.Level
	logLevel := conf.GetLogLevel()
	if err := level.Set(logLevel); err != nil {
		level = zapcore.InfoLevel
	}
	encoderCfg := zap.NewProductionEncoderConfig()

	atomicLevel := zap.NewAtomicLevelAt(level)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atomicLevel,
	)
	logger := zap.New(core)
	return logger
}

func Fatal(msg string, fields ...zap.Field) {
	std.Fatal(msg, fields...)
}

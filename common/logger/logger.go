package logger

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"sync"

	"go.uber.org/zap"
)

var (
	l       *zap.Logger
	once    sync.Once
	appName string
)

func Init(w io.Writer) {
	once.Do(func() {
		appName = os.Getenv("APP_NAME")

		config := zapcore.EncoderConfig{}
		config.EncodeTime = zapcore.ISO8601TimeEncoder

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(config),
			zapcore.AddSync(w),
			zapcore.InfoLevel,
		)

		l = zap.New(core)
	})
}

func Get() *zap.Logger {
	return l
}

func Info(msg string, fields ...zap.Field) {
	fields = append([]zap.Field{
		zap.String("message", msg),
		zap.String("service", appName),
	}, fields...)

	l.Info("info", fields...)
}

func Errorf(format string, err error, fields ...zap.Field) {
	fields = append([]zap.Field{
		zap.Error(fmt.Errorf(format, err)),
		zap.String("service", appName),
	}, fields...)

	l.Error("error", fields...)
}

func Fatalf(format string, a ...any) {
	l.Fatal("fatal", zap.Error(
		fmt.Errorf(format, a),
	))
}

package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
		zap.Time("timestamp", time.Now()),
	}, fields...)

	l.Info("info", fields...)
}

func Errorf(format string, err error, fields ...zap.Field) {
	fields = append([]zap.Field{
		zap.Error(fmt.Errorf(format, err)),
		zap.String("service", appName),
		zap.Time("timestamp", time.Now()),
	}, fields...)

	l.Error("error", fields...)
}

func Fatalf(format string, a ...any) {
	l.Fatal("fatal", zap.Error(
		fmt.Errorf(format, a...),
	),
		zap.Time("timestamp", time.Now()),
	)

}

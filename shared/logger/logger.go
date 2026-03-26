package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(env string) {
	var encoder zapcore.EncoderConfig

	if env == "production" {
		encoder = zap.NewProductionEncoderConfig()
	} else {
		encoder = zap.NewDevelopmentEncoderConfig()
	}

	encoder.TimeKey = "timestamp"
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)

	log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func Get() *zap.Logger {
	if log == nil {
		// fallback if Init wasn't called
		log, _ = zap.NewDevelopment()
	}
	return log
}

func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

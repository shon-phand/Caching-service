package logger

import (
	"os"

	"github.com/shon-phand/CryptoServer/utils/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.Logger
)
var logHome string = os.Getenv("HOME") + "/logs.txt"

func init() {
	logConfig := zap.Config{
		OutputPaths: []string{"stdout", logHome},
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:    "level",
			TimeKey:     "time",
			MessageKey:  "msg",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		},
	}
	var err error
	if log, err = logConfig.Build(); err != nil {
		panic(err)
	}
}

func Info(restErr *errors.RestErr, stacktrace error, tags ...zap.Field) {
	tags = append(tags, zap.Int("HTTPStatus", restErr.Status))
	tags = append(tags, zap.String("Code", restErr.Code))
	tags = append(tags, zap.String("Message", restErr.Message))
	tags = append(tags, zap.NamedError("stacktrace", stacktrace))
	log.Info("info-error", tags...)
	log.Sync()
}

func Error(restErr *errors.RestErr, stacktrace error, tags ...zap.Field) {
	tags = append(tags, zap.Int("HTTPStatus", restErr.Status))
	tags = append(tags, zap.String("Code", restErr.Code))
	tags = append(tags, zap.String("Message", restErr.Message))
	tags = append(tags, zap.NamedError("staccktrace", stacktrace))
	log.Error("service-error", tags...)
	log.Sync()
}

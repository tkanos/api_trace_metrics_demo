package logger

import (
	"github.com/pkg/errors"
	"github.com/tkanos/api_trace_metrics_demo/app/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//Log ...
var (
	logStdOut *zap.Logger
	logStdErr *zap.Logger
)

//Init : instantiate a new logger
func Init() {
	conf := zap.NewProductionConfig()
	conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.EncoderConfig.MessageKey = "message"
	conf.EncoderConfig.TimeKey = "timestamp"

	level := zap.NewAtomicLevel()
	level.SetLevel(zap.DebugLevel)

	conf.Level = level
	conf.OutputPaths = []string{"stdout"}
	logStdOut, _ = conf.Build()

	conf.OutputPaths = []string{"stderr"}
	logStdErr, _ = conf.Build()
}

//Error logs errors into std err
func Error(msg string, err error, fields ...zapcore.Field) {
	logStdErr.Error(errors.Wrap(err, msg).Error(), fields...)
}

//Debug logs info into std out
func Debug(msg string, fields ...zapcore.Field) {
	if config.Verbose {
		logStdOut.Debug(msg, fields...)
	}
}

//Info logs info into std out
func Info(msg string, fields ...zapcore.Field) {
	logStdOut.Info(msg, fields...)
}

// StringField ...
func StringField(name string, value string) zapcore.Field {
	return zap.String(name, value)
}

// IntField ...
func IntField(name string, id int) zapcore.Field {
	return zap.Int(name, id)
}

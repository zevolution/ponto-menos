package logger

import (
	"fmt"
	"ponto-menos/cmd/cli/config/env"

	"bytes"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DEFAULT_MAX_SIZE_IN_MB  int64  = 5
	DEFAULT_MAX_AGE_IN_DAYS int64  = 7
	DEFAULT_LOG_FILE_NAME   string = "ponto-menos.log"
	DEFAULT_LOG_PATH        string = "."
	DEFAULT_LOG_LEVEL       string = "DEBUG"
)

func init() {
	Configure()
}

func Configure() {
	core := zapcore.NewTee(configureFileLogger())
	logger := zap.New(core)
	zap.ReplaceGlobals(logger)
}

func configureFileLogger() zapcore.Core {
	var b bytes.Buffer
	b.WriteString(env.GetOrDefault("logger.path", DEFAULT_LOG_PATH).(string))
	b.WriteString("/")
	b.WriteString(env.GetOrDefault("logger.filename", DEFAULT_LOG_FILE_NAME).(string))
	logfile := b.String()
	b.Reset()

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logfile,
		MaxSize:    int(env.GetOrDefault("logger.maxSize", DEFAULT_MAX_SIZE_IN_MB).(int64)),
		MaxBackups: 3,
		MaxAge:     int(env.GetOrDefault("logger.maxAge", DEFAULT_MAX_AGE_IN_DAYS).(int64)),
	})

	fileEncoderCfg := zap.NewDevelopmentEncoderConfig()
	fileEncoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	fileEncoder := zapcore.NewConsoleEncoder(fileEncoderCfg)
	level, err := zapcore.ParseLevel(env.GetOrDefault("logger.level", DEFAULT_LOG_LEVEL).(string))
	if err != nil {
		panic(fmt.Sprintf("Invalid zapcore level: %v\n", env.GetOrDefault("logger.level", DEFAULT_LOG_LEVEL).(string)))
	}

	zap.L().Debug(fmt.Sprintf("The ponto-menos file-logger was created with '%v' level\n", level))

	return zapcore.NewCore(fileEncoder, file, zap.NewAtomicLevelAt(level))
}

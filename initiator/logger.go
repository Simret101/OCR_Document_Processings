package initiator

import (
	"os"

	"aidoc/platform/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Info struct {
	Logger logger.Logger
}

func InitLogger() logger.Logger {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/info/info.log",
		MaxSize:    10,
		MaxAge:     28,
		MaxBackups: 3,
	})
	e := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/error/error.log",
		MaxSize:    10,
		MaxAge:     28,
		MaxBackups: 3,
	})

	
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.InfoLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	console := zapcore.Lock(os.Stdout)
	core := zapcore.NewTee(
		zapcore.NewCore(getEncoder(), w, infoLevel),
		zapcore.NewCore(getEncoder(), e, errorLevel),
		zapcore.NewCore(getEncoder(), console, zapcore.DebugLevel),
	)
	return logger.New(zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)))
}

func getEncoder() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(config)
}

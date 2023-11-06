package logger

import (
	"context"
	"os"
	"path"

	"kredit-plus/app/constants"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	SugarLogger       *zap.SugaredLogger
	AccessSugarLogger *zap.SugaredLogger
)

func InitLogger() {
	writerSyncer := getLogWriter()
	accessWriterSyncer := getAccessLogWriter()

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	core := zapcore.NewTee(
		zapcore.NewCore(getConsoleEncoder(), zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		zapcore.NewCore(getFileEncoder(), writerSyncer, zapcore.InfoLevel),
	)
	logger := zap.New(core, zap.AddCaller())
	SugarLogger = logger.Sugar()

	accessCore := zapcore.NewTee(
		zapcore.NewCore(getConsoleEncoder(), zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		zapcore.NewCore(getFileEncoder(), accessWriterSyncer, zapcore.InfoLevel),
	)
	accessLogger := zap.New(accessCore)
	AccessSugarLogger = accessLogger.Sugar()
}

func getConsoleEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getFileEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Include the full path of the caller's file in the log
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder

	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	logFilePath := constants.Config.LogConfig.LOG_FILE_PATH
	logFileName := constants.Config.LogConfig.LOG_FILE_NAME
	logFileMaxSize := constants.Config.LogConfig.LOG_FILE_MAXSIZE
	logFileMaxBackups := constants.Config.LogConfig.LOG_FILE_MAXBACKUP
	logFileMaxAge := constants.Config.LogConfig.LOG_FILE_MAXAGE
	logFile := path.Join(logFilePath, logFileName)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    logFileMaxSize,
		MaxBackups: logFileMaxBackups,
		MaxAge:     logFileMaxAge,
		Compress:   true,
		LocalTime:  true,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// Logger returns a zap logger with as much context as possible
func Logger(ctx context.Context) *zap.SugaredLogger {
	newLogger := SugarLogger
	if ctxCorrelationID, ok := ctx.Value(constants.CORRELATION_KEY_ID).(string); ok {
		newLogger = newLogger.With(zap.String(constants.CORRELATION_KEY_ID.String(), ctxCorrelationID))
	}

	return newLogger
}

func getAccessLogWriter() zapcore.WriteSyncer {
	accessLogFilePath := constants.Config.LogConfig.ACCESS_LOG_FILE_PATH
	accessLogFileName := constants.Config.LogConfig.ACCESS_LOG_FILE_NAME
	accessLogFileMaxSize := constants.Config.LogConfig.ACCESS_LOG_FILE_MAXSIZE
	accessLogFileMaxBackups := constants.Config.LogConfig.ACCESS_LOG_FILE_MAXBACKUP
	accessLogFileMaxAge := constants.Config.LogConfig.ACCESS_LOG_FILE_MAXAGE
	accessLogFile := path.Join(accessLogFilePath, accessLogFileName)
	accessLumberJackLogger := &lumberjack.Logger{
		Filename:   accessLogFile,
		MaxSize:    accessLogFileMaxSize,
		MaxBackups: accessLogFileMaxBackups,
		MaxAge:     accessLogFileMaxAge,
		Compress:   true,
		LocalTime:  true,
	}
	return zapcore.AddSync(accessLumberJackLogger)
}

// AccessLogger returns a separate zap logger for access logs
func AccessLogger() *zap.SugaredLogger {
	return AccessSugarLogger
}

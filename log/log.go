package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var Logger *zap.Logger

// InitZap 初始化zap logger
func InitZap() *zap.Logger {
	Logger = zap.New(getEncoderCore(), zap.AddCaller())
	Logger.Info("logger init success")
	return Logger
}

func getEncoderCore() zapcore.Core {
	core := zapcore.NewCore(getEncoder(), getWriteSyncer(), zap.DebugLevel)
	return core
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

func getEncoderConfig() zapcore.EncoderConfig {
	zap.NewProductionEncoderConfig()
	config := zapcore.EncoderConfig{
		CallerKey:     "caller_line", // 打印文件名和行数
		LevelKey:      "level_name",
		MessageKey:    "msg",
		TimeKey:       "ts",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		// 自定义时间格式
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString("[" + t.Format("2006/01/02 15:04:05.000") + "]")
		},
		// loglevel 大写编码器
		EncodeLevel: func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString("[" + level.CapitalString() + "]")
		},
		// 全路径编码器
		EncodeCaller: func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(caller.TrimmedPath())
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	return config
}

func getWriteSyncer() zapcore.WriteSyncer {
	//file, err := os.OpenFile("./info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	//if err != nil {
	//	return nil
	//}
	return zapcore.AddSync(os.Stdout)
}
